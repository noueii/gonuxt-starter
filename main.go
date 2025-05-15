package main

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/noueii/gonuxt-starter/api"
	db "github.com/noueii/gonuxt-starter/db/out"
	"github.com/noueii/gonuxt-starter/gapi"
	"github.com/noueii/gonuxt-starter/pb"
	"github.com/noueii/gonuxt-starter/util"
	"github.com/pressly/goose/v3"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

//go:embed db/schema/*.sql
var embedMigrations embed.FS

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	environment, err := util.LoadEnv()

	if err != nil {
		log.Fatal().Err(err).Msg("could not load environment")
	}

	if environment == util.Development {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	cfg, err := util.LoadConfig(environment)

	if err != nil {
		log.Fatal().Err(err).Msg("could not load config")
	}

	conn, err := sql.Open(cfg.DbDriver, cfg.DbURL)

	if err != nil {
		log.Fatal().Msg("cannot connect to db:")
	}

	err = runMigrations(cfg, conn)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot run migrations:")
	}

	queries := db.New(conn)

	waitGroup, ctx := errgroup.WithContext(ctx)

	runGatewayServer(ctx, waitGroup, cfg, queries)
	runGRPCServer(ctx, waitGroup, cfg, queries)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("wait group error")
	}

}

func runMigrations(config *util.Config, db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect(config.DbDriver); err != nil {
		return err
	}

	if err := goose.Up(db, config.DBMigrationsLocation); err != nil {
		return err
	}
	return nil
}

func runGRPCServer(ctx context.Context, waitGroup *errgroup.Group, config *util.Config, queries *db.Queries) {
	server, err := gapi.NewServer(config, queries)
	if err != nil {
		fmt.Println(err)
		log.Fatal().Err(err).Msg("cannot create server:")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GRPCLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterGoNuxtServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCAddr)
	if err != nil {
		fmt.Println(err)
		log.Fatal().Err(err).Msg("cannot create gRPC listener:")
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("starting gRPC server at %s", listener.Addr().String())
		err = grpcServer.Serve(listener)

		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}

			log.Fatal().Err(err).Msg("cannot start gRPC server")
			return err
		}

		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("gracefully shutting down gRPC server")

		grpcServer.GracefulStop()
		log.Info().Msg("gRPC server shutdown")

		return nil
	})

}

func runGatewayServer(ctx context.Context, waitGroup *errgroup.Group, config *util.Config, queries *db.Queries) {
	server, err := gapi.NewServer(config, queries)
	if err != nil {
		fmt.Println(err)
		log.Fatal().Err(err).Msg("cannot create server:")
	}

	muxOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions:   protojson.MarshalOptions{UseProtoNames: true},
		UnmarshalOptions: protojson.UnmarshalOptions{DiscardUnknown: true},
	})

	grpcMux := runtime.NewServeMux(muxOption)

	err = pb.RegisterGoNuxtHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register handler server:")
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), gapi.ResponseWriterKey, w)

		grpcMux.ServeHTTP(w, r.WithContext(ctx))
	}))

	crs := cors.New(cors.Options{
		AllowedOrigins: config.CORSAllowedOrigins,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders: []string{
			"Authorization",
			"Content-Type",
		},
		AllowCredentials: true,
	})

	corsHandler := crs.Handler(gapi.HttpLogger(mux))

	httpServer := &http.Server{
		Handler: corsHandler,
		Addr:    config.HTTPAddr,
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("starting HTTP gateway server at %s", httpServer.Addr)

		err = httpServer.ListenAndServe()

		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}

			fmt.Println(err)
			log.Error().Err(err).Msg("failed to start HTTP gateway server:")
			return err
		}

		return nil

	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("gracefully shutting down HTTP gateway server")

		err = httpServer.Shutdown(context.Background())

		if err != nil {
			log.Error().Err(err).Msg("failed to shutdown HTTP gateway server")
			return err
		}

		return nil
	})

}

func runGinServer(ctx context.Context, waitGroup *errgroup.Group, config *util.Config, queries *db.Queries) {
	httpServer, err := api.NewServer(config, queries)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create http server:")
		fmt.Println(err)
	}

	waitGroup.Go(func() error {
		fmt.Println("Starting server")
		err = httpServer.Start(config.HTTPAddr)

		if err != nil {
			log.Fatal().Err(err).Msg("cannot start HTTP server")
			return err
		}

		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("gracefully shutting down HTTP server")

		err = httpServer.Shutdown()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}

			log.Error().Err(err).Msg("could not shutdown HTTP server")
			return err
		}

		return nil
	})

}
