package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/emotab87/vibe_coding/backend/internal/entities"
	"github.com/emotab87/vibe_coding/backend/internal/middleware"
	"github.com/emotab87/vibe_coding/backend/internal/repositories"
	"github.com/emotab87/vibe_coding/backend/internal/services"
)

// Temporary stub handlers - to be implemented in future issues

// User authentication handlers
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "User registration not yet implemented")
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "User login not yet implemented")
}

func GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "Get current user not yet implemented")
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "Update user not yet implemented")
}

// Article handlers
func ListArticlesHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "List articles not yet implemented")
}

func GetArticleHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "Get article not yet implemented")
}

func CreateArticleHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "Create article not yet implemented")
}

func UpdateArticleHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "Update article not yet implemented")
}

func DeleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "Delete article not yet implemented")
}

// Comment handlers
func ListCommentsHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "List comments not yet implemented")
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "Create comment not yet implemented")
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "Delete comment not yet implemented")
}

// Profile handlers
func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "Get profile not yet implemented")
}

// Helper functions

// writeNotImplemented returns "not implemented" responses
func writeNotImplemented(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)

	response := map[string]string{
		"error":   "Not implemented",
		"message": message,
	}

	json.NewEncoder(w).Encode(response)
}

// writeJSON writes a JSON response
func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// writeError writes an error response
func writeError(w http.ResponseWriter, statusCode int, message string) {
	response := map[string]string{
		"error": message,
	}
	writeJSON(w, statusCode, response)
}

// writeValidationErrors writes validation error response
func writeValidationErrors(w http.ResponseWriter, validationErrors *entities.ValidationErrors) {
	response := map[string]interface{}{
		"errors": validationErrors.Errors,
	}
	writeJSON(w, http.StatusBadRequest, response)
}

// parseJSON parses JSON request body into the provided struct
func parseJSON(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	
	return nil
}

// getUserIDFromContext extracts user ID from request context
func getUserIDFromContext(r *http.Request) (int64, error) {
	userID := r.Context().Value(middleware.UserIDContextKey)
	if userID == nil {
		return 0, fmt.Errorf("user ID not found in context")
	}
	
	switch v := userID.(type) {
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case string:
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid user ID format: %w", err)
		}
		return id, nil
	default:
		return 0, fmt.Errorf("invalid user ID type: %T", userID)
	}
}

// getUsernameFromContext extracts username from request context
func getUsernameFromContext(r *http.Request) (string, error) {
	username := r.Context().Value(middleware.UsernameContextKey)
	if username == nil {
		return "", fmt.Errorf("username not found in context")
	}
	
	usernameStr, ok := username.(string)
	if !ok {
		return "", fmt.Errorf("invalid username type: %T", username)
	}
	
	return usernameStr, nil
}

// extractToken extracts JWT token from Authorization header
func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}
	
	if !strings.HasPrefix(authHeader, "Token ") {
		return "", fmt.Errorf("invalid authorization header format")
	}
	
	token := strings.TrimPrefix(authHeader, "Token ")
	if token == "" {
		return "", fmt.Errorf("missing token")
	}
	
	return token, nil
}