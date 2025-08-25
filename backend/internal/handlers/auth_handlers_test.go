package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/emotab87/vibe_coding/backend/internal/database"
	"github.com/emotab87/vibe_coding/backend/internal/entities"
	"github.com/emotab87/vibe_coding/backend/internal/middleware"
	"github.com/emotab87/vibe_coding/backend/internal/repositories"
	"github.com/emotab87/vibe_coding/backend/internal/services"
)

func setupTestDB(t *testing.T) *database.DB {
	// Create temporary database file
	tempDir, err := os.MkdirTemp("", "auth_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	
	dbPath := filepath.Join(tempDir, "test.db")
	
	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	
	// Run migrations
	if err := db.Migrate("../../../migrations"); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
	
	return db
}

func setupTestHandlers(t *testing.T) (*AuthHandlers, *database.DB) {
	db := setupTestDB(t)
	userRepo := repositories.NewUserRepository(db)
	jwtService := services.NewJWTService("test-secret-key", 24)
	handlers := NewAuthHandlers(userRepo, jwtService)
	
	return handlers, db
}

func cleanupTestDB(db *database.DB) {
	if db != nil {
		db.Close()
	}
}

func TestAuthHandlers_RegisterUser(t *testing.T) {
	handlers, db := setupTestHandlers(t)
	defer cleanupTestDB(db)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectUser     bool
	}{
		{
			name: "Valid registration",
			requestBody: map[string]interface{}{
				"user": map[string]interface{}{
					"username": "testuser",
					"email":    "test@example.com",
					"password": "password123",
				},
			},
			expectedStatus: http.StatusCreated,
			expectUser:     true,
		},
		{
			name: "Invalid JSON format",
			requestBody: map[string]interface{}{
				"invalid": "format",
			},
			expectedStatus: http.StatusBadRequest,
			expectUser:     false,
		},
		{
			name: "Missing required fields",
			requestBody: map[string]interface{}{
				"user": map[string]interface{}{
					"username": "testuser",
					// missing email and password
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectUser:     false,
		},
		{
			name: "Invalid email format",
			requestBody: map[string]interface{}{
				"user": map[string]interface{}{
					"username": "testuser",
					"email":    "invalid-email",
					"password": "password123",
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectUser:     false,
		},
		{
			name: "Password too short",
			requestBody: map[string]interface{}{
				"user": map[string]interface{}{
					"username": "testuser",
					"email":    "test@example.com",
					"password": "123",
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectUser:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Call handler
			handlers.RegisterUser(w, req)
			
			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			
			// For successful registration, verify response structure
			if tt.expectUser && w.Code == http.StatusCreated {
				var response map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				
				user, ok := response["user"].(map[string]interface{})
				if !ok {
					t.Fatal("Response should contain user object")
				}
				
				if _, ok := user["username"].(string); !ok {
					t.Error("User should have username")
				}
				
				if _, ok := user["email"].(string); !ok {
					t.Error("User should have email")
				}
				
				if _, ok := user["token"].(string); !ok {
					t.Error("User should have token")
				}
			}
		})
	}
}

func TestAuthHandlers_LoginUser(t *testing.T) {
	handlers, db := setupTestHandlers(t)
	defer cleanupTestDB(db)

	// First, register a test user
	registerBody := map[string]interface{}{
		"user": map[string]interface{}{
			"username": "testuser",
			"email":    "test@example.com",
			"password": "password123",
		},
	}
	
	body, _ := json.Marshal(registerBody)
	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handlers.RegisterUser(w, req)
	
	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to register test user: %d", w.Code)
	}

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectUser     bool
	}{
		{
			name: "Valid login",
			requestBody: map[string]interface{}{
				"user": map[string]interface{}{
					"email":    "test@example.com",
					"password": "password123",
				},
			},
			expectedStatus: http.StatusOK,
			expectUser:     true,
		},
		{
			name: "Invalid email",
			requestBody: map[string]interface{}{
				"user": map[string]interface{}{
					"email":    "nonexistent@example.com",
					"password": "password123",
				},
			},
			expectedStatus: http.StatusUnauthorized,
			expectUser:     false,
		},
		{
			name: "Invalid password",
			requestBody: map[string]interface{}{
				"user": map[string]interface{}{
					"email":    "test@example.com",
					"password": "wrongpassword",
				},
			},
			expectedStatus: http.StatusUnauthorized,
			expectUser:     false,
		},
		{
			name: "Missing email",
			requestBody: map[string]interface{}{
				"user": map[string]interface{}{
					"password": "password123",
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectUser:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Call handler
			handlers.LoginUser(w, req)
			
			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			
			// For successful login, verify response structure
			if tt.expectUser && w.Code == http.StatusOK {
				var response map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				
				user, ok := response["user"].(map[string]interface{})
				if !ok {
					t.Fatal("Response should contain user object")
				}
				
				if token, ok := user["token"].(string); !ok || token == "" {
					t.Error("User should have non-empty token")
				}
			}
		})
	}
}

func TestAuthHandlers_GetCurrentUser(t *testing.T) {
	handlers, db := setupTestHandlers(t)
	defer cleanupTestDB(db)

	// Register and get token for test user
	registerBody := map[string]interface{}{
		"user": map[string]interface{}{
			"username": "testuser",
			"email":    "test@example.com",
			"password": "password123",
		},
	}
	
	body, _ := json.Marshal(registerBody)
	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handlers.RegisterUser(w, req)
	
	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to register test user: %d", w.Code)
	}
	
	var registerResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &registerResponse)
	user := registerResponse["user"].(map[string]interface{})
	token := user["token"].(string)

	tests := []struct {
		name           string
		token          string
		expectedStatus int
		setContext     bool
	}{
		{
			name:           "Valid token with context",
			token:          token,
			expectedStatus: http.StatusOK,
			setContext:     true,
		},
		{
			name:           "Missing context",
			token:          token,
			expectedStatus: http.StatusUnauthorized,
			setContext:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest(http.MethodGet, "/api/user", nil)
			req.Header.Set("Authorization", "Token "+tt.token)
			
			// Set context if needed (normally done by middleware)
			if tt.setContext {
				ctx := req.Context()
				ctx = context.WithValue(ctx, middleware.UserIDContextKey, int64(1))
				ctx = context.WithValue(ctx, middleware.UsernameContextKey, "testuser")
				req = req.WithContext(ctx)
			}
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Call handler
			handlers.GetCurrentUser(w, req)
			
			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			
			// For successful request, verify response structure
			if tt.setContext && w.Code == http.StatusOK {
				var response map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				
				user, ok := response["user"].(map[string]interface{})
				if !ok {
					t.Fatal("Response should contain user object")
				}
				
				if username, ok := user["username"].(string); !ok || username != "testuser" {
					t.Error("User should have correct username")
				}
			}
		})
	}
}

func TestAuthHandlers_UpdateUser(t *testing.T) {
	handlers, db := setupTestHandlers(t)
	defer cleanupTestDB(db)

	// Register test user first
	registerBody := map[string]interface{}{
		"user": map[string]interface{}{
			"username": "testuser",
			"email":    "test@example.com",
			"password": "password123",
		},
	}
	
	body, _ := json.Marshal(registerBody)
	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handlers.RegisterUser(w, req)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		setContext     bool
	}{
		{
			name: "Valid update with context",
			requestBody: map[string]interface{}{
				"user": map[string]interface{}{
					"bio": "Updated bio",
				},
			},
			expectedStatus: http.StatusOK,
			setContext:     true,
		},
		{
			name: "Update without context",
			requestBody: map[string]interface{}{
				"user": map[string]interface{}{
					"bio": "Updated bio",
				},
			},
			expectedStatus: http.StatusUnauthorized,
			setContext:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/api/user", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			// Set context if needed (normally done by middleware)
			if tt.setContext {
				ctx := req.Context()
				ctx = context.WithValue(ctx, middleware.UserIDContextKey, int64(1))
				ctx = context.WithValue(ctx, middleware.UsernameContextKey, "testuser")
				req = req.WithContext(ctx)
			}
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Call handler
			handlers.UpdateUser(w, req)
			
			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestAuthHandlers_DuplicateRegistration(t *testing.T) {
	handlers, db := setupTestHandlers(t)
	defer cleanupTestDB(db)

	// Register first user
	registerBody := map[string]interface{}{
		"user": map[string]interface{}{
			"username": "testuser",
			"email":    "test@example.com",
			"password": "password123",
		},
	}
	
	body, _ := json.Marshal(registerBody)
	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handlers.RegisterUser(w, req)
	
	if w.Code != http.StatusCreated {
		t.Fatalf("First registration failed: %d", w.Code)
	}

	// Try to register with same email
	req2 := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	handlers.RegisterUser(w2, req2)
	
	if w2.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for duplicate email, got %d", http.StatusBadRequest, w2.Code)
	}

	// Try to register with same username but different email
	duplicateUsernameBody := map[string]interface{}{
		"user": map[string]interface{}{
			"username": "testuser", // same username
			"email":    "different@example.com", // different email
			"password": "password123",
		},
	}
	
	body3, _ := json.Marshal(duplicateUsernameBody)
	req3 := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewReader(body3))
	req3.Header.Set("Content-Type", "application/json")
	w3 := httptest.NewRecorder()
	handlers.RegisterUser(w3, req3)
	
	if w3.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for duplicate username, got %d", http.StatusBadRequest, w3.Code)
	}
}