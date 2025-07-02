package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aungmyozaw92/go-api-setup/internal/domain"
	"github.com/aungmyozaw92/go-api-setup/internal/usecase"
	"github.com/gorilla/mux"
)

// UserHandler handles user-related requests
type UserHandler struct {
	userUsecase usecase.UserUsecase
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// GetProfile returns the current user's profile
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from JWT context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		writeErrorResponse(w, "Invalid user context", http.StatusUnauthorized)
		return
	}

	user, err := h.userUsecase.GetProfile(r.Context(), userID)
	if err != nil {
		if err.Error() == "user not found" {
			writeErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		}
		writeErrorResponse(w, "Failed to get profile", http.StatusInternalServerError)
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "Profile retrieved successfully",
		"user":    user,
	}, http.StatusOK)
}

// CreateUser creates a new user (admin function)
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.userUsecase.CreateUser(r.Context(), &req)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			writeErrorResponse(w, err.Error(), http.StatusConflict)
			return
		}
		writeErrorResponse(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "User created successfully",
		"user":    user,
	}, http.StatusCreated)
}

// GetUser returns a specific user by ID
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	userIDStr, exists := vars["id"]
	if !exists {
		writeErrorResponse(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userUsecase.GetUserByID(r.Context(), uint(userID))
	if err != nil {
		if err.Error() == "user not found" {
			writeErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		}
		writeErrorResponse(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "User retrieved successfully",
		"user":    user,
	}, http.StatusOK)
}

// UpdateUser updates the current user's profile
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from JWT context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		writeErrorResponse(w, "Invalid user context", http.StatusUnauthorized)
		return
	}

	var req domain.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate password length if provided
	if req.Password != "" && len(req.Password) < 6 {
		writeErrorResponse(w, "Password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	user, err := h.userUsecase.UpdateUser(r.Context(), userID, &req)
	if err != nil {
		switch err.Error() {
		case "user not found":
			writeErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		case "email already exists":
			writeErrorResponse(w, err.Error(), http.StatusConflict)
			return
		default:
			writeErrorResponse(w, "Failed to update user", http.StatusInternalServerError)
			return
		}
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "User updated successfully",
		"user":    user,
	}, http.StatusOK)
}

// UpdateUserByID updates a specific user by ID (admin function)
func (h *UserHandler) UpdateUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	userIDStr, exists := vars["id"]
	if !exists {
		writeErrorResponse(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req domain.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate password length if provided
	if req.Password != "" && len(req.Password) < 6 {
		writeErrorResponse(w, "Password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	user, err := h.userUsecase.UpdateUser(r.Context(), uint(userID), &req)
	if err != nil {
		switch err.Error() {
		case "user not found":
			writeErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		case "email already exists":
			writeErrorResponse(w, err.Error(), http.StatusConflict)
			return
		default:
			writeErrorResponse(w, "Failed to update user", http.StatusInternalServerError)
			return
		}
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "User updated successfully",
		"user":    user,
	}, http.StatusOK)
}

// DeleteUser deletes the current user's account
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from JWT context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		writeErrorResponse(w, "Invalid user context", http.StatusUnauthorized)
		return
	}

	err := h.userUsecase.DeleteUser(r.Context(), userID)
	if err != nil {
		if err.Error() == "user not found" {
			writeErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		}
		writeErrorResponse(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "User deleted successfully",
	}, http.StatusOK)
}

// DeleteUserByID deletes a specific user by ID (admin function)
func (h *UserHandler) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	userIDStr, exists := vars["id"]
	if !exists {
		writeErrorResponse(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.userUsecase.DeleteUser(r.Context(), uint(userID))
	if err != nil {
		if err.Error() == "user not found" {
			writeErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		}
		writeErrorResponse(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "User deleted successfully",
	}, http.StatusOK)
}

// GetAllUsers returns all users with pagination
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters for pagination
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10 // default limit
	offset := 0 // default offset

	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	users, err := h.userUsecase.GetAllUsers(r.Context(), limit, offset)
	if err != nil {
		writeErrorResponse(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "Users retrieved successfully",
		"users":   users,
		"count":   len(users),
		"limit":   limit,
		"offset":  offset,
	}, http.StatusOK)
} 