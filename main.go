package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/noueii/gonuxt-starter/api"
	db "github.com/noueii/gonuxt-starter/db/out"
	"github.com/noueii/gonuxt-starter/gapi"
	"github.com/noueii/gonuxt-starter/pb"
	"github.com/noueii/gonuxt-starter/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {

	cfg, err := util.LoadConfig(util.Development)

	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Connecting db")

	conn, err := sql.Open(cfg.DbDriver, cfg.DbURL)

	if err != nil {
		log.Fatal("cannot connect to db:", err)
		fmt.Println(err)
	}

	queries := db.New(conn)

	go runGatewayServer(cfg, queries)
	runGRPCServer(cfg, queries)

}

func runGRPCServer(config *util.Config, queries *db.Queries) {
	server, err := gapi.NewServer(config, queries)
	if err != nil {
		fmt.Println(err)
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGoNuxtServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCAddr)
	if err != nil {
		fmt.Println(err)
		log.Fatal("cannot create gRPC listener:", err)
	}

	log.Printf("starting gRPC server at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		fmt.Println(err)
		log.Fatal("failed to start gRPC server:", err)
	}

}

func runGatewayServer(config *util.Config, queries *db.Queries) {
	server, err := gapi.NewServer(config, queries)
	if err != nil {
		fmt.Println(err)
		log.Fatal("cannot create server:", err)
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
		log.Fatal("cannot register handler server:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HTTPAddr)
	if err != nil {
		fmt.Println(err)
		log.Fatal("cannot create gRPC listener:", err)
	}

	log.Printf("starting HTTP gateway server at %s", listener.Addr().String())

	err = http.Serve(listener, mux)
	if err != nil {
		fmt.Println(err)
		log.Fatal("failed to start HTTP gateway server:", err)
	}

}

func runGinServer(config *util.Config, queries *db.Queries) {
	httpServer, err := api.NewServer(config, queries)

	if err != nil {
		log.Fatal("cannot create http server:", err)
		fmt.Println(err)
	}

	fmt.Println("Starting server")
	err = httpServer.Start(config.HTTPAddr)
	if err != nil {
		log.Fatal("cannot start http server:", err)
		fmt.Println(err)
	}
}
