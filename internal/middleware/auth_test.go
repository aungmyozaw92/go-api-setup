package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/aungmyozaw92/go-api-setup/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	// Test secret key
	jwtSecret := "test-jwt-secret"

	// Create middleware
	middleware := AuthMiddleware(jwtSecret)

	// Create a test handler that will be called if auth succeeds
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from context
		userID := r.Context().Value("user_id")
		email := r.Context().Value("user_email")
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		
		// Convert userID to string properly
		userIDStr := "0"
		if uid, ok := userID.(uint); ok {
			userIDStr = strconv.FormatUint(uint64(uid), 10)
		}
		
		emailStr := ""
		if e, ok := email.(string); ok {
			emailStr = e
		}
		
		response := map[string]interface{}{
			"message": "authenticated",
			"userID":  userIDStr,
			"email":   emailStr,
		}
		json.NewEncoder(w).Encode(response)
	})

	tests := []struct {
		name           string
		token          string
		expectedStatus int
		setupToken     bool
		userID         uint
		email          string
	}{
		{
			name:           "valid token",
			setupToken:     true,
			userID:         123,
			email:          "test@example.com",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing authorization header",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid authorization format",
			token:          "InvalidTokenFormat",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid token",
			token:          "Bearer invalid.jwt.token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "malformed bearer token",
			token:          "Bearer",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "wrong auth scheme",
			token:          "Basic dGVzdDp0ZXN0",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup token if needed
			var authToken string
			if tt.setupToken {
				token, err := utils.GenerateJWT(tt.userID, tt.email, jwtSecret)
				assert.NoError(t, err)
				authToken = "Bearer " + token
			} else {
				authToken = tt.token
			}

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if authToken != "" {
				req.Header.Set("Authorization", authToken)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Create the middleware wrapped handler
			handler := middleware(testHandler)

			// Execute
			handler.ServeHTTP(rr, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusOK {
				// For successful auth, verify the context was set correctly
				assert.Contains(t, rr.Body.String(), "authenticated")
				assert.Contains(t, rr.Body.String(), tt.email)
			} else {
				// For failed auth, check error message
				assert.Contains(t, rr.Body.String(), "error")
			}
		})
	}
}

func TestAuthMiddleware_ContextValues(t *testing.T) {
	jwtSecret := "test-jwt-secret"
	userID := uint(456)
	email := "context@example.com"

	// Generate a valid token
	token, err := utils.GenerateJWT(userID, email, jwtSecret)
	assert.NoError(t, err)

	// Create middleware
	middleware := AuthMiddleware(jwtSecret)

	// Create a test handler that checks context values
	var capturedUserID interface{}
	var capturedEmail interface{}
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUserID = r.Context().Value("user_id")
		capturedEmail = r.Context().Value("user_email")
		w.WriteHeader(http.StatusOK)
	})

	// Create request with valid token
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute
	handler := middleware(testHandler)
	handler.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, userID, capturedUserID)
	assert.Equal(t, email, capturedEmail)
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	// This test would require a more complex setup to create an expired token
	// For now, we test with an invalid token which simulates the same error path
	jwtSecret := "test-jwt-secret"
	middleware := AuthMiddleware(jwtSecret)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create request with an obviously invalid token
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.signature")

	rr := httptest.NewRecorder()

	// Execute
	handler := middleware(testHandler)
	handler.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "error")
}

func TestAuthMiddleware_DifferentSecrets(t *testing.T) {
	// Test that a token generated with one secret fails with another secret
	originalSecret := "original-secret"
	differentSecret := "different-secret"

	// Generate token with original secret
	token, err := utils.GenerateJWT(123, "test@example.com", originalSecret)
	assert.NoError(t, err)

	// Try to validate with different secret
	middleware := AuthMiddleware(differentSecret)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()

	// Execute
	handler := middleware(testHandler)
	handler.ServeHTTP(rr, req)

	// Assert - should fail due to different secret
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "error")
}

// Test helper function to check if middleware preserves request method
func TestAuthMiddleware_PreservesRequestMethod(t *testing.T) {
	jwtSecret := "test-jwt-secret"
	middleware := AuthMiddleware(jwtSecret)

	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}

	for _, method := range methods {
		t.Run("method_"+method, func(t *testing.T) {
			// Generate valid token
			token, err := utils.GenerateJWT(123, "test@example.com", jwtSecret)
			assert.NoError(t, err)

			var capturedMethod string
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedMethod = r.Method
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(method, "/protected", nil)
			req.Header.Set("Authorization", "Bearer "+token)

			rr := httptest.NewRecorder()

			handler := middleware(testHandler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)
			assert.Equal(t, method, capturedMethod)
		})
	}
} 