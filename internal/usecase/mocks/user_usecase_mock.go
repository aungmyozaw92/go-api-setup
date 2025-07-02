package mocks

import (
	"context"

	"github.com/aungmyozaw92/go-api-setup/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockUserUsecase is a mock implementation of UserUsecaseInterface
type MockUserUsecase struct {
	mock.Mock
}

// Register mocks the Register method
func (m *MockUserUsecase) Register(ctx context.Context, req *domain.UserRequest) (*domain.UserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserResponse), args.Error(1)
}

// Login mocks the Login method
func (m *MockUserUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LoginResponse), args.Error(1)
}

// GetProfile mocks the GetProfile method
func (m *MockUserUsecase) GetProfile(ctx context.Context, userID uint) (*domain.UserResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserResponse), args.Error(1)
}

// UpdateUser mocks the UpdateUser method
func (m *MockUserUsecase) UpdateUser(ctx context.Context, userID uint, req *domain.UpdateUserRequest) (*domain.UserResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserResponse), args.Error(1)
}

// DeleteUser mocks the DeleteUser method
func (m *MockUserUsecase) DeleteUser(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// CreateUser mocks the CreateUser method
func (m *MockUserUsecase) CreateUser(ctx context.Context, req *domain.UserRequest) (*domain.UserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserResponse), args.Error(1)
}

// GetAllUsers mocks the GetAllUsers method
func (m *MockUserUsecase) GetAllUsers(ctx context.Context, limit, offset int) ([]*domain.UserResponse, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.UserResponse), args.Error(1)
}

// GetUserByID mocks the GetUserByID method
func (m *MockUserUsecase) GetUserByID(ctx context.Context, userID uint) (*domain.UserResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserResponse), args.Error(1)
} 