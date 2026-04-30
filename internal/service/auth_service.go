package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/kevinbrivio/yuki-rental-backend/internal/domain"
	"github.com/kevinbrivio/yuki-rental-backend/internal/dto"
	"github.com/kevinbrivio/yuki-rental-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) error
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	Logout(ctx context.Context, sessionID string) error
}

func NewAuthService(ur repository.UserRepository, sr repository.SessionRepository) AuthService {
	return &authService{
		ur: ur,
		sr: sr,
	}
}

type authService struct {
	ur repository.UserRepository
	sr repository.SessionRepository
}

func (as *authService) Register(ctx context.Context, req dto.RegisterRequest) error {
	// validate input
	email := req.Email
	if email == "" {
		return fmt.Errorf("email is empty")
	}

	// check email is exist
	existingUser, err := as.ur.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil { // there is a user with same email
		return fmt.Errorf("email already taken")
	}

	// hash pw
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		DateOfBirth:  req.DateOfBirth,
		Gender:       req.Gender,
		PasswordHash: string(hash),
		Role:         domain.Customer,
		IsActive:     true,
	}

	if err := as.ur.CreateUser(ctx, user); err != nil {
		return err
	}

	return nil
}

func (as *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Check email already registered or not
	existingUser, err := as.ur.GetUserByEmail(ctx, req.Email)
	if err != nil { // there is no email exist
		return nil, fmt.Errorf("email is not registered")
	}

	// compare the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("wrong password")
	}

	// generate random token
	token := make([]byte, 32)
	rand.Read(token)
	rawToken := base64.URLEncoding.EncodeToString(token)

	// hash it for storage
	hash := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hash[:])

	// Create session
	session := &domain.Session{
		UserID:    existingUser.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		IPAddress: req.IPAddress,
		UserAgent: req.UserAgent,
	}

	err = as.sr.CreateSession(ctx, session)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Session: session,
		RawToken: rawToken,
	}, nil
}

func (as *authService) Logout(ctx context.Context, sessionID string) error {
	// update the revoke time
	if err := as.sr.RevokeSession(ctx, sessionID); err != nil {
		return err
	}
	
	return nil
}
