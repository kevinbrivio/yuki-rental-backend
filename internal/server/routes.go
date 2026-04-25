package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

func (s *Server) routes() chi.Router {
	r := chi.NewRouter()
	
	// middleware goes here
	// 
	// router goes here
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		// Check db
		err := s.db.Ping(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{"status": "ok", "database": "down"})
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "ok", "database": "up"})
		}
	})
	
	return r
}