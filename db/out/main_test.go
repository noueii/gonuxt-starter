package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/noueii/gonuxt-starter/util"
)

var testQueries *Queries

func TestMain(m *testing.M) {

	cfg, err := util.LoadConfig(util.Test, "../")

	if err != nil {
		log.Fatal("error loading environment variables: ", err)
	}

	conn, err := sql.Open(cfg.DbDriver, cfg.DbURL)

	if err != nil {
		log.Fatal("can't connect to db: ", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
