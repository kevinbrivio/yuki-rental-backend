package repository

import (
	"context"

	"github.com/kevinbrivio/yuki-rental-backend/internal/domain"
	"gorm.io/gorm"
)

type UserRepository interface { // Hosting database operations, not business logic
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository { // constructor
	return &userRepository{db: db}
}

func (ur *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	return ur.db.WithContext(ctx).Create(user).Error
}

func (ur *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := ur.db.WithContext(ctx).Where("email = ?", email).Find(&user).Error
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}