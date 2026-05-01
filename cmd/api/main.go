package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/kevinbrivio/yuki-rental-backend/internal/config"
	"github.com/kevinbrivio/yuki-rental-backend/internal/database"
	"github.com/kevinbrivio/yuki-rental-backend/internal/handler"
	"github.com/kevinbrivio/yuki-rental-backend/internal/repository"
	"github.com/kevinbrivio/yuki-rental-backend/internal/server"
	"github.com/kevinbrivio/yuki-rental-backend/internal/service"
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

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("getting underlying error: %w", err)
	}
	defer sqlDB.Close()

	// 4. Repositories
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	// 5. Services
	authService := service.NewAuthService(userRepo, sessionRepo)

	// 6. Handlers
	handlers := handler.Handlers{
		Auth: handler.NewAuthHandler(authService),
	}

	// 7. Create server with router
	srv := server.New(cfg, db, slogLogger, &handlers)

	return srv.Run()
}
