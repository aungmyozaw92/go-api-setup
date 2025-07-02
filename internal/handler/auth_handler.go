package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aungmyozaw92/go-api-setup/internal/domain"
	"github.com/aungmyozaw92/go-api-setup/internal/usecase"
)

// AuthHandler handles authentication related requests
type AuthHandler struct {
	userUsecase usecase.UserUsecase
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userUsecase usecase.UserUsecase) *AuthHandler {
	return &AuthHandler{
		userUsecase: userUsecase,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Name == "" || req.Email == "" || req.Password == "" {
		writeErrorResponse(w, "Name, email, and password are required", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 6 {
		writeErrorResponse(w, "Password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	user, err := h.userUsecase.Register(r.Context(), &req)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			writeErrorResponse(w, err.Error(), http.StatusConflict)
			return
		}
		writeErrorResponse(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "User registered successfully",
		"user":    user,
	}, http.StatusCreated)
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Email == "" || req.Password == "" {
		writeErrorResponse(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	loginResponse, err := h.userUsecase.Login(r.Context(), &req)
	if err != nil {
		if err.Error() == "invalid email or password" {
			writeErrorResponse(w, err.Error(), http.StatusUnauthorized)
			return
		}
		writeErrorResponse(w, "Login failed", http.StatusInternalServerError)
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "Login successful",
		"token":   loginResponse.Token,
		"user":    loginResponse.User,
	}, http.StatusOK)
}

// writeErrorResponse writes an error response in JSON format
func writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := map[string]string{
		"error": message,
	}
	
	json.NewEncoder(w).Encode(response)
}

// writeSuccessResponse writes a success response in JSON format
func writeSuccessResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
} 