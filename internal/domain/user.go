package domain

import (
	"time"

	"gorm.io/gorm"
)

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
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id,omitempty"`
	Email string `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"not null" json:"-"`
	PhoneNumber string `gorm:"uniqueIndex;not null" json:"phone_number"`
	Role UserRole `gorm:"type:user_role;not null" json:"role"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Gender GenderType `gorm:"type:gender_type;not null" json:"gender"`
	AvatarURL string `json:"avatar_url"`
	IsActive bool `json:"is_active"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}