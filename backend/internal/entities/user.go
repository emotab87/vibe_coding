package entities

import (
	"regexp"
	"strings"
	"time"
)

// User represents a user in the system
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	ImageURL string `json:"image"`
	
	// Internal fields (not exposed in API)
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

// UserRegistration represents user registration request
type UserRegistration struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserLogin represents user login request
type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserUpdate represents user update request
type UserUpdate struct {
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	Bio      *string `json:"bio,omitempty"`
	ImageURL *string `json:"image,omitempty"`
	Password *string `json:"password,omitempty"`
}

// UserResponse represents user data returned by API
type UserResponse struct {
	User UserData `json:"user"`
}

// UserData represents user data in API response
type UserData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	ImageURL string `json:"image"`
	Token    string `json:"token"`
}

// ValidationError represents validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (ve *ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve.Errors {
		messages = append(messages, err.Field+": "+err.Message)
	}
	return strings.Join(messages, ", ")
}

// Validate validates user registration data
func (ur *UserRegistration) Validate() *ValidationErrors {
	var errors []ValidationError

	// Username validation
	if ur.Username == "" {
		errors = append(errors, ValidationError{
			Field:   "username",
			Message: "username is required",
		})
	} else if len(ur.Username) < 3 {
		errors = append(errors, ValidationError{
			Field:   "username",
			Message: "username must be at least 3 characters long",
		})
	} else if len(ur.Username) > 50 {
		errors = append(errors, ValidationError{
			Field:   "username",
			Message: "username must be less than 50 characters long",
		})
	} else if !isValidUsername(ur.Username) {
		errors = append(errors, ValidationError{
			Field:   "username",
			Message: "username can only contain letters, numbers, and underscores",
		})
	}

	// Email validation
	if ur.Email == "" {
		errors = append(errors, ValidationError{
			Field:   "email",
			Message: "email is required",
		})
	} else if !isValidEmail(ur.Email) {
		errors = append(errors, ValidationError{
			Field:   "email",
			Message: "email format is invalid",
		})
	}

	// Password validation
	if ur.Password == "" {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "password is required",
		})
	} else if len(ur.Password) < 6 {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "password must be at least 6 characters long",
		})
	} else if len(ur.Password) > 100 {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "password must be less than 100 characters long",
		})
	}

	if len(errors) > 0 {
		return &ValidationErrors{Errors: errors}
	}
	return nil
}

// Validate validates user login data
func (ul *UserLogin) Validate() *ValidationErrors {
	var errors []ValidationError

	// Email validation
	if ul.Email == "" {
		errors = append(errors, ValidationError{
			Field:   "email",
			Message: "email is required",
		})
	} else if !isValidEmail(ul.Email) {
		errors = append(errors, ValidationError{
			Field:   "email",
			Message: "email format is invalid",
		})
	}

	// Password validation
	if ul.Password == "" {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "password is required",
		})
	}

	if len(errors) > 0 {
		return &ValidationErrors{Errors: errors}
	}
	return nil
}

// Validate validates user update data
func (uu *UserUpdate) Validate() *ValidationErrors {
	var errors []ValidationError

	// Username validation (if provided)
	if uu.Username != nil {
		username := *uu.Username
		if username != "" {
			if len(username) < 3 {
				errors = append(errors, ValidationError{
					Field:   "username",
					Message: "username must be at least 3 characters long",
				})
			} else if len(username) > 50 {
				errors = append(errors, ValidationError{
					Field:   "username",
					Message: "username must be less than 50 characters long",
				})
			} else if !isValidUsername(username) {
				errors = append(errors, ValidationError{
					Field:   "username",
					Message: "username can only contain letters, numbers, and underscores",
				})
			}
		}
	}

	// Email validation (if provided)
	if uu.Email != nil {
		email := *uu.Email
		if email != "" && !isValidEmail(email) {
			errors = append(errors, ValidationError{
				Field:   "email",
				Message: "email format is invalid",
			})
		}
	}

	// Password validation (if provided)
	if uu.Password != nil {
		password := *uu.Password
		if password != "" {
			if len(password) < 6 {
				errors = append(errors, ValidationError{
					Field:   "password",
					Message: "password must be at least 6 characters long",
				})
			} else if len(password) > 100 {
				errors = append(errors, ValidationError{
					Field:   "password",
					Message: "password must be less than 100 characters long",
				})
			}
		}
	}

	// Bio validation (if provided)
	if uu.Bio != nil && len(*uu.Bio) > 500 {
		errors = append(errors, ValidationError{
			Field:   "bio",
			Message: "bio must be less than 500 characters long",
		})
	}

	if len(errors) > 0 {
		return &ValidationErrors{Errors: errors}
	}
	return nil
}

// ToUserData converts User to UserData with token
func (u *User) ToUserData(token string) UserData {
	return UserData{
		Username: u.Username,
		Email:    u.Email,
		Bio:      u.Bio,
		ImageURL: u.ImageURL,
		Token:    token,
	}
}

// ToUserResponse converts User to UserResponse with token
func (u *User) ToUserResponse(token string) UserResponse {
	return UserResponse{
		User: u.ToUserData(token),
	}
}

// Helper functions
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return emailRegex.MatchString(strings.ToLower(email))
}

func isValidUsername(username string) bool {
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return usernameRegex.MatchString(username)
}