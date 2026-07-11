// Package database provides access to the application's SQLite database.
package database

import (
	"context"
	"database/sql"
	"log/slog"
	"runtime"
	"time"

	// Import the SQLite driver
	// The underscore import registers the driver with database/sql
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

type Database struct {
	readDB  *sql.DB
	writeDB *sql.DB
}

func InitSQLite(dbPath string) (*Database, error) {
	slog.Info("initializing SQLite database", "path", dbPath)

	// PRAGMA settings on startup
	dsn := dbPath + "?_journal_mode=WAL&_busy_timeout=5000&_synchronous=NORMAL&_cache_size=-64000&_foreign_keys=ON_txlock=immediate"
	readDB, err := sql.Open("sqlite3", dsn)
	if err != nil {
		slog.Error("database initialization failed", "path", dbPath, "err", err)
		return nil, err
	}

	maxReader := max(4, runtime.NumCPU())

	readDB.SetMaxOpenConns(maxReader)
	readDB.SetMaxIdleConns(maxReader)
	readDB.SetConnMaxLifetime(time.Hour)

	writeDB, err := sql.Open("sqlite3", dsn)
	if err != nil {
		slog.Error("database initialization failed", "path", dbPath, "err", err)
		return nil, err
	}

	// SQLite supports only one writer at a time, by limiting the connection to one we avoid a SQL SERVER BUSY error
	writeDB.SetMaxOpenConns(1)
	writeDB.SetMaxIdleConns(1)
	writeDB.SetConnMaxLifetime(time.Hour)

	// we create a context so the ping call doesnt hang indefinitely in case the pg server is not reachable
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := readDB.PingContext(ctx); err != nil {
		_ = readDB.Close()
		slog.Error("failed to ping database", "err", err)
		return nil, err
	}

	err = runMigration(readDB, goose.DialectSQLite3)
	if err != nil {
		_ = readDB.Close()
		slog.Error("failed to run migrations", "err", err)
		return nil, err
	}

	return &Database{
		readDB,
		writeDB,
	}, nil
}
