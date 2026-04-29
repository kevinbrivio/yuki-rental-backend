package database

import (
	"context"
	"fmt"

	"github.com/kevinbrivio/yuki-rental-backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)	

func New(ctx context.Context, cfg config.DBConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	
	// verifies the connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("getting underlying db: %w", err)
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}
	
	
	return db, nil
}