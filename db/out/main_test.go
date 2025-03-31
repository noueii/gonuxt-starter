package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var testQueries *Queries

const (
	dbDriver = "postgres"
	dbSource = "postgres://root:gonuxtsecret@localhost:5432/gonuxt?sslmode=disable"
)

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("can't connect to db: ", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
