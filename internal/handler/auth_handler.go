package handler

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/kevinbrivio/yuki-rental-backend/internal/dto"
	"github.com/kevinbrivio/yuki-rental-backend/internal/middleware"
	"github.com/kevinbrivio/yuki-rental-backend/internal/response"
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
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	
	// 2. Call service
	err := ah.authService.Register(r.Context(), req)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusCreated, map[string]string{"message": "registered successfully"})
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	// Split the port because r.RemoteAddr combines host:port
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	req.IPAddress = host
	req.UserAgent = r.Header.Get("User-Agent")

	res, err := ah.authService.Login(r.Context(), req)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name: "session_token",
		Value: res.RawToken,
		Path: "/",
		HttpOnly: true, // => This means JS Cannot read the cookie
		Secure: false, // => Use https (TRUE if prod)
		SameSite: http.SameSiteStrictMode,
		Expires: res.Session.ExpiresAt,
	})

	response.WriteJSON(w, http.StatusOK, map[string]any{"message": "login successfully"})
}

func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// 1. Retrieve the sessionID from url query
	sessionID, ok := r.Context().Value(middleware.SessionIDKey).(string)
	if !ok|| sessionID == "" {
		response.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if sessionID == "" { // when there's no session
		response.WriteError(w, http.StatusBadRequest, "invalid session")
		return
	}
	
	if err := ah.authService.Logout(r.Context(), sessionID); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"message": "logout successfully"})
}
