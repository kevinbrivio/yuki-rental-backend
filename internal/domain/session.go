package domain

import "time"

type Session struct {
	ID string `json:"id,omitempty"`
	UserID string `json:"user_id"`
	TokenHash string `json:"-"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	ExpiresAt time.Time `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at"`
	CreatedAt time.Time `json:"created_at"`
}