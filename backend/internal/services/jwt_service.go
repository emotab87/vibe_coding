package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/emotab87/vibe_coding/backend/internal/entities"
)

// JWTService handles JWT token operations
type JWTService interface {
	GenerateToken(user *entities.User) (string, error)
	ValidateToken(tokenString string) (*jwt.MapClaims, error)
	ParseToken(tokenString string) (*jwt.Token, error)
	GetUserIDFromToken(tokenString string) (int64, error)
	GetUsernameFromToken(tokenString string) (string, error)
}

// jwtService implements JWTService
type jwtService struct {
	secretKey    []byte
	tokenExpiry  time.Duration
	signingMethod jwt.SigningMethod
}

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// NewJWTService creates a new JWT service
func NewJWTService(secretKey string, tokenExpiryHours int) JWTService {
	return &jwtService{
		secretKey:     []byte(secretKey),
		tokenExpiry:   time.Duration(tokenExpiryHours) * time.Hour,
		signingMethod: jwt.SigningMethodHS256,
	}
}

// GenerateToken generates a JWT token for a user
func (s *jwtService) GenerateToken(user *entities.User) (string, error) {
	now := time.Now()
	expirationTime := now.Add(s.tokenExpiry)

	claims := &JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   fmt.Sprintf("user:%d", user.ID),
			Issuer:    "conduit-api",
		},
	}

	token := jwt.NewWithClaims(s.signingMethod, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *jwtService) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := s.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return &claims, nil
}

// ParseToken parses a JWT token
func (s *jwtService) ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		} else if method != s.signingMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", method)
		}

		return s.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	return token, nil
}

// GetUserIDFromToken extracts user ID from token
func (s *jwtService) GetUserIDFromToken(tokenString string) (int64, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}

	userIDClaim, exists := (*claims)["user_id"]
	if !exists {
		return 0, fmt.Errorf("user_id not found in token")
	}

	// Handle both int64 and float64 (JSON number conversion)
	switch v := userIDClaim.(type) {
	case float64:
		return int64(v), nil
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("user_id has invalid type: %T", userIDClaim)
	}
}

// GetUsernameFromToken extracts username from token
func (s *jwtService) GetUsernameFromToken(tokenString string) (string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	usernameClaim, exists := (*claims)["username"]
	if !exists {
		return "", fmt.Errorf("username not found in token")
	}

	username, ok := usernameClaim.(string)
	if !ok {
		return "", fmt.Errorf("username has invalid type: %T", usernameClaim)
	}

	return username, nil
}

// IsTokenExpired checks if a token is expired
func IsTokenExpired(err error) bool {
	if err == nil {
		return false
	}

	// Check for expiration error
	return containsString(err.Error(), "token is expired") ||
		   containsString(err.Error(), "exp")
}

// IsTokenInvalid checks if a token is invalid (malformed, wrong signature, etc.)
func IsTokenInvalid(err error) bool {
	if err == nil {
		return false
	}

	// Check for various token validation errors
	return containsString(err.Error(), "token is malformed") ||
		   containsString(err.Error(), "signature is invalid") ||
		   containsString(err.Error(), "unexpected signing method") ||
		   containsString(err.Error(), "invalid token")
}

// Helper function to check string contains (case-insensitive)
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(toLowerCase(s), toLowerCase(substr)) >= 0
}

// Helper function to convert to lowercase
func toLowerCase(s string) string {
	result := make([]rune, len(s))
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			result[i] = r + 32
		} else {
			result[i] = r
		}
	}
	return string(result)
}

// Helper function to find substring
func findSubstring(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}