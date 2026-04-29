package domain

import "time"

type Session struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id,omitempty"`
	UserID string `gorm:"type:uuid;not null" json:"user_id"`
	TokenHash string `gorm:"not null" json:"-"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	ExpiresAt time.Time `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at"`
	CreatedAt time.Time `json:"created_at"`
}