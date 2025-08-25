package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/emotab87/vibe_coding/backend/internal/entities"
	"github.com/emotab87/vibe_coding/backend/internal/repositories"
)

// CommentHandlers handles comment-related HTTP requests
type CommentHandlers struct {
	commentRepo repositories.CommentRepository
	articleRepo repositories.ArticleRepository
}

// NewCommentHandlers creates a new comment handlers instance
func NewCommentHandlers(commentRepo repositories.CommentRepository, articleRepo repositories.ArticleRepository) *CommentHandlers {
	return &CommentHandlers{
		commentRepo: commentRepo,
		articleRepo: articleRepo,
	}
}

// CreateComment handles comment creation
func (h *CommentHandlers) CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, err := getUserIDFromContext(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get slug from URL path
	vars := mux.Vars(r)
	slug := vars["slug"]
	if slug == "" {
		writeError(w, http.StatusBadRequest, "Missing article slug")
		return
	}

	// Check if article exists and get its ID
	article, err := h.articleRepo.GetBySlug(slug)
	if err != nil {
		if containsString(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "Article not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get article")
		return
	}

	// Parse request body
	var req struct {
		Comment entities.CommentCreate `json:"comment"`
	}

	if err := parseJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Validate comment data
	if validationErr := req.Comment.Validate(); validationErr != nil {
		writeValidationErrors(w, validationErr)
		return
	}

	// Create comment
	comment, err := h.commentRepo.Create(userID, article.ID, &req.Comment)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create comment")
		return
	}

	// Return comment response
	response := comment.ToCommentResponse()
	writeJSON(w, http.StatusCreated, response)
}

// GetCommentsByArticle handles comment listing for an article
func (h *CommentHandlers) GetCommentsByArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get slug from URL path
	vars := mux.Vars(r)
	slug := vars["slug"]
	if slug == "" {
		writeError(w, http.StatusBadRequest, "Missing article slug")
		return
	}

	// Check if article exists
	_, err := h.articleRepo.GetBySlug(slug)
	if err != nil {
		if containsString(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "Article not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get article")
		return
	}

	// Get comments for the article
	comments, err := h.commentRepo.GetByArticleSlug(slug)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get comments")
		return
	}

	// Return comments response
	response := entities.CommentsResponse{
		Comments: comments,
	}
	writeJSON(w, http.StatusOK, response)
}

// DeleteComment handles comment deletion
func (h *CommentHandlers) DeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, err := getUserIDFromContext(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get slug and comment ID from URL path
	vars := mux.Vars(r)
	slug := vars["slug"]
	commentIDStr := vars["id"]
	
	if slug == "" {
		writeError(w, http.StatusBadRequest, "Missing article slug")
		return
	}
	
	if commentIDStr == "" {
		writeError(w, http.StatusBadRequest, "Missing comment ID")
		return
	}

	// Parse comment ID
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	// Check if article exists
	_, err = h.articleRepo.GetBySlug(slug)
	if err != nil {
		if containsString(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "Article not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get article")
		return
	}

	// Check if comment exists
	existingComment, err := h.commentRepo.GetByID(commentID)
	if err != nil {
		if containsString(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "Comment not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get comment")
		return
	}

	// Check if user is the author
	if existingComment.AuthorID != userID {
		writeError(w, http.StatusForbidden, "You can only delete your own comments")
		return
	}

	// Delete comment
	if err := h.commentRepo.Delete(commentID); err != nil {
		if containsString(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "Comment not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to delete comment")
		return
	}

	// Return 204 No Content for successful deletion
	w.WriteHeader(http.StatusNoContent)
}

// Helper function to check string contains (case-insensitive)
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(toLowerCase(s), toLowerCase(substr)) >= 0
}

// Helper function to convert to lowercase
func toLowerCase(s string) string {
	result := make([]rune, len(s))
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			result[i] = r + 32
		} else {
			result[i] = r
		}
	}
	return string(result)
}

// Helper function to find substring
func findSubstring(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}

	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}