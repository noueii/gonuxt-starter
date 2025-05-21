package tests

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	db "github.com/noueii/gonuxt-starter/internal/db/out"
	"github.com/noueii/gonuxt-starter/internal/util"
)

var testQueries *db.Queries

func TestMain(m *testing.M) {

	cfg, err := util.LoadConfig(util.Test, "../")

	if err != nil {
		log.Fatal("error loading environment variables: ", err)
	}

	conn, err := sql.Open(cfg.DbDriver, cfg.DbURL)

	if err != nil {
		log.Fatal("can't connect to db: ", err)
	}

	testQueries = db.New(conn)

	os.Exit(m.Run())
}
