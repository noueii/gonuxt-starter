package dbrun

import (
	"database/sql"
	"embed"

	"github.com/noueii/gonuxt-starter/internal/util"
	"github.com/pressly/goose/v3"
)

//go:embed schema/*.sql
var embedMigrations embed.FS

func RunMigrations(config *util.Config, db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect(config.DbDriver); err != nil {
		return err
	}

	if err := goose.Up(db, config.DBMigrationsLocation); err != nil {
		return err
	}
	return nil
}
