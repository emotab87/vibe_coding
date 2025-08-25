package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for our application
type Config struct {
	Environment     string
	Port            string
	Host            string
	DatabasePath    string
	JWTSecret       string
	JWTExpiryHours  int
	CORSOrigins     string
	LogLevel        string
	LogFormat       string
	BcryptRounds    int
	DebugSQL        bool
	DebugCORS       bool
	AIREnabled      bool
}

// LoadConfig loads configuration from environment variables with sensible defaults
func LoadConfig() *Config {
	return &Config{
		Environment:     getEnvOrDefault("ENV", "development"),
		Port:            getEnvOrDefault("PORT", "8080"),
		Host:            getEnvOrDefault("HOST", "localhost"),
		DatabasePath:    getEnvOrDefault("DB_PATH", "./data/conduit.db"),
		JWTSecret:       getEnvOrDefault("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production"),
		JWTExpiryHours:  getEnvIntOrDefault("JWT_EXPIRY_HOURS", 72),
		CORSOrigins:     getEnvOrDefault("CORS_ORIGINS", "http://localhost:3000"),
		LogLevel:        getEnvOrDefault("LOG_LEVEL", "debug"),
		LogFormat:       getEnvOrDefault("LOG_FORMAT", "json"),
		BcryptRounds:    getEnvIntOrDefault("BCRYPT_ROUNDS", 12),
		DebugSQL:        getEnvBoolOrDefault("DEBUG_SQL", true),
		DebugCORS:       getEnvBoolOrDefault("DEBUG_CORS", true),
		AIREnabled:      getEnvBoolOrDefault("AIR_ENABLED", true),
	}
}

// ServerAddress returns the full server address (host:port)
func (c *Config) ServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// IsDevelopment returns true if we're in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if we're in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// Validate checks if all required configuration is present
func (c *Config) Validate() error {
	if c.JWTSecret == "" || c.JWTSecret == "your-super-secret-jwt-key-change-this-in-production" {
		if c.IsProduction() {
			return fmt.Errorf("JWT_SECRET must be set in production")
		}
	}

	if c.Port == "" {
		return fmt.Errorf("PORT must be set")
	}

	return nil
}

// Helper functions for environment variable parsing

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}