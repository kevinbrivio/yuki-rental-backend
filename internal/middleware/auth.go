package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/kevinbrivio/yuki-rental-backend/internal/repository"
	"github.com/kevinbrivio/yuki-rental-backend/internal/response"
)

type contextKey string
const SessionIDKey contextKey = "sessionID"

func Auth(sr repository.SessionRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Get the cookie
			cookie, err := r.Cookie("session_token")
			if err != nil {
				response.WriteError(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			
			// 2. Hash the token
			hash := sha256.Sum256([]byte(cookie.Value))
			tokenHash := hex.EncodeToString(hash[:])
			// 3. Check in DB
			session, err := sr.GetSessionByToken(r.Context(), tokenHash)
			if err != nil {
				response.WriteError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			// 4. Inject session ID to context
			ctx := context.WithValue(r.Context(), SessionIDKey, session.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}