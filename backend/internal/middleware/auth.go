package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// ContextKey type for context keys
type ContextKey string

const (
	// UserIDContextKey is the key for user ID in context
	UserIDContextKey ContextKey = "user_id"
	// UsernameContextKey is the key for username in context
	UsernameContextKey ContextKey = "username"
)

// AuthMiddleware validates JWT tokens and adds user info to context
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeUnauthorizedError(w, "Missing authorization header")
				return
			}

			// Check if it starts with "Token "
			if !strings.HasPrefix(authHeader, "Token ") {
				writeUnauthorizedError(w, "Invalid authorization header format")
				return
			}

			// Extract the token
			tokenString := strings.TrimPrefix(authHeader, "Token ")
			if tokenString == "" {
				writeUnauthorizedError(w, "Missing token")
				return
			}

			// Parse and validate the token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Validate the signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			})

			if err != nil {
				writeUnauthorizedError(w, "Invalid token")
				return
			}

			if !token.Valid {
				writeUnauthorizedError(w, "Token is not valid")
				return
			}

			// Extract claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				writeUnauthorizedError(w, "Invalid token claims")
				return
			}

			// Get user info from claims
			userID, ok := claims["user_id"]
			if !ok {
				writeUnauthorizedError(w, "Missing user_id in token")
				return
			}

			username, ok := claims["username"]
			if !ok {
				writeUnauthorizedError(w, "Missing username in token")
				return
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
			ctx = context.WithValue(ctx, UsernameContextKey, username)

			// Call next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// writeUnauthorizedError writes a 401 Unauthorized response
func writeUnauthorizedError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	response := ErrorResponse{
		Error: message,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// If JSON encoding fails, fall back to plain text
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Unauthorized"))
	}
}