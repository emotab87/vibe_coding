package services

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/emotab87/vibe_coding/backend/internal/entities"
)

func TestJWTService_GenerateToken(t *testing.T) {
	service := NewJWTService("test-secret-key", 24)
	
	user := &entities.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}
	
	token, err := service.GenerateToken(user)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if token == "" {
		t.Fatal("Expected token to be generated, got empty string")
	}
}

func TestJWTService_ValidateToken(t *testing.T) {
	service := NewJWTService("test-secret-key", 24)
	
	user := &entities.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}
	
	// Generate token
	token, err := service.GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	// Validate token
	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if claims == nil {
		t.Fatal("Expected claims to be returned")
	}
	
	// Check user_id claim
	userIDClaim, exists := (*claims)["user_id"]
	if !exists {
		t.Fatal("Expected user_id claim to exist")
	}
	
	userID, ok := userIDClaim.(float64) // JSON numbers are float64
	if !ok {
		t.Fatalf("Expected user_id to be float64, got %T", userIDClaim)
	}
	
	if int64(userID) != user.ID {
		t.Fatalf("Expected user_id %d, got %d", user.ID, int64(userID))
	}
	
	// Check username claim
	usernameClaim, exists := (*claims)["username"]
	if !exists {
		t.Fatal("Expected username claim to exist")
	}
	
	username, ok := usernameClaim.(string)
	if !ok {
		t.Fatalf("Expected username to be string, got %T", usernameClaim)
	}
	
	if username != user.Username {
		t.Fatalf("Expected username %s, got %s", user.Username, username)
	}
}

func TestJWTService_ValidateToken_InvalidToken(t *testing.T) {
	service := NewJWTService("test-secret-key", 24)
	
	tests := []struct {
		name  string
		token string
	}{
		{"Empty token", ""},
		{"Invalid format", "invalid.token.format"},
		{"Wrong signature", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIn0.wrong_signature"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.ValidateToken(tt.token)
			if err == nil {
				t.Error("Expected error for invalid token, got nil")
			}
		})
	}
}

func TestJWTService_GetUserIDFromToken(t *testing.T) {
	service := NewJWTService("test-secret-key", 24)
	
	user := &entities.User{
		ID:       123,
		Username: "testuser",
		Email:    "test@example.com",
	}
	
	token, err := service.GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	userID, err := service.GetUserIDFromToken(token)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if userID != user.ID {
		t.Fatalf("Expected user ID %d, got %d", user.ID, userID)
	}
}

func TestJWTService_GetUsernameFromToken(t *testing.T) {
	service := NewJWTService("test-secret-key", 24)
	
	user := &entities.User{
		ID:       1,
		Username: "testuser123",
		Email:    "test@example.com",
	}
	
	token, err := service.GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	username, err := service.GetUsernameFromToken(token)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if username != user.Username {
		t.Fatalf("Expected username %s, got %s", user.Username, username)
	}
}

func TestJWTService_ParseToken(t *testing.T) {
	service := NewJWTService("test-secret-key", 24)
	
	user := &entities.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}
	
	tokenString, err := service.GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	token, err := service.ParseToken(tokenString)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if !token.Valid {
		t.Fatal("Expected token to be valid")
	}
	
	if token.Method != jwt.SigningMethodHS256 {
		t.Fatalf("Expected signing method HS256, got %v", token.Method)
	}
}

func TestJWTService_ExpiredToken(t *testing.T) {
	// Create service with very short expiry
	service := NewJWTService("test-secret-key", 0) // 0 hours = immediate expiry
	
	user := &entities.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}
	
	// Generate token (it will be expired immediately due to 0 hour expiry)
	token, err := service.GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	// Wait a moment to ensure expiration
	time.Sleep(1 * time.Millisecond)
	
	// Try to validate expired token
	_, err = service.ValidateToken(token)
	if err == nil {
		t.Fatal("Expected error for expired token, got nil")
	}
	
	if !IsTokenExpired(err) {
		t.Fatalf("Expected token expired error, got: %v", err)
	}
}

func TestJWTService_DifferentSecrets(t *testing.T) {
	service1 := NewJWTService("secret-key-1", 24)
	service2 := NewJWTService("secret-key-2", 24)
	
	user := &entities.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}
	
	// Generate token with first service
	token, err := service1.GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	// Try to validate with second service (different secret)
	_, err = service2.ValidateToken(token)
	if err == nil {
		t.Fatal("Expected error when validating token with different secret, got nil")
	}
	
	if !IsTokenInvalid(err) {
		t.Fatalf("Expected token invalid error, got: %v", err)
	}
}

func TestIsTokenExpired(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"Nil error", nil, false},
		{"Expired error", jwt.ErrTokenExpired, true},
		{"Other error", jwt.ErrSignatureInvalid, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsTokenExpired(tt.err)
			if result != tt.expected {
				t.Fatalf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsTokenInvalid(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"Nil error", nil, false},
		{"Signature invalid", jwt.ErrSignatureInvalid, true},
		{"Expired error", jwt.ErrTokenExpired, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsTokenInvalid(tt.err)
			if result != tt.expected {
				t.Fatalf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}