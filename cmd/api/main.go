package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/kevinbrivio/yuki-rental-backend/internal/config"
	"github.com/kevinbrivio/yuki-rental-backend/internal/database"
	"github.com/kevinbrivio/yuki-rental-backend/internal/server"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error %v", err)
		os.Exit(1)
	}
}

func run() error {
	// 1. Load the config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	
	// 2. Logger
	// Use JSONHandler if Prod, else use TextHandler
	var logger *zap.Logger
	if cfg.App.IsProd() {
		logger, err = zap.NewProduction()
		if err != nil {
			return err
		}
	} else {
		logger, err = zap.NewDevelopment()
		if err != nil {
			return err
		}
	}
	
	slogLogger := slog.New(zapslog.NewHandler(logger.Core()))
	
	// 3. Create to database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	db, err := database.New(ctx, cfg.DB)
	if err != nil {
		return fmt.Errorf("creating database: %w", err)
	}
	defer db.Close()

	// 4. Create server then run it
	srv := server.New(cfg, db, slogLogger)
	
	return srv.Run()
}
