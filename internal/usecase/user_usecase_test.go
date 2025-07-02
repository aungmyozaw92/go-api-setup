package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aungmyozaw92/go-api-setup/internal/domain"
	"github.com/aungmyozaw92/go-api-setup/internal/repository/mocks"
	"github.com/aungmyozaw92/go-api-setup/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserUsecaseTestSuite struct {
	suite.Suite
	mockRepo    *mocks.MockUserRepository
	usecase     UserUsecase
	ctx         context.Context
	jwtSecret   string
}

func (suite *UserUsecaseTestSuite) SetupTest() {
	suite.mockRepo = new(mocks.MockUserRepository)
	suite.jwtSecret = "test-jwt-secret"
	suite.usecase = NewUserUsecase(suite.mockRepo, suite.jwtSecret)
	suite.ctx = context.Background()
}

func (suite *UserUsecaseTestSuite) TearDownTest() {
	suite.mockRepo.AssertExpectations(suite.T())
}

// Test Register
func (suite *UserUsecaseTestSuite) TestRegister_Success() {
	req := &domain.UserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	// Mock expectations
	suite.mockRepo.On("GetByEmail", suite.ctx, req.Email).Return(nil, nil)
	suite.mockRepo.On("Create", suite.ctx, mock.AnythingOfType("*domain.User")).Return(nil).Run(func(args mock.Arguments) {
		user := args.Get(1).(*domain.User)
		user.ID = 1
		user.CreatedAt = time.Now()
	})

	// Execute
	result, err := suite.usecase.Register(suite.ctx, req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), uint(1), result.ID)
	assert.Equal(suite.T(), req.Name, result.Name)
	assert.Equal(suite.T(), req.Email, result.Email)
}

func (suite *UserUsecaseTestSuite) TestRegister_EmailAlreadyExists() {
	req := &domain.UserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	existingUser := &domain.User{
		ID:    1,
		Email: req.Email,
	}

	// Mock expectations
	suite.mockRepo.On("GetByEmail", suite.ctx, req.Email).Return(existingUser, nil)

	// Execute
	result, err := suite.usecase.Register(suite.ctx, req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "user with this email already exists")
}

func (suite *UserUsecaseTestSuite) TestRegister_RepositoryError() {
	req := &domain.UserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	// Mock expectations
	suite.mockRepo.On("GetByEmail", suite.ctx, req.Email).Return(nil, errors.New("database error"))

	// Execute
	result, err := suite.usecase.Register(suite.ctx, req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "failed to check existing user")
}

// Test Login
func (suite *UserUsecaseTestSuite) TestLogin_Success() {
	req := &domain.LoginRequest{
		Email:    "john@example.com",
		Password: "password123",
	}

	// Create a user with hashed password - generate proper hash
	hashedPassword, err := utils.HashPassword(req.Password)
	suite.NoError(err)
	
	user := &domain.User{
		ID:       1,
		Name:     "John Doe",
		Email:    req.Email,
		Password: hashedPassword,
	}

	// Mock expectations
	suite.mockRepo.On("GetByEmail", suite.ctx, req.Email).Return(user, nil)

	// Execute
	result, err := suite.usecase.Login(suite.ctx, req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.NotEmpty(suite.T(), result.Token)
	assert.Equal(suite.T(), user.ID, result.User.ID)
	assert.Equal(suite.T(), user.Name, result.User.Name)
	assert.Equal(suite.T(), user.Email, result.User.Email)
}

func (suite *UserUsecaseTestSuite) TestLogin_UserNotFound() {
	req := &domain.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	// Mock expectations
	suite.mockRepo.On("GetByEmail", suite.ctx, req.Email).Return(nil, nil)

	// Execute
	result, err := suite.usecase.Login(suite.ctx, req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "invalid email or password")
}

func (suite *UserUsecaseTestSuite) TestLogin_WrongPassword() {
	req := &domain.LoginRequest{
		Email:    "john@example.com",
		Password: "wrongpassword",
	}

	// Create a user with hashed password - hash the correct password, not the wrong one
	hashedPassword, err := utils.HashPassword("password123")
	suite.NoError(err)
	
	user := &domain.User{
		ID:       1,
		Email:    req.Email,
		Password: hashedPassword,
	}

	// Mock expectations
	suite.mockRepo.On("GetByEmail", suite.ctx, req.Email).Return(user, nil)

	// Execute
	result, err := suite.usecase.Login(suite.ctx, req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "invalid email or password")
}

// Test GetProfile
func (suite *UserUsecaseTestSuite) TestGetProfile_Success() {
	userID := uint(1)
	user := &domain.User{
		ID:        userID,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
	}

	// Mock expectations
	suite.mockRepo.On("GetByID", suite.ctx, userID).Return(user, nil)

	// Execute
	result, err := suite.usecase.GetProfile(suite.ctx, userID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), user.ID, result.ID)
	assert.Equal(suite.T(), user.Name, result.Name)
	assert.Equal(suite.T(), user.Email, result.Email)
}

func (suite *UserUsecaseTestSuite) TestGetProfile_UserNotFound() {
	userID := uint(999)

	// Mock expectations
	suite.mockRepo.On("GetByID", suite.ctx, userID).Return(nil, nil)

	// Execute
	result, err := suite.usecase.GetProfile(suite.ctx, userID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "user not found")
}

// Test UpdateUser
func (suite *UserUsecaseTestSuite) TestUpdateUser_Success() {
	userID := uint(1)
	updateReq := &domain.UpdateUserRequest{
		Name:  "John Updated",
		Email: "john.updated@example.com",
	}

	existingUser := &domain.User{
		ID:    userID,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Mock expectations
	suite.mockRepo.On("GetByID", suite.ctx, userID).Return(existingUser, nil)
	suite.mockRepo.On("GetByEmail", suite.ctx, updateReq.Email).Return(nil, nil)
	suite.mockRepo.On("Update", suite.ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	// Execute
	result, err := suite.usecase.UpdateUser(suite.ctx, userID, updateReq)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), updateReq.Name, result.Name)
	assert.Equal(suite.T(), updateReq.Email, result.Email)
}

func (suite *UserUsecaseTestSuite) TestUpdateUser_EmailAlreadyExists() {
	userID := uint(1)
	updateReq := &domain.UpdateUserRequest{
		Email: "existing@example.com",
	}

	existingUser := &domain.User{
		ID:    userID,
		Email: "john@example.com",
	}

	conflictUser := &domain.User{
		ID:    2,
		Email: updateReq.Email,
	}

	// Mock expectations
	suite.mockRepo.On("GetByID", suite.ctx, userID).Return(existingUser, nil)
	suite.mockRepo.On("GetByEmail", suite.ctx, updateReq.Email).Return(conflictUser, nil)

	// Execute
	result, err := suite.usecase.UpdateUser(suite.ctx, userID, updateReq)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "email already exists")
}

// Test DeleteUser
func (suite *UserUsecaseTestSuite) TestDeleteUser_Success() {
	userID := uint(1)
	user := &domain.User{
		ID:    userID,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Mock expectations
	suite.mockRepo.On("GetByID", suite.ctx, userID).Return(user, nil)
	suite.mockRepo.On("Delete", suite.ctx, userID).Return(nil)

	// Execute
	err := suite.usecase.DeleteUser(suite.ctx, userID)

	// Assert
	assert.NoError(suite.T(), err)
}

func (suite *UserUsecaseTestSuite) TestDeleteUser_UserNotFound() {
	userID := uint(999)

	// Mock expectations
	suite.mockRepo.On("GetByID", suite.ctx, userID).Return(nil, nil)

	// Execute
	err := suite.usecase.DeleteUser(suite.ctx, userID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "user not found")
}

// Test GetAllUsers
func (suite *UserUsecaseTestSuite) TestGetAllUsers_Success() {
	limit := 10
	offset := 0
	users := []*domain.User{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
	}

	// Mock expectations
	suite.mockRepo.On("GetAll", suite.ctx, limit, offset).Return(users, nil)

	// Execute
	result, err := suite.usecase.GetAllUsers(suite.ctx, limit, offset)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), users[0].ID, result[0].ID)
	assert.Equal(suite.T(), users[1].ID, result[1].ID)
}

func (suite *UserUsecaseTestSuite) TestGetAllUsers_RepositoryError() {
	limit := 10
	offset := 0

	// Mock expectations
	suite.mockRepo.On("GetAll", suite.ctx, limit, offset).Return(nil, errors.New("database error"))

	// Execute
	result, err := suite.usecase.GetAllUsers(suite.ctx, limit, offset)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "failed to get users")
}

// Run the test suite
func TestUserUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
} 