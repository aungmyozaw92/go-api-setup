package repository

import (
	"context"
	"errors"

	"github.com/aungmyozaw92/go-api-setup/internal/domain"
	"gorm.io/gorm"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uint) error
	GetAll(ctx context.Context, limit, offset int) ([]*domain.User, error)
}

// userRepository implements UserRepository interface
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create creates a new user in the database
func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil instead of error for not found
		}
		return nil, err
	}
	return &user, nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil instead of error for not found
		}
		return nil, err
	}
	return &user, nil
}

// Update updates a user in the database
func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return err
	}
	return nil
}

// Delete deletes a user from the database (soft delete)
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&domain.User{}, id).Error; err != nil {
		return err
	}
	return nil
}

// GetAll retrieves all users with pagination
func (r *userRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	query := r.db.WithContext(ctx).Model(&domain.User{})
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
} 