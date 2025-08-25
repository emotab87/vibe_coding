package handlers

import (
	"net/http"

	"github.com/emotab87/vibe_coding/backend/internal/entities"
	"github.com/emotab87/vibe_coding/backend/internal/repositories"
	"github.com/emotab87/vibe_coding/backend/internal/services"
)

// AuthHandlers handles authentication-related HTTP requests
type AuthHandlers struct {
	userRepo   repositories.UserRepository
	jwtService services.JWTService
}

// NewAuthHandlers creates a new auth handlers instance
func NewAuthHandlers(userRepo repositories.UserRepository, jwtService services.JWTService) *AuthHandlers {
	return &AuthHandlers{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// RegisterUser handles user registration
func (h *AuthHandlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse request body
	var req struct {
		User entities.UserRegistration `json:"user"`
	}

	if err := parseJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Validate user data
	if validationErr := req.User.Validate(); validationErr != nil {
		writeValidationErrors(w, validationErr)
		return
	}

	// Check if email already exists
	if exists, err := h.userRepo.EmailExists(req.User.Email); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	} else if exists {
		writeError(w, http.StatusBadRequest, "User with this email already exists")
		return
	}

	// Check if username already exists
	if exists, err := h.userRepo.UsernameExists(req.User.Username); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	} else if exists {
		writeError(w, http.StatusBadRequest, "User with this username already exists")
		return
	}

	// Create user
	user, err := h.userRepo.Create(&req.User)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate JWT token
	token, err := h.jwtService.GenerateToken(user)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Return user response
	response := user.ToUserResponse(token)
	writeJSON(w, http.StatusCreated, response)
}

// LoginUser handles user login
func (h *AuthHandlers) LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse request body
	var req struct {
		User entities.UserLogin `json:"user"`
	}

	if err := parseJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Validate login data
	if validationErr := req.User.Validate(); validationErr != nil {
		writeValidationErrors(w, validationErr)
		return
	}

	// Get user by email
	user, err := h.userRepo.GetByEmail(req.User.Email)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Verify password
	if !h.userRepo.VerifyPassword(user, req.User.Password) {
		writeError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate JWT token
	token, err := h.jwtService.GenerateToken(user)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Return user response
	response := user.ToUserResponse(token)
	writeJSON(w, http.StatusOK, response)
}

// GetCurrentUser handles getting current user info
func (h *AuthHandlers) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, err := getUserIDFromContext(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get user from database
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		writeError(w, http.StatusNotFound, "User not found")
		return
	}

	// Extract token from request header
	token, err := extractToken(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	// Return user response with current token
	response := user.ToUserResponse(token)
	writeJSON(w, http.StatusOK, response)
}

// UpdateUser handles updating current user info
func (h *AuthHandlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, err := getUserIDFromContext(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse request body
	var req struct {
		User entities.UserUpdate `json:"user"`
	}

	if err := parseJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Validate update data
	if validationErr := req.User.Validate(); validationErr != nil {
		writeValidationErrors(w, validationErr)
		return
	}

	// Check email uniqueness if email is being updated
	if req.User.Email != nil {
		if exists, err := h.userRepo.EmailExists(*req.User.Email); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		} else if exists {
			// Check if it's not the current user's email
			currentUser, err := h.userRepo.GetByID(userID)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "Internal server error")
				return
			}
			if currentUser.Email != *req.User.Email {
				writeError(w, http.StatusBadRequest, "Email already exists")
				return
			}
		}
	}

	// Check username uniqueness if username is being updated
	if req.User.Username != nil {
		if exists, err := h.userRepo.UsernameExists(*req.User.Username); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		} else if exists {
			// Check if it's not the current user's username
			currentUser, err := h.userRepo.GetByID(userID)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "Internal server error")
				return
			}
			if currentUser.Username != *req.User.Username {
				writeError(w, http.StatusBadRequest, "Username already exists")
				return
			}
		}
	}

	// Update user
	updatedUser, err := h.userRepo.Update(userID, &req.User)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	// Generate new JWT token (in case username changed)
	token, err := h.jwtService.GenerateToken(updatedUser)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Return updated user response
	response := updatedUser.ToUserResponse(token)
	writeJSON(w, http.StatusOK, response)
}