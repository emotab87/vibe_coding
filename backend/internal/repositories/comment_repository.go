package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/emotab87/vibe_coding/backend/internal/database"
	"github.com/emotab87/vibe_coding/backend/internal/entities"
)

// CommentRepository defines the interface for comment data operations
type CommentRepository interface {
	Create(authorID, articleID int64, comment *entities.CommentCreate) (*entities.Comment, error)
	GetByArticleSlug(slug string) ([]entities.Comment, error)
	GetByID(id int64) (*entities.Comment, error)
	Delete(id int64) error
	IsAuthor(commentID, userID int64) (bool, error)
}

// commentRepository implements CommentRepository using direct SQL
type commentRepository struct {
	db       *database.DB
	userRepo UserRepository
}

// NewCommentRepository creates a new comment repository
func NewCommentRepository(db *database.DB, userRepo UserRepository) CommentRepository {
	return &commentRepository{
		db:       db,
		userRepo: userRepo,
	}
}

// Create creates a new comment
func (r *commentRepository) Create(authorID, articleID int64, commentCreate *entities.CommentCreate) (*entities.Comment, error) {
	now := time.Now()

	query := `
		INSERT INTO comments (body, author_id, article_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id, body, author_id, article_id, created_at, updated_at
	`

	comment := &entities.Comment{}
	err := r.db.QueryRow(query,
		commentCreate.Body,
		authorID,
		articleID,
		now,
		now,
	).Scan(
		&comment.ID,
		&comment.Body,
		&comment.AuthorID,
		&comment.ArticleID,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// Load author information
	if err := r.loadAuthor(comment); err != nil {
		return nil, fmt.Errorf("failed to load author: %w", err)
	}

	return comment, nil
}

// GetByArticleSlug retrieves all comments for an article by slug
func (r *commentRepository) GetByArticleSlug(slug string) ([]entities.Comment, error) {
	query := `
		SELECT c.id, c.body, c.author_id, c.article_id, c.created_at, c.updated_at
		FROM comments c
		JOIN articles a ON c.article_id = a.id
		WHERE a.slug = ?
		ORDER BY c.created_at ASC
	`

	rows, err := r.db.Query(query, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()

	var comments []entities.Comment
	for rows.Next() {
		var comment entities.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.Body,
			&comment.AuthorID,
			&comment.ArticleID,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}

		// Load author information
		if err := r.loadAuthor(&comment); err != nil {
			return nil, fmt.Errorf("failed to load author: %w", err)
		}

		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over comments: %w", err)
	}

	return comments, nil
}

// GetByID retrieves a comment by ID
func (r *commentRepository) GetByID(id int64) (*entities.Comment, error) {
	query := `
		SELECT id, body, author_id, article_id, created_at, updated_at
		FROM comments 
		WHERE id = ?
	`

	comment := &entities.Comment{}
	err := r.db.QueryRow(query, id).Scan(
		&comment.ID,
		&comment.Body,
		&comment.AuthorID,
		&comment.ArticleID,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("comment not found")
		}
		return nil, fmt.Errorf("failed to get comment by ID: %w", err)
	}

	// Load author information
	if err := r.loadAuthor(comment); err != nil {
		return nil, fmt.Errorf("failed to load author: %w", err)
	}

	return comment, nil
}

// Delete deletes a comment
func (r *commentRepository) Delete(id int64) error {
	query := "DELETE FROM comments WHERE id = ?"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("comment not found")
	}

	return nil
}

// IsAuthor checks if a user is the author of a comment
func (r *commentRepository) IsAuthor(commentID, userID int64) (bool, error) {
	query := "SELECT author_id FROM comments WHERE id = ?"

	var authorID int64
	err := r.db.QueryRow(query, commentID).Scan(&authorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check comment author: %w", err)
	}

	return authorID == userID, nil
}

// loadAuthor loads author information for a comment
func (r *commentRepository) loadAuthor(comment *entities.Comment) error {
	author, err := r.userRepo.GetByID(comment.AuthorID)
	if err != nil {
		return err
	}

	// Create author data without sensitive information
	comment.Author = &entities.User{
		ID:       author.ID,
		Username: author.Username,
		Bio:      author.Bio,
		ImageURL: author.ImageURL,
	}

	return nil
}