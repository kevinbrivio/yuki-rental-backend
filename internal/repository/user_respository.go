package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kevinbrivio/yuki-rental-backend/internal/domain"
)

type UserRepository interface { // Hosting database operations, not business logic
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository { // constructor
	return &userRepository{db: db}
}

func (ur *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	
}

func (ur *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	
}