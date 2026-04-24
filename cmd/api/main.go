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
	var logger *slog.Logger
	if cfg.App.IsProd() {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	
	// 3. Create to database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	db, err := database.New(ctx, cfg.DB)
	if err != nil {
		return fmt.Errorf("creating database: %w", err)
	}
	defer db.Close()

	// 4. Create server then run it
	srv := server.New(cfg, db, logger)
	
	return srv.Run()
}
