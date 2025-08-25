package entities

import (
	"strings"
	"time"
)

// Comment represents a comment in the system
type Comment struct {
	ID        int64     `json:"id"`
	Body      string    `json:"body"`
	AuthorID  int64     `json:"-"`
	Author    *User     `json:"author,omitempty"`
	ArticleID int64     `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CommentCreate represents comment creation request
type CommentCreate struct {
	Body string `json:"body"`
}

// CommentResponse represents single comment API response
type CommentResponse struct {
	Comment Comment `json:"comment"`
}

// CommentsResponse represents multiple comments API response
type CommentsResponse struct {
	Comments []Comment `json:"comments"`
}

// Validate validates comment creation data
func (cc *CommentCreate) Validate() *ValidationErrors {
	var errors []ValidationError

	// Body validation
	if cc.Body == "" {
		errors = append(errors, ValidationError{
			Field:   "body",
			Message: "body is required",
		})
	} else if len(strings.TrimSpace(cc.Body)) < 1 {
		errors = append(errors, ValidationError{
			Field:   "body",
			Message: "body cannot be empty",
		})
	} else if len(cc.Body) > 10000 {
		errors = append(errors, ValidationError{
			Field:   "body",
			Message: "body must be less than 10000 characters long",
		})
	}

	if len(errors) > 0 {
		return &ValidationErrors{Errors: errors}
	}
	return nil
}

// ToCommentResponse converts Comment to CommentResponse
func (c *Comment) ToCommentResponse() CommentResponse {
	return CommentResponse{
		Comment: *c,
	}
}