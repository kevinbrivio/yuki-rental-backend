package server

import (
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kevinbrivio/yuki-rental-backend/internal/config"
)

type Server struct {
	cfg *config.Config
	db *pgxpool.Pool
	log *slog.Logger
}

func New(cfg *config.Config, db *pgxpool.Pool, log *slog.Logger) *Server {
	log.Info("Server created")
	return &Server{
		cfg: cfg,
		db: db,
		log: log,
	}
}