package dto

import (
	"time"

	"github.com/kevinbrivio/yuki-rental-backend/internal/domain"
)

type RegisterRequest struct {
	Email       string
	Password    string
	PhoneNumber string
	FirstName   string
	LastName    string
	DateOfBirth time.Time
	Gender      domain.GenderType
}

type LoginRequest struct {
	Email     string
	Password  string
	IPAddress string
	UserAgent string
}

type LoginResponse struct {
	Session  *domain.Session
	RawToken string
}
