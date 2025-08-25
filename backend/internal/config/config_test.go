package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Save original environment
	originalEnv := map[string]string{
		"ENV":      os.Getenv("ENV"),
		"PORT":     os.Getenv("PORT"),
		"HOST":     os.Getenv("HOST"),
		"DB_PATH":  os.Getenv("DB_PATH"),
	}

	// Clean up after test
	defer func() {
		for key, value := range originalEnv {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	// Test with default values
	t.Run("DefaultValues", func(t *testing.T) {
		// Clear environment variables
		os.Unsetenv("ENV")
		os.Unsetenv("PORT")
		os.Unsetenv("HOST")
		os.Unsetenv("DB_PATH")

		cfg := LoadConfig()

		if cfg.Environment != "development" {
			t.Errorf("Expected Environment 'development', got '%s'", cfg.Environment)
		}

		if cfg.Port != "8080" {
			t.Errorf("Expected Port '8080', got '%s'", cfg.Port)
		}

		if cfg.Host != "localhost" {
			t.Errorf("Expected Host 'localhost', got '%s'", cfg.Host)
		}

		if cfg.DatabasePath != "./data/conduit.db" {
			t.Errorf("Expected DatabasePath './data/conduit.db', got '%s'", cfg.DatabasePath)
		}
	})

	// Test with environment variables
	t.Run("EnvironmentVariables", func(t *testing.T) {
		os.Setenv("ENV", "production")
		os.Setenv("PORT", "3000")
		os.Setenv("HOST", "0.0.0.0")
		os.Setenv("DB_PATH", "/var/lib/conduit.db")

		cfg := LoadConfig()

		if cfg.Environment != "production" {
			t.Errorf("Expected Environment 'production', got '%s'", cfg.Environment)
		}

		if cfg.Port != "3000" {
			t.Errorf("Expected Port '3000', got '%s'", cfg.Port)
		}

		if cfg.Host != "0.0.0.0" {
			t.Errorf("Expected Host '0.0.0.0', got '%s'", cfg.Host)
		}

		if cfg.DatabasePath != "/var/lib/conduit.db" {
			t.Errorf("Expected DatabasePath '/var/lib/conduit.db', got '%s'", cfg.DatabasePath)
		}
	})
}

func TestServerAddress(t *testing.T) {
	cfg := &Config{
		Host: "localhost",
		Port: "8080",
	}

	expected := "localhost:8080"
	if addr := cfg.ServerAddress(); addr != expected {
		t.Errorf("Expected server address '%s', got '%s'", expected, addr)
	}
}

func TestIsDevelopment(t *testing.T) {
	tests := []struct {
		env      string
		expected bool
	}{
		{"development", true},
		{"production", false},
		{"staging", false},
		{"", false},
	}

	for _, test := range tests {
		cfg := &Config{Environment: test.env}
		if result := cfg.IsDevelopment(); result != test.expected {
			t.Errorf("Environment '%s': expected IsDevelopment() %v, got %v",
				test.env, test.expected, result)
		}
	}
}

func TestValidate(t *testing.T) {
	t.Run("ValidDevelopmentConfig", func(t *testing.T) {
		cfg := &Config{
			Environment: "development",
			Port:        "8080",
			JWTSecret:   "test-secret",
		}

		if err := cfg.Validate(); err != nil {
			t.Errorf("Expected valid config, got error: %v", err)
		}
	})

	t.Run("InvalidProductionConfig", func(t *testing.T) {
		cfg := &Config{
			Environment: "production",
			Port:        "8080",
			JWTSecret:   "your-super-secret-jwt-key-change-this-in-production",
		}

		if err := cfg.Validate(); err == nil {
			t.Error("Expected validation error for production config with default JWT secret")
		}
	})

	t.Run("MissingPort", func(t *testing.T) {
		cfg := &Config{
			Environment: "development",
			Port:        "",
			JWTSecret:   "test-secret",
		}

		if err := cfg.Validate(); err == nil {
			t.Error("Expected validation error for missing port")
		}
	})
}