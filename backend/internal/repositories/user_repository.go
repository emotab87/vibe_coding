package repositories

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/emotab87/vibe_coding/backend/internal/database"
	"github.com/emotab87/vibe_coding/backend/internal/entities"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(user *entities.UserRegistration) (*entities.User, error)
	GetByEmail(email string) (*entities.User, error)
	GetByUsername(username string) (*entities.User, error)
	GetByID(id int64) (*entities.User, error)
	Update(id int64, updates *entities.UserUpdate) (*entities.User, error)
	EmailExists(email string) (bool, error)
	UsernameExists(username string) (bool, error)
	VerifyPassword(user *entities.User, password string) bool
}

// userRepository implements UserRepository using direct SQL
type userRepository struct {
	db *database.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *database.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create creates a new user
func (r *userRepository) Create(userReg *entities.UserRegistration) (*entities.User, error) {
	// Hash password
	hashedPassword, err := hashPassword(userReg.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	now := time.Now()
	
	query := `
		INSERT INTO users (username, email, password_hash, bio, image_url, created_at, updated_at)
		VALUES (?, ?, ?, '', '', ?, ?)
		RETURNING id, username, email, bio, image_url, created_at, updated_at
	`
	
	user := &entities.User{}
	err = r.db.QueryRow(query, 
		userReg.Username, 
		userReg.Email, 
		hashedPassword,
		now,
		now,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Bio,
		&user.ImageURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if isUniqueConstraintError(err) {
			return nil, fmt.Errorf("user with this email or username already exists")
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	
	user.PasswordHash = hashedPassword
	return user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(email string) (*entities.User, error) {
	query := `
		SELECT id, username, email, password_hash, bio, image_url, created_at, updated_at
		FROM users 
		WHERE email = ?
	`
	
	user := &entities.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Bio,
		&user.ImageURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	
	return user, nil
}

// GetByUsername retrieves a user by username
func (r *userRepository) GetByUsername(username string) (*entities.User, error) {
	query := `
		SELECT id, username, email, password_hash, bio, image_url, created_at, updated_at
		FROM users 
		WHERE username = ?
	`
	
	user := &entities.User{}
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Bio,
		&user.ImageURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	
	return user, nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(id int64) (*entities.User, error) {
	query := `
		SELECT id, username, email, password_hash, bio, image_url, created_at, updated_at
		FROM users 
		WHERE id = ?
	`
	
	user := &entities.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Bio,
		&user.ImageURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	
	return user, nil
}

// Update updates user information
func (r *userRepository) Update(id int64, updates *entities.UserUpdate) (*entities.User, error) {
	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	
	if updates.Username != nil {
		setParts = append(setParts, "username = ?")
		args = append(args, *updates.Username)
	}
	
	if updates.Email != nil {
		setParts = append(setParts, "email = ?")
		args = append(args, *updates.Email)
	}
	
	if updates.Bio != nil {
		setParts = append(setParts, "bio = ?")
		args = append(args, *updates.Bio)
	}
	
	if updates.ImageURL != nil {
		setParts = append(setParts, "image_url = ?")
		args = append(args, *updates.ImageURL)
	}
	
	if updates.Password != nil {
		hashedPassword, err := hashPassword(*updates.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		setParts = append(setParts, "password_hash = ?")
		args = append(args, hashedPassword)
	}
	
	if len(setParts) == 0 {
		// No updates requested, just return current user
		return r.GetByID(id)
	}
	
	// Add updated_at and user ID
	setParts = append(setParts, "updated_at = ?")
	args = append(args, time.Now())
	args = append(args, id)
	
	query := fmt.Sprintf(`
		UPDATE users 
		SET %s
		WHERE id = ?
		RETURNING id, username, email, password_hash, bio, image_url, created_at, updated_at
	`, joinStrings(setParts, ", "))
	
	user := &entities.User{}
	err := r.db.QueryRow(query, args...).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Bio,
		&user.ImageURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		if isUniqueConstraintError(err) {
			return nil, fmt.Errorf("username or email already exists")
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	
	return user, nil
}

// EmailExists checks if an email already exists
func (r *userRepository) EmailExists(email string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM users WHERE email = ?"
	
	err := r.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	
	return count > 0, nil
}

// UsernameExists checks if a username already exists
func (r *userRepository) UsernameExists(username string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM users WHERE username = ?"
	
	err := r.db.QueryRow(query, username).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}
	
	return count > 0, nil
}

// VerifyPassword verifies a password against the stored hash
func (r *userRepository) VerifyPassword(user *entities.User, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) == nil
}

// Helper functions

// hashPassword hashes a password using bcrypt
func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// isUniqueConstraintError checks if the error is a unique constraint violation
func isUniqueConstraintError(err error) bool {
	return err != nil && 
		(containsString(err.Error(), "UNIQUE constraint failed") ||
		 containsString(err.Error(), "unique constraint"))
}

// containsString checks if a string contains a substring (case-insensitive)
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && 
		findSubstring(strings.ToLower(s), strings.ToLower(substr)) >= 0
}

// findSubstring finds the index of a substring
func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// joinStrings joins strings with a separator
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}
	
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}