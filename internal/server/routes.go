package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/kevinbrivio/yuki-rental-backend/internal/middleware"
)

func (s *Server) setupRoutes() chi.Router {
	r := chi.NewRouter()

	// middleware goes here
	r.Use(middleware.Logger(s.log))
	r.Use(middleware.Recovery(s.log))

	// Health
	r.Get("/health", s.handleHealth)

	// Auth Handler
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", s.handlers.Auth.Register)
		r.Post("/login", s.handlers.Auth.Login)

		r.Group(func (r chi.Router) {
			// r.Use(authMiddleware) => TODO: Build this next
			r.Post("/logout", s.handlers.Auth.Logout)
		})
	})

	return r
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	w.Header().Set("Content-Type", "application/json")

	// Check db
	sqlDB, err := s.db.DB()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"status": "failed", "database": "down"})
	}
	
	err = sqlDB.PingContext(ctx)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"status": "failed", "database": "down"})
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "database": "up"})
	}
}
