package domain

import "time"

type UserRole string
const (
	Admin UserRole = "admin"
	Customer UserRole = "customer"
)

type GenderType string
const (
	Male GenderType = "male"
	Female GenderType = "female"
)

type User struct {
	ID string `json:"id,omitempty"`
	Email string `json:"email"`
	PasswordHash string `json:"-"`
	PhoneNumber string `json:"phone_number"`
	Role UserRole `json:"role"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Gender GenderType `json:"gender"`
	AvatarURL string `json:"avatar_url"`
	IsActive bool `json:"is_active"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at"`
}