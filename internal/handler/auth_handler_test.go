package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aungmyozaw92/go-api-setup/internal/domain"
	"github.com/aungmyozaw92/go-api-setup/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AuthHandlerTestSuite struct {
	suite.Suite
	mockUsecase *mocks.MockUserUsecase
	handler     *AuthHandler
}

func (suite *AuthHandlerTestSuite) SetupTest() {
	suite.mockUsecase = new(mocks.MockUserUsecase)
	suite.handler = NewAuthHandler(suite.mockUsecase)
}

func (suite *AuthHandlerTestSuite) TearDownTest() {
	suite.mockUsecase.AssertExpectations(suite.T())
}

// Test Register Handler
func (suite *AuthHandlerTestSuite) TestRegister_Success() {
	reqBody := &domain.UserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	expectedResponse := &domain.UserResponse{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
	}

	// Setup mock
	suite.mockUsecase.On("Register", mock.Anything, reqBody).Return(expectedResponse, nil)

	// Create request
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.Register(rr, req)

	// Assert
	assert.Equal(suite.T(), http.StatusCreated, rr.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), "User registered successfully", response["message"])
	assert.NotNil(suite.T(), response["user"])
	
	userData := response["user"].(map[string]interface{})
	assert.Equal(suite.T(), float64(1), userData["id"]) // JSON unmarshals numbers as float64
	assert.Equal(suite.T(), "John Doe", userData["name"])
	assert.Equal(suite.T(), "john@example.com", userData["email"])
}

func (suite *AuthHandlerTestSuite) TestRegister_InvalidJSON() {
	// Create request with invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	
	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.Register(rr, req)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, rr.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), "Invalid request body", response["error"])
}

func (suite *AuthHandlerTestSuite) TestRegister_ValidationError() {
	reqBody := &domain.UserRequest{
		Name:     "", // Empty name should cause validation error
		Email:    "john@example.com",
		Password: "password123",
	}

	// Create request
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.Register(rr, req)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, rr.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Contains(suite.T(), response["error"], "Name, email, and password are required")
}

func (suite *AuthHandlerTestSuite) TestRegister_UsecaseError() {
	reqBody := &domain.UserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	// Setup mock to return error
	suite.mockUsecase.On("Register", mock.Anything, reqBody).Return(nil, errors.New("user with this email already exists"))

	// Create request
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.Register(rr, req)

	// Assert
	assert.Equal(suite.T(), http.StatusConflict, rr.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), "user with this email already exists", response["error"])
}

// Test Login Handler
func (suite *AuthHandlerTestSuite) TestLogin_Success() {
	reqBody := &domain.LoginRequest{
		Email:    "john@example.com",
		Password: "password123",
	}

	expectedResponse := &domain.LoginResponse{
		Token: "jwt-token-here",
		User: domain.UserResponse{
			ID:    1,
			Name:  "John Doe",
			Email: "john@example.com",
		},
	}

	// Setup mock
	suite.mockUsecase.On("Login", mock.Anything, reqBody).Return(expectedResponse, nil)

	// Create request
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.Login(rr, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), "Login successful", response["message"])
	assert.Equal(suite.T(), "jwt-token-here", response["token"])
	assert.NotNil(suite.T(), response["user"])
}

func (suite *AuthHandlerTestSuite) TestLogin_InvalidJSON() {
	// Create request with invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	
	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.Login(rr, req)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, rr.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), "Invalid request body", response["error"])
}

func (suite *AuthHandlerTestSuite) TestLogin_ValidationError() {
	reqBody := &domain.LoginRequest{
		Email:    "", // Empty email should cause validation error
		Password: "password123",
	}

	// Create request
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.Login(rr, req)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, rr.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Contains(suite.T(), response["error"], "Email and password are required")
}

func (suite *AuthHandlerTestSuite) TestLogin_InvalidCredentials() {
	reqBody := &domain.LoginRequest{
		Email:    "john@example.com",
		Password: "wrongpassword",
	}

	// Setup mock to return error
	suite.mockUsecase.On("Login", mock.Anything, reqBody).Return(nil, errors.New("invalid email or password"))

	// Create request
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute
	suite.handler.Login(rr, req)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, rr.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), "invalid email or password", response["error"])
}



// Run the test suite
func TestAuthHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerTestSuite))
} 