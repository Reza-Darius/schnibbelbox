package database

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type PGDatabase struct {
	pool *pgxpool.Pool
}

func InitPG() (*PGDatabase, error) {
	dbPath := os.Getenv("DATABASE_URL")
	if dbPath == "" {
		slog.Error("couldnt retrieve database URL")
		return nil, errors.New("couldnt retrieve DATABASE_URL env variable")
	}

	slog.Info("initializing PG database", "path", dbPath)

	pool, err := pgxpool.New(context.Background(), dbPath)
	if err != nil {
		slog.Error("database initialization failed", "path", dbPath, "err", err)
		return nil, err
	}

	// we create a context so the ping call doesnt hang indefinitely in case the pg server is not reachable
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		slog.Error("failed to ping database", "err", err)
		return nil, err
	}

	if err := goose.SetDialect(string(goose.DialectPostgres)); err != nil {
		pool.Close()
		slog.Error("failed to set goose dialect:", "err", err)
		return nil, err
	}

  goose.SetBaseFS(migrationFiles)

	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	err = goose.Up(db, "migrations")
	if err != nil {
		pool.Close()
		slog.Error("failed to run migrations", "err", err)
		return nil, err
	}

	return &PGDatabase{
		pool,
	}, nil
}
