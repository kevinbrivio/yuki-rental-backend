package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kevinbrivio/yuki-rental-backend/internal/config"
	"github.com/kevinbrivio/yuki-rental-backend/internal/handler"
	"gorm.io/gorm"
)

type Server struct {
	cfg *config.Config
	db *gorm.DB
	log *slog.Logger
	handlers *handler.Handlers
}

func New(cfg *config.Config, db *gorm.DB, log *slog.Logger, h *handler.Handlers) *Server {
	log.Info("Server created")
	return &Server{
		cfg: cfg,
		db: db,
		log: log,
		handlers: h,
	}
}

func (s *Server) Run() error {
	// 1. Create new http.Server:
	srv := &http.Server{
		Addr: s.cfg.App.Addr(),
		Handler: s.setupRoutes(),
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout: 60 * time.Second,
	}
	
	// 2. Start listening with goroutines
	errChan := make(chan error, 1)
	
	s.log.Info("Server starting", slog.String("addr", srv.Addr))
	go func() {
		if err := srv.ListenAndServe(); err != nil  && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()	
	
	// 3. Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
		case err := <-errChan:
			// server failed to start (port in use, etc.)
			s.log.Error("Server failed to start due to -> %w", slog.Any("error", err))
			return fmt.Errorf("server error: %w", err)
		case sig := <-quit:
			// got ctrl+c from SIGTERM
			s.log.Info("shutdown signal received", slog.String("signal", sig.String()))
	}
	
	// 4. Call srv.Shutdown(ctx) with timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	
	return nil
}
