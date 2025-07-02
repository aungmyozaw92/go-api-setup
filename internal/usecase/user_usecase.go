package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/aungmyozaw92/go-api-setup/internal/domain"
	"github.com/aungmyozaw92/go-api-setup/internal/repository"
	"github.com/aungmyozaw92/go-api-setup/pkg/utils"
)

// UserUsecase defines the interface for user business logic
type UserUsecase interface {
	Register(ctx context.Context, req *domain.UserRequest) (*domain.UserResponse, error)
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error)
	CreateUser(ctx context.Context, req *domain.UserRequest) (*domain.UserResponse, error)
	GetProfile(ctx context.Context, userID uint) (*domain.UserResponse, error)
	GetUserByID(ctx context.Context, userID uint) (*domain.UserResponse, error)
	UpdateUser(ctx context.Context, userID uint, req *domain.UpdateUserRequest) (*domain.UserResponse, error)
	DeleteUser(ctx context.Context, userID uint) error
	GetAllUsers(ctx context.Context, limit, offset int) ([]*domain.UserResponse, error)
}

// userUsecase implements UserUsecase interface
type userUsecase struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

// NewUserUsecase creates a new user usecase
func NewUserUsecase(userRepo repository.UserRepository, jwtSecret string) UserUsecase {
	return &userUsecase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Register handles user registration
func (u *userUsecase) Register(ctx context.Context, req *domain.UserRequest) (*domain.UserResponse, error) {
	// Check if user already exists
	existingUser, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Return user response
	return &domain.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

// Login handles user authentication
func (u *userUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	// Get user by email
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Check password
	if err := utils.CheckPassword(req.Password, user.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Email, u.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Return login response
	return &domain.LoginResponse{
		Token: token,
		User: domain.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// CreateUser creates a new user (admin function)
func (u *userUsecase) CreateUser(ctx context.Context, req *domain.UserRequest) (*domain.UserResponse, error) {
	// Check if user already exists
	existingUser, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Return user response
	return &domain.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

// GetProfile gets the current user's profile
func (u *userUsecase) GetProfile(ctx context.Context, userID uint) (*domain.UserResponse, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return &domain.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

// GetUserByID gets a user by ID
func (u *userUsecase) GetUserByID(ctx context.Context, userID uint) (*domain.UserResponse, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return &domain.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

// UpdateUser updates a user's information
func (u *userUsecase) UpdateUser(ctx context.Context, userID uint, req *domain.UpdateUserRequest) (*domain.UserResponse, error) {
	// Get existing user
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Check if email is being changed and if it's already taken
	if req.Email != "" && req.Email != user.Email {
		existingUser, err := u.userRepo.GetByEmail(ctx, req.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing email: %w", err)
		}
		if existingUser != nil {
			return nil, errors.New("email already exists")
		}
		user.Email = req.Email
	}

	// Update fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}

	// Update password if provided
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = hashedPassword
	}

	// Save updated user
	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &domain.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

// DeleteUser deletes a user
func (u *userUsecase) DeleteUser(ctx context.Context, userID uint) error {
	// Check if user exists
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Delete user
	if err := u.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// GetAllUsers gets all users with pagination
func (u *userUsecase) GetAllUsers(ctx context.Context, limit, offset int) ([]*domain.UserResponse, error) {
	users, err := u.userRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	var userResponses []*domain.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, &domain.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		})
	}

	return userResponses, nil
} 