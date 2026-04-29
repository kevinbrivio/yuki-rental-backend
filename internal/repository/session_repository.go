package repository

import (
	"context"
	"time"

	"github.com/kevinbrivio/yuki-rental-backend/internal/domain"
	"gorm.io/gorm"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session *domain.Session) error
	GetSessionByToken(ctx context.Context, tokenHash string) (*domain.Session, error)
	RevokeSession(ctx context.Context, sessionID string) error
}

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (sr *sessionRepository) CreateSession(ctx context.Context, session *domain.Session) error {
	return sr.db.WithContext(ctx).Create(session).Error
}

func (sr *sessionRepository) GetSessionByToken(ctx context.Context, tokenHash string) (*domain.Session, error) {
	var session domain.Session
	// When a user is active (has session), revoked_at should be null
	err := sr.db.WithContext(ctx).Where("token_hash = ? AND expires_at > ? AND revoked_at is NULL", tokenHash, time.Now()).
		First(&session).
		Error
	if err != nil {
		return nil, err
	}
	
	return &session, nil
}

func (sr *sessionRepository) RevokeSession(ctx context.Context, sessionID string) error {
	// Because .Update doesn't know about which table to look up
	// Use .Model() to point GORM which table... 
	return sr.db.WithContext(ctx).Model(&domain.Session{}).Where("id = ? ", sessionID).
		Update("revoked_at", time.Now()).
		Error
}