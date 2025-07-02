package mocks

import (
	"context"

	"github.com/aungmyozaw92/go-api-setup/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// GetByEmail mocks the GetByEmail method
func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// GetByID mocks the GetByID method
func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// Update mocks the Update method
func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// GetAll mocks the GetAll method
func (m *MockUserRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
} 