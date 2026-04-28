package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"
)

func Recovery(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec != nil {
				logger.Error("panic recovered", "error", rec, "stack", string(debug.Stack()))
				jsonBody, _ := json.Marshal(map[string]string {
					"error": "Internal server error",
				})
				
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonBody)
			}
		}()
		next.ServeHTTP(w, r)
	})
	}
}
