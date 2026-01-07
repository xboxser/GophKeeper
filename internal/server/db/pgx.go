package db

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"gophkeeper/internal/config"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB interface {
	Close()
	Query(context.Context, string, ...any) (pgx.Rows, error)
}

type DBPgx struct {
	pool *pgxpool.Pool
}

func NewDBPgx(ctx context.Context, cfg config.ServerConfig) (*DBPgx, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, err
	}

	// настройка пула
	poolCfg.MaxConns = 10
	poolCfg.MinConns = 2
	poolCfg.MaxConnLifetime = time.Hour
	poolCfg.MaxConnIdleTime = 10 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	// проверка соединения
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	err = runMigrations(cfg.DSN)
	if err != nil {
		pool.Close()
		return nil, err
	}

	return &DBPgx{pool: pool}, nil
}

func (db *DBPgx) Close() {
	db.pool.Close()
}

func (db *DBPgx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return db.pool.Query(ctx, sql, args...)
}

//go:embed migrations/*.sql
var migrationsDir embed.FS

func runMigrations(dsn string) error {
	d, err := iofs.New(migrationsDir, "migrations")
	if err != nil {
		return fmt.Errorf("failed to return an iofs driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		return fmt.Errorf("failed to get a new migrate instance: %w", err)
	}

	if err := m.Up(); err != nil {

		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations to the DB: %w", err)
		}
	}
	return nil
}
