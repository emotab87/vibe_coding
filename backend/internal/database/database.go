package database

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// DB wraps sql.DB to provide additional functionality
type DB struct {
	*sql.DB
	path string
}

// NewDB creates a new database connection
func NewDB(databasePath string) (*DB, error) {
	// Ensure directory exists
	dir := filepath.Dir(databasePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	sqlDB, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure SQLite connection
	sqlDB.SetMaxOpenConns(1) // SQLite works best with single connection
	sqlDB.SetMaxIdleConns(1)

	// Enable foreign key constraints
	if _, err := sqlDB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Enable WAL mode for better concurrency
	if _, err := sqlDB.Exec("PRAGMA journal_mode = WAL"); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	db := &DB{
		DB:   sqlDB,
		path: databasePath,
	}

	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.DB != nil {
		return db.DB.Close()
	}
	return nil
}

// Ping checks database connectivity
func (db *DB) Ping() error {
	return db.DB.Ping()
}

// Migrate runs database migrations from the migrations directory
func (db *DB) Migrate(migrationsDir string) error {
	// Create migrations table if it doesn't exist
	if err := db.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of migration files
	migrationFiles, err := getMigrationFiles(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations: %w", err)
	}

	// Get applied migrations
	appliedMigrations, err := db.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Apply pending migrations
	for _, file := range migrationFiles {
		if _, applied := appliedMigrations[file]; !applied {
			if err := db.applyMigration(migrationsDir, file); err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", file, err)
			}
			log.Printf("✅ Applied migration: %s", file)
		}
	}

	log.Printf("🎉 Database migrations completed successfully")
	return nil
}

// createMigrationsTable creates the migrations tracking table
func (db *DB) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename TEXT PRIMARY KEY,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.DB.Exec(query)
	return err
}

// getMigrationFiles returns sorted list of migration files
func getMigrationFiles(migrationsDir string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(migrationsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".sql") {
			files = append(files, d.Name())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort files to ensure proper migration order
	sort.Strings(files)
	return files, nil
}

// getAppliedMigrations returns a map of applied migration filenames
func (db *DB) getAppliedMigrations() (map[string]bool, error) {
	rows, err := db.DB.Query("SELECT filename FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			return nil, err
		}
		applied[filename] = true
	}

	return applied, rows.Err()
}

// applyMigration applies a single migration file
func (db *DB) applyMigration(migrationsDir, filename string) error {
	// Read migration file
	filePath := filepath.Join(migrationsDir, filename)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Extract UP migration (ignore DOWN for now)
	migrationSQL := extractUpMigration(string(content))
	if migrationSQL == "" {
		return fmt.Errorf("no UP migration found in %s", filename)
	}

	// Begin transaction
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute migration
	if _, err := tx.Exec(migrationSQL); err != nil {
		return err
	}

	// Record migration as applied
	if _, err := tx.Exec(
		"INSERT INTO schema_migrations (filename) VALUES (?)",
		filename,
	); err != nil {
		return err
	}

	return tx.Commit()
}

// extractUpMigration extracts the UP migration from the content
func extractUpMigration(content string) string {
	lines := strings.Split(content, "\n")
	var upLines []string
	inUpSection := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "-- +migrate Up" {
			inUpSection = true
			continue
		}

		if trimmed == "-- +migrate Down" {
			break
		}

		if inUpSection && !strings.HasPrefix(trimmed, "--") {
			upLines = append(upLines, line)
		}
	}

	return strings.Join(upLines, "\n")
}

// Transaction helper method
func (db *DB) Transaction(fn func(*sql.Tx) error) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}