package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kevinbrivio/yuki-rental-backend/internal/dto"
	"github.com/kevinbrivio/yuki-rental-backend/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(as service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: as,
	}
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// 1. Parse json body into dto
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	
	// 2. Call service
	err := ah.authService.Register(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "registered successfully"})
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.IPAddress = r.RemoteAddr
	req.UserAgent = r.Header.Get("User-Agent")

	res, err := ah.authService.Login(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name: "session_token",
		Value: res.RawToken,
		Path: "/",
		HttpOnly: true, // => This means JS Cannot read the cookie
		Secure: true, // => Use https
		SameSite: http.SameSiteStrictMode,
		Expires: res.Session.ExpiresAt,
	})

	writeJSON(w, http.StatusOK, map[string]any{"message": "login successfully"})
}

func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// 1. Retrieve the sessionID from url query
	sessionID := r.Context().Value("sessionID").(string)
	if sessionID == "" { // when there's no session
		writeError(w, http.StatusBadRequest, "invalid session")
		return
	}
	
	if err := ah.authService.Logout(r.Context(), sessionID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "logout successfully"})
}