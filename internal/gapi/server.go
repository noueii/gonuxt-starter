package gapi

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	db "github.com/noueii/gonuxt-starter/internal/db/out"
	"github.com/noueii/gonuxt-starter/internal/pb"
	"github.com/noueii/gonuxt-starter/internal/token"
	"github.com/noueii/gonuxt-starter/internal/util"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

type Server struct {
	pb.UnimplementedGoNuxtServer
	pb.UnimplementedAuthServer
	config     *util.Config
	db         *db.Queries
	tokenMaker token.Maker
}

func NewServer(config *util.Config, db *db.Queries) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, err
	}

	server := &Server{
		db:         db,
		config:     config,
		tokenMaker: tokenMaker,
	}

	return server, nil
}

func RunGRPCServer(ctx context.Context, waitGroup *errgroup.Group, config *util.Config, queries *db.Queries) {
	server, err := NewServer(config, queries)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server:")
	}

	grpcLogger := grpc.UnaryInterceptor(GRPCLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterGoNuxtServer(grpcServer, server)
	pb.RegisterAuthServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCAddr)
	if err != nil {
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

func RunGatewayServer(ctx context.Context, waitGroup *errgroup.Group, config *util.Config, queries *db.Queries) {
	server, err := NewServer(config, queries)
	if err != nil {
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

	err = pb.RegisterAuthHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register handler server:")

	}

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ResponseWriterKey, w)

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

	corsHandler := crs.Handler(HttpLogger(mux))

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
