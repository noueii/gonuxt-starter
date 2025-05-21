package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	dbrun "github.com/noueii/gonuxt-starter/internal/db"
	db "github.com/noueii/gonuxt-starter/internal/db/out"
	"github.com/noueii/gonuxt-starter/internal/gapi"
	"github.com/noueii/gonuxt-starter/internal/util"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

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

	if environment == util.Production {
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

	err = dbrun.RunMigrations(cfg, conn)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot run migrations:")
	}

	queries := db.New(conn)

	waitGroup, ctx := errgroup.WithContext(ctx)

	gapi.RunGatewayServer(ctx, waitGroup, cfg, queries)
	gapi.RunGRPCServer(ctx, waitGroup, cfg, queries)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("wait group error")
	}

}
