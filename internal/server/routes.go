package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/kevinbrivio/yuki-rental-backend/internal/middleware"
)

func (s *Server) routes() chi.Router {
	r := chi.NewRouter()

	// middleware goes here
	r.Use(middleware.Logger(s.log))
	r.Use(middleware.Recovery(s.log))

	// router goes here
	r.Get("/health", s.handleHealth)

	return r
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	w.Header().Set("Content-Type", "application/json")

	// Check db
	err := s.db.Ping(ctx)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "database": "down"})
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "database": "up"})
	}
}
