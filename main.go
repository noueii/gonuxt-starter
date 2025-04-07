package main

import (
	"context"
	"database/sql"
	"embed"
	"net"
	"net/http"
	"os"

	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/noueii/gonuxt-starter/api"
	db "github.com/noueii/gonuxt-starter/db/out"
	"github.com/noueii/gonuxt-starter/gapi"
	"github.com/noueii/gonuxt-starter/pb"
	"github.com/noueii/gonuxt-starter/util"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

//go:embed db/schema/*.sql
var embedMigrations embed.FS

func main() {
	environment, err := util.LoadEnv()

	if environment == util.Development {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if err != nil {
		log.Fatal().Err(err).Msg("could not load environment")
	}

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

	go runGatewayServer(cfg, queries)
	runGRPCServer(cfg, queries)

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

func runGRPCServer(config *util.Config, queries *db.Queries) {
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

	log.Printf("starting gRPC server at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		fmt.Println(err)
		log.Fatal().Err(err).Msg("failed to start gRPC server:")
	}

}

func runGatewayServer(config *util.Config, queries *db.Queries) {
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterGoNuxtHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register handler server:")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HTTPAddr)
	if err != nil {
		fmt.Println(err)
		log.Fatal().Err(err).Msg("cannot create gRPC listener:")
	}

	log.Info().Msgf("starting HTTP gateway server at %s", listener.Addr().String())

	handler := gapi.HttpLogger(mux)

	err = http.Serve(listener, handler)
	if err != nil {
		fmt.Println(err)
		log.Fatal().Err(err).Msg("failed to start HTTP gateway server:")
	}

}

func runGinServer(config *util.Config, queries *db.Queries) {
	httpServer, err := api.NewServer(config, queries)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create http server:")
		fmt.Println(err)
	}

	fmt.Println("Starting server")
	err = httpServer.Start(config.HTTPAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start http server:")
		fmt.Println(err)
	}
}
