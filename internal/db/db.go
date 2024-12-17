package db

import (
	"context"
	"embed"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
)

var embedMigrations embed.FS

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	PoolSize int
}

func Connect(ctx context.Context, cfg DBConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	if cfg.PoolSize > 0 {
		config.MaxConns = int32(cfg.PoolSize)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := PingDB(ctx, pool); err != nil {
		pool.Close()
		return nil, err
	}

	if err := MigrateDB(dsn); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Info().Msg("Database connection established and migrations applied.")
	return pool, nil
}

func PingDB(ctx context.Context, pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return pool.Ping(ctx)
}

func MigrateDB(dsn string) error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.UpFS(nil, dsn, embedMigrations, "migrations"); err != nil {
		return err
	}

	return nil
}
