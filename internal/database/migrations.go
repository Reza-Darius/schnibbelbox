package database

import (
	"database/sql"
	"embed"
	"log/slog"

	"github.com/pressly/goose/v3"
)

//go:embed migrations
var migrationFiles embed.FS

func runMigration(db *sql.DB, dialect goose.Dialect) error {
	if err := goose.SetDialect(string(dialect)); err != nil {
		slog.Error("failed to set goose dialect:", "err", err)
		return err
	}

  goose.SetBaseFS(migrationFiles)

	err := goose.Up(db, "migrations")
	if err != nil {
		slog.Error("failed to run migrations", "err", err)
		return err
	}
	return nil
}
