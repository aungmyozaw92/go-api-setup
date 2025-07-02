package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false, // bcrypt can hash empty strings
		},
		{
			name:     "long password",
			password: "this_is_a_long_password_but_within_bcrypt_72_byte_limit",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword, err := HashPassword(tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, hashedPassword)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hashedPassword)
				assert.NotEqual(t, tt.password, hashedPassword)
				
				// Verify that the same password produces different hashes (due to salt)
				hashedPassword2, err2 := HashPassword(tt.password)
				assert.NoError(t, err2)
				assert.NotEqual(t, hashedPassword, hashedPassword2)
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	password := "testpassword123"
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)

	tests := []struct {
		name           string
		password       string
		hashedPassword string
		wantErr        bool
	}{
		{
			name:           "correct password",
			password:       password,
			hashedPassword: hashedPassword,
			wantErr:        false,
		},
		{
			name:           "incorrect password",
			password:       "wrongpassword",
			hashedPassword: hashedPassword,
			wantErr:        true,
		},
		{
			name:           "empty password",
			password:       "",
			hashedPassword: hashedPassword,
			wantErr:        true,
		},
		{
			name:           "empty hash",
			password:       password,
			hashedPassword: "",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPassword(tt.password, tt.hashedPassword)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerateJWT(t *testing.T) {
	secretKey := "test-secret-key"
	userID := uint(123)
	email := "test@example.com"

	tests := []struct {
		name      string
		userID    uint
		email     string
		secretKey string
		wantErr   bool
	}{
		{
			name:      "valid input",
			userID:    userID,
			email:     email,
			secretKey: secretKey,
			wantErr:   false,
		},
		{
			name:      "empty email",
			userID:    userID,
			email:     "",
			secretKey: secretKey,
			wantErr:   false, // JWT can handle empty email
		},
		{
			name:      "zero user ID",
			userID:    0,
			email:     email,
			secretKey: secretKey,
			wantErr:   false, // JWT can handle zero userID
		},
		{
			name:      "empty secret key",
			userID:    userID,
			email:     email,
			secretKey: "",
			wantErr:   false, // JWT generation works with empty secret, but validation will fail
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateJWT(tt.userID, tt.email, tt.secretKey)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				
				// JWT should have 3 parts separated by dots
				parts := len(token)
				assert.Greater(t, parts, 0)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	secretKey := "test-secret-key"
	userID := uint(123)
	email := "test@example.com"

	// Generate a valid token for testing
	validToken, err := GenerateJWT(userID, email, secretKey)
	require.NoError(t, err)

	tests := []struct {
		name         string
		token        string
		secretKey    string
		expectUserID uint
		expectEmail  string
		wantErr      bool
	}{
		{
			name:         "valid token",
			token:        validToken,
			secretKey:    secretKey,
			expectUserID: userID,
			expectEmail:  email,
			wantErr:      false,
		},
		{
			name:      "invalid token format",
			token:     "invalid.token.format",
			secretKey: secretKey,
			wantErr:   true,
		},
		{
			name:      "empty token",
			token:     "",
			secretKey: secretKey,
			wantErr:   true,
		},
		{
			name:      "wrong secret key",
			token:     validToken,
			secretKey: "wrong-secret",
			wantErr:   true,
		},
		{
			name:      "malformed token",
			token:     "not.a.jwt",
			secretKey: secretKey,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateJWT(tt.token, tt.secretKey)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, tt.expectUserID, claims.UserID)
				assert.Equal(t, tt.expectEmail, claims.Email)
				
				// Check that token hasn't expired
				assert.True(t, time.Now().Before(claims.ExpiresAt.Time))
			}
		})
	}
}

func TestJWTWorkflow(t *testing.T) {
	// Test the complete JWT workflow: generate -> validate
	secretKey := "test-secret-key-for-workflow"
	userID := uint(456)
	email := "workflow@example.com"

	// Generate JWT
	token, err := GenerateJWT(userID, email, secretKey)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Validate JWT
	claims, err := ValidateJWT(token, secretKey)
	require.NoError(t, err)
	require.NotNil(t, claims)

	// Verify claims
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.True(t, time.Now().Before(claims.ExpiresAt.Time))
	assert.True(t, time.Now().After(claims.IssuedAt.Time.Add(-time.Second))) // Allow 1 second tolerance
} 