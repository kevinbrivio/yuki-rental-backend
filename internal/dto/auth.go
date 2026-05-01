package dto

import (
	"time"

	"github.com/kevinbrivio/yuki-rental-backend/internal/domain"
)

type RegisterRequest struct {
	Email       string            `json:"email"`
	Password    string            `json:"password"`
	PhoneNumber string            `json:"phone_number"`
	FirstName   string            `json:"first_name"`
	LastName    string            `json:"last_name"`
	DateOfBirth time.Time         `json:"date_of_birth"`
	Gender      domain.GenderType `json:"gender"`
}

type LoginRequest struct {
	Email     string `json:"email"`
	Password  string `json:"passsword"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
}

type LoginResponse struct {
	Session  *domain.Session `json:"session"`
	RawToken string          `json:"raw_token"`
}
