// Package database provides access to the application's SQLite database.
package database

import (
	"database/sql"
	"log/slog"
	"time"

	// Import the SQLite driver
	// The underscore import registers the driver with database/sql
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

type Database struct {
	db *sql.DB
}

func Init(dbPath string, migrationDir string) (*Database, error) {
	slog.Info("initializing database", "path", dbPath)

	// PRAGMA settings on startup
	dsn := dbPath + "?_journal_mode=WAL&_busy_timeout=5000&_synchronous=NORMAL&_cache_size=-64000&_foreign_keys=ON"
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		slog.Error("database initialization failed", "path", dbPath, "err", err)
		return nil, err
	}

	db.SetMaxOpenConns(1) // SQLite supports only one writer at a time
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		slog.Error("failed to ping database", "err", err)
		return nil, err
	}

	if err := goose.SetDialect("sqlite3"); err != nil {
		slog.Error("failed to set goose dialect:", "err", err)
		return nil, err
	}

	err = goose.Up(db, migrationDir)
	if err != nil {
		slog.Error("failed to run migrations", "err", err)
		return nil, err
	}

	return &Database{
		db,
	}, nil
}
