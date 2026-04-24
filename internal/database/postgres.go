package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kevinbrivio/yuki-rental-backend/internal/config"
)	

func New(ctx context.Context, cfg config.DBConfig) (*pgxpool.Pool, error) {
	// parse dsn
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("parsing db config: %w", err)
	}
	
	// creates the pool
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("creating pool error: %w", err)
	}
	
	// verifies the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close() // clean up
		return nil, fmt.Errorf("pinging database: %w", err)
	}
	
	return pool, nil
}