package entities

import (
	"testing"
)

func TestUserRegistrationValidate(t *testing.T) {
	tests := []struct {
		name     string
		user     UserRegistration
		wantErr  bool
		errorMsg string
	}{
		{
			name: "Valid registration",
			user: UserRegistration{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "Empty username",
			user: UserRegistration{
				Username: "",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr:  true,
			errorMsg: "username is required",
		},
		{
			name: "Username too short",
			user: UserRegistration{
				Username: "ab",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr:  true,
			errorMsg: "username must be at least 3 characters long",
		},
		{
			name: "Username too long",
			user: UserRegistration{
				Username: "this_is_a_very_long_username_that_exceeds_fifty_chars",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr:  true,
			errorMsg: "username must be less than 50 characters long",
		},
		{
			name: "Invalid username characters",
			user: UserRegistration{
				Username: "test-user!",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr:  true,
			errorMsg: "username can only contain letters, numbers, and underscores",
		},
		{
			name: "Empty email",
			user: UserRegistration{
				Username: "testuser",
				Email:    "",
				Password: "password123",
			},
			wantErr:  true,
			errorMsg: "email is required",
		},
		{
			name: "Invalid email format",
			user: UserRegistration{
				Username: "testuser",
				Email:    "invalid-email",
				Password: "password123",
			},
			wantErr:  true,
			errorMsg: "email format is invalid",
		},
		{
			name: "Empty password",
			user: UserRegistration{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "",
			},
			wantErr:  true,
			errorMsg: "password is required",
		},
		{
			name: "Password too short",
			user: UserRegistration{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "12345",
			},
			wantErr:  true,
			errorMsg: "password must be at least 6 characters long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRegistration.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && len(tt.errorMsg) > 0 {
				found := false
				for _, validationErr := range err.Errors {
					if validationErr.Message == tt.errorMsg {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error message '%s' not found in validation errors: %v", tt.errorMsg, err.Errors)
				}
			}
		})
	}
}

func TestUserLoginValidate(t *testing.T) {
	tests := []struct {
		name     string
		user     UserLogin
		wantErr  bool
		errorMsg string
	}{
		{
			name: "Valid login",
			user: UserLogin{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "Empty email",
			user: UserLogin{
				Email:    "",
				Password: "password123",
			},
			wantErr:  true,
			errorMsg: "email is required",
		},
		{
			name: "Invalid email format",
			user: UserLogin{
				Email:    "invalid-email",
				Password: "password123",
			},
			wantErr:  true,
			errorMsg: "email format is invalid",
		},
		{
			name: "Empty password",
			user: UserLogin{
				Email:    "test@example.com",
				Password: "",
			},
			wantErr:  true,
			errorMsg: "password is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("UserLogin.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && len(tt.errorMsg) > 0 {
				found := false
				for _, validationErr := range err.Errors {
					if validationErr.Message == tt.errorMsg {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error message '%s' not found in validation errors: %v", tt.errorMsg, err.Errors)
				}
			}
		})
	}
}

func TestUserUpdateValidate(t *testing.T) {
	tests := []struct {
		name     string
		user     UserUpdate
		wantErr  bool
		errorMsg string
	}{
		{
			name: "Valid update - all fields",
			user: UserUpdate{
				Username: stringPtr("newuser"),
				Email:    stringPtr("new@example.com"),
				Bio:      stringPtr("New bio"),
				ImageURL: stringPtr("https://example.com/image.jpg"),
				Password: stringPtr("newpassword123"),
			},
			wantErr: false,
		},
		{
			name: "Valid update - partial fields",
			user: UserUpdate{
				Bio: stringPtr("Updated bio only"),
			},
			wantErr: false,
		},
		{
			name: "Invalid username too short",
			user: UserUpdate{
				Username: stringPtr("ab"),
			},
			wantErr:  true,
			errorMsg: "username must be at least 3 characters long",
		},
		{
			name: "Invalid email format",
			user: UserUpdate{
				Email: stringPtr("invalid-email"),
			},
			wantErr:  true,
			errorMsg: "email format is invalid",
		},
		{
			name: "Invalid password too short",
			user: UserUpdate{
				Password: stringPtr("12345"),
			},
			wantErr:  true,
			errorMsg: "password must be at least 6 characters long",
		},
		{
			name: "Bio too long",
			user: UserUpdate{
				Bio: stringPtr(generateLongString(501)),
			},
			wantErr:  true,
			errorMsg: "bio must be less than 500 characters long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("UserUpdate.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && len(tt.errorMsg) > 0 {
				found := false
				for _, validationErr := range err.Errors {
					if validationErr.Message == tt.errorMsg {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error message '%s' not found in validation errors: %v", tt.errorMsg, err.Errors)
				}
			}
		})
	}
}

func TestUserToUserData(t *testing.T) {
	user := &User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Bio:      "Test bio",
		ImageURL: "https://example.com/image.jpg",
	}

	token := "test-jwt-token"
	userData := user.ToUserData(token)

	if userData.Username != user.Username {
		t.Errorf("Expected Username %s, got %s", user.Username, userData.Username)
	}
	if userData.Email != user.Email {
		t.Errorf("Expected Email %s, got %s", user.Email, userData.Email)
	}
	if userData.Bio != user.Bio {
		t.Errorf("Expected Bio %s, got %s", user.Bio, userData.Bio)
	}
	if userData.ImageURL != user.ImageURL {
		t.Errorf("Expected ImageURL %s, got %s", user.ImageURL, userData.ImageURL)
	}
	if userData.Token != token {
		t.Errorf("Expected Token %s, got %s", token, userData.Token)
	}
}

func TestUserToUserResponse(t *testing.T) {
	user := &User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Bio:      "Test bio",
		ImageURL: "https://example.com/image.jpg",
	}

	token := "test-jwt-token"
	userResponse := user.ToUserResponse(token)

	if userResponse.User.Username != user.Username {
		t.Errorf("Expected Username %s, got %s", user.Username, userResponse.User.Username)
	}
	if userResponse.User.Token != token {
		t.Errorf("Expected Token %s, got %s", token, userResponse.User.Token)
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"user.name@domain.co.uk", true},
		{"test+tag@example.org", true},
		{"invalid-email", false},
		{"@example.com", false},
		{"test@", false},
		{"", false},
		{"test.example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			result := isValidEmail(tt.email)
			if result != tt.valid {
				t.Errorf("isValidEmail(%s) = %v, want %v", tt.email, result, tt.valid)
			}
		})
	}
}

func TestIsValidUsername(t *testing.T) {
	tests := []struct {
		username string
		valid    bool
	}{
		{"testuser", true},
		{"test_user", true},
		{"user123", true},
		{"User123", true},
		{"test-user", false},
		{"test user", false},
		{"test@user", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.username, func(t *testing.T) {
			result := isValidUsername(tt.username)
			if result != tt.valid {
				t.Errorf("isValidUsername(%s) = %v, want %v", tt.username, result, tt.valid)
			}
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func generateLongString(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = 'a'
	}
	return string(result)
}