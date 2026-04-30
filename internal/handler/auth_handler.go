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
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "registered successfully"})
}