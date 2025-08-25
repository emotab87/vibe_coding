package repositories

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/emotab87/vibe_coding/backend/internal/database"
	"github.com/emotab87/vibe_coding/backend/internal/entities"
)

// ArticleRepository defines the interface for article data operations
type ArticleRepository interface {
	Create(authorID int64, article *entities.ArticleCreate) (*entities.Article, error)
	GetBySlug(slug string) (*entities.Article, error)
	GetByID(id int64) (*entities.Article, error)
	Update(id int64, updates *entities.ArticleUpdate) (*entities.Article, error)
	Delete(id int64) error
	List(query *entities.ArticleListQuery) ([]entities.Article, int, error)
	SlugExists(slug string) (bool, error)
	GetExistingSlugs(baseSlug string) ([]string, error)
	IsAuthor(articleID, userID int64) (bool, error)
}

// articleRepository implements ArticleRepository using direct SQL
type articleRepository struct {
	db       *database.DB
	userRepo UserRepository
}

// NewArticleRepository creates a new article repository
func NewArticleRepository(db *database.DB, userRepo UserRepository) ArticleRepository {
	return &articleRepository{
		db:       db,
		userRepo: userRepo,
	}
}

// Create creates a new article
func (r *articleRepository) Create(authorID int64, articleCreate *entities.ArticleCreate) (*entities.Article, error) {
	// Generate base slug
	baseSlug := entities.GenerateSlug(articleCreate.Title)
	if baseSlug == "" {
		return nil, fmt.Errorf("failed to generate slug from title")
	}

	// Get existing slugs to ensure uniqueness
	existingSlugs, err := r.GetExistingSlugs(baseSlug)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing slugs: %w", err)
	}

	// Ensure unique slug
	uniqueSlug := entities.EnsureUniqueSlug(baseSlug, existingSlugs)

	now := time.Now()

	query := `
		INSERT INTO articles (slug, title, description, body, author_id, favorites_count, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, 0, ?, ?)
		RETURNING id, slug, title, description, body, author_id, favorites_count, created_at, updated_at
	`

	article := &entities.Article{}
	err = r.db.QueryRow(query,
		uniqueSlug,
		articleCreate.Title,
		articleCreate.Description,
		articleCreate.Body,
		authorID,
		now,
		now,
	).Scan(
		&article.ID,
		&article.Slug,
		&article.Title,
		&article.Description,
		&article.Body,
		&article.AuthorID,
		&article.FavoritesCount,
		&article.CreatedAt,
		&article.UpdatedAt,
	)

	if err != nil {
		if isUniqueConstraintError(err) {
			return nil, fmt.Errorf("article with this slug already exists")
		}
		return nil, fmt.Errorf("failed to create article: %w", err)
	}

	// Load author information
	if err := r.loadAuthor(article); err != nil {
		return nil, fmt.Errorf("failed to load author: %w", err)
	}

	return article, nil
}

// GetBySlug retrieves an article by slug
func (r *articleRepository) GetBySlug(slug string) (*entities.Article, error) {
	query := `
		SELECT id, slug, title, description, body, author_id, favorites_count, created_at, updated_at
		FROM articles 
		WHERE slug = ?
	`

	article := &entities.Article{}
	err := r.db.QueryRow(query, slug).Scan(
		&article.ID,
		&article.Slug,
		&article.Title,
		&article.Description,
		&article.Body,
		&article.AuthorID,
		&article.FavoritesCount,
		&article.CreatedAt,
		&article.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("article not found")
		}
		return nil, fmt.Errorf("failed to get article by slug: %w", err)
	}

	// Load author information
	if err := r.loadAuthor(article); err != nil {
		return nil, fmt.Errorf("failed to load author: %w", err)
	}

	return article, nil
}

// GetByID retrieves an article by ID
func (r *articleRepository) GetByID(id int64) (*entities.Article, error) {
	query := `
		SELECT id, slug, title, description, body, author_id, favorites_count, created_at, updated_at
		FROM articles 
		WHERE id = ?
	`

	article := &entities.Article{}
	err := r.db.QueryRow(query, id).Scan(
		&article.ID,
		&article.Slug,
		&article.Title,
		&article.Description,
		&article.Body,
		&article.AuthorID,
		&article.FavoritesCount,
		&article.CreatedAt,
		&article.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("article not found")
		}
		return nil, fmt.Errorf("failed to get article by ID: %w", err)
	}

	// Load author information
	if err := r.loadAuthor(article); err != nil {
		return nil, fmt.Errorf("failed to load author: %w", err)
	}

	return article, nil
}

// Update updates an article
func (r *articleRepository) Update(id int64, updates *entities.ArticleUpdate) (*entities.Article, error) {
	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}

	if updates.Title != nil {
		// If title is being updated, we need to generate a new slug
		baseSlug := entities.GenerateSlug(*updates.Title)
		if baseSlug == "" {
			return nil, fmt.Errorf("failed to generate slug from new title")
		}

		// Get existing slugs to ensure uniqueness (excluding current article)
		existingSlugs, err := r.getExistingSlugsExcluding(baseSlug, id)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing slugs: %w", err)
		}

		uniqueSlug := entities.EnsureUniqueSlug(baseSlug, existingSlugs)
		setParts = append(setParts, "title = ?", "slug = ?")
		args = append(args, *updates.Title, uniqueSlug)
	}

	if updates.Description != nil {
		setParts = append(setParts, "description = ?")
		args = append(args, *updates.Description)
	}

	if updates.Body != nil {
		setParts = append(setParts, "body = ?")
		args = append(args, *updates.Body)
	}

	if len(setParts) == 0 {
		// No updates requested, just return current article
		return r.GetByID(id)
	}

	// Add updated_at and article ID
	setParts = append(setParts, "updated_at = ?")
	args = append(args, time.Now())
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE articles 
		SET %s
		WHERE id = ?
		RETURNING id, slug, title, description, body, author_id, favorites_count, created_at, updated_at
	`, joinStrings(setParts, ", "))

	article := &entities.Article{}
	err := r.db.QueryRow(query, args...).Scan(
		&article.ID,
		&article.Slug,
		&article.Title,
		&article.Description,
		&article.Body,
		&article.AuthorID,
		&article.FavoritesCount,
		&article.CreatedAt,
		&article.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("article not found")
		}
		if isUniqueConstraintError(err) {
			return nil, fmt.Errorf("slug already exists")
		}
		return nil, fmt.Errorf("failed to update article: %w", err)
	}

	// Load author information
	if err := r.loadAuthor(article); err != nil {
		return nil, fmt.Errorf("failed to load author: %w", err)
	}

	return article, nil
}

// Delete deletes an article
func (r *articleRepository) Delete(id int64) error {
	query := "DELETE FROM articles WHERE id = ?"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete article: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("article not found")
	}

	return nil
}

// List retrieves articles with pagination and filtering
func (r *articleRepository) List(query *entities.ArticleListQuery) ([]entities.Article, int, error) {
	// Set default values
	if query.Limit <= 0 {
		query.Limit = 20
	}
	if query.Limit > 100 {
		query.Limit = 100
	}
	if query.Offset < 0 {
		query.Offset = 0
	}

	// Build WHERE clause
	whereParts := []string{}
	args := []interface{}{}

	if query.Author != "" {
		whereParts = append(whereParts, "u.username = ?")
		args = append(args, query.Author)
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = "WHERE " + joinStrings(whereParts, " AND ")
	}

	// Get total count
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM articles a
		JOIN users u ON a.author_id = u.id
		%s
	`, whereClause)

	var totalCount int
	err := r.db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// Get articles
	articlesQuery := fmt.Sprintf(`
		SELECT a.id, a.slug, a.title, a.description, a.body, a.author_id, a.favorites_count, a.created_at, a.updated_at
		FROM articles a
		JOIN users u ON a.author_id = u.id
		%s
		ORDER BY a.created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	// Add limit and offset to args
	queryArgs := append(args, query.Limit, query.Offset)

	rows, err := r.db.Query(articlesQuery, queryArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query articles: %w", err)
	}
	defer rows.Close()

	var articles []entities.Article
	for rows.Next() {
		var article entities.Article
		err := rows.Scan(
			&article.ID,
			&article.Slug,
			&article.Title,
			&article.Description,
			&article.Body,
			&article.AuthorID,
			&article.FavoritesCount,
			&article.CreatedAt,
			&article.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan article: %w", err)
		}

		// Load author information
		if err := r.loadAuthor(&article); err != nil {
			return nil, 0, fmt.Errorf("failed to load author: %w", err)
		}

		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate over articles: %w", err)
	}

	return articles, totalCount, nil
}

// SlugExists checks if a slug already exists
func (r *articleRepository) SlugExists(slug string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM articles WHERE slug = ?"

	err := r.db.QueryRow(query, slug).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check slug existence: %w", err)
	}

	return count > 0, nil
}

// GetExistingSlugs gets existing slugs that start with the base slug
func (r *articleRepository) GetExistingSlugs(baseSlug string) ([]string, error) {
	query := "SELECT slug FROM articles WHERE slug LIKE ? ORDER BY slug"
	pattern := baseSlug + "%"

	rows, err := r.db.Query(query, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to query existing slugs: %w", err)
	}
	defer rows.Close()

	var slugs []string
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return nil, fmt.Errorf("failed to scan slug: %w", err)
		}
		slugs = append(slugs, slug)
	}

	return slugs, rows.Err()
}

// getExistingSlugsExcluding gets existing slugs excluding a specific article ID
func (r *articleRepository) getExistingSlugsExcluding(baseSlug string, excludeID int64) ([]string, error) {
	query := "SELECT slug FROM articles WHERE slug LIKE ? AND id != ? ORDER BY slug"
	pattern := baseSlug + "%"

	rows, err := r.db.Query(query, pattern, excludeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query existing slugs: %w", err)
	}
	defer rows.Close()

	var slugs []string
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return nil, fmt.Errorf("failed to scan slug: %w", err)
		}
		slugs = append(slugs, slug)
	}

	return slugs, rows.Err()
}

// IsAuthor checks if a user is the author of an article
func (r *articleRepository) IsAuthor(articleID, userID int64) (bool, error) {
	query := "SELECT author_id FROM articles WHERE id = ?"

	var authorID int64
	err := r.db.QueryRow(query, articleID).Scan(&authorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check article author: %w", err)
	}

	return authorID == userID, nil
}

// loadAuthor loads author information for an article
func (r *articleRepository) loadAuthor(article *entities.Article) error {
	author, err := r.userRepo.GetByID(article.AuthorID)
	if err != nil {
		return err
	}

	// Create author data without sensitive information
	article.Author = &entities.User{
		ID:       author.ID,
		Username: author.Username,
		Bio:      author.Bio,
		ImageURL: author.ImageURL,
	}

	return nil
}

// Helper functions

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