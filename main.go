package main

import (
	"database/sql"
	"log"

	"fmt"

	_ "github.com/lib/pq"
	"github.com/noueii/gonuxt-starter/api"
	db "github.com/noueii/gonuxt-starter/db/out"
	"github.com/noueii/gonuxt-starter/util"
)

func main() {

	cfg, err := util.LoadConfig(util.Production)

	if err != nil {
		fmt.Println(err.Error())
	}

	conn, err := sql.Open(cfg.DbDriver, cfg.DbURL)

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	queries := db.New(conn)
	httpServer := api.NewServer(queries)

	err = httpServer.Start(cfg.HTTPAddr)
	if err != nil {
		log.Fatal("cannot start http server:", err)
	}

}
