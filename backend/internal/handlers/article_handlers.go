package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/emotab87/vibe_coding/backend/internal/entities"
	"github.com/emotab87/vibe_coding/backend/internal/repositories"
)

// ArticleHandlers handles article-related HTTP requests
type ArticleHandlers struct {
	articleRepo repositories.ArticleRepository
}

// NewArticleHandlers creates a new article handlers instance
func NewArticleHandlers(articleRepo repositories.ArticleRepository) *ArticleHandlers {
	return &ArticleHandlers{
		articleRepo: articleRepo,
	}
}

// CreateArticle handles article creation
func (h *ArticleHandlers) CreateArticle(w http.ResponseWriter, r *http.Request) {
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

	// Parse request body
	var req struct {
		Article entities.ArticleCreate `json:"article"`
	}

	if err := parseJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Validate article data
	if validationErr := req.Article.Validate(); validationErr != nil {
		writeValidationErrors(w, validationErr)
		return
	}

	// Create article
	article, err := h.articleRepo.Create(userID, &req.Article)
	if err != nil {
		if containsString(err.Error(), "already exists") {
			writeError(w, http.StatusConflict, "Article with this title already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to create article")
		return
	}

	// Return article response
	response := article.ToArticleResponse()
	writeJSON(w, http.StatusCreated, response)
}

// GetArticle handles article retrieval by slug
func (h *ArticleHandlers) GetArticle(w http.ResponseWriter, r *http.Request) {
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

	// Get article by slug
	article, err := h.articleRepo.GetBySlug(slug)
	if err != nil {
		if containsString(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "Article not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get article")
		return
	}

	// Return article response
	response := article.ToArticleResponse()
	writeJSON(w, http.StatusOK, response)
}

// UpdateArticle handles article updates
func (h *ArticleHandlers) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
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

	// Get existing article to check authorization
	existingArticle, err := h.articleRepo.GetBySlug(slug)
	if err != nil {
		if containsString(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "Article not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get article")
		return
	}

	// Check if user is the author
	if existingArticle.AuthorID != userID {
		writeError(w, http.StatusForbidden, "You can only update your own articles")
		return
	}

	// Parse request body
	var req struct {
		Article entities.ArticleUpdate `json:"article"`
	}

	if err := parseJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Validate update data
	if validationErr := req.Article.Validate(); validationErr != nil {
		writeValidationErrors(w, validationErr)
		return
	}

	// Update article
	updatedArticle, err := h.articleRepo.Update(existingArticle.ID, &req.Article)
	if err != nil {
		if containsString(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "Article not found")
			return
		}
		if containsString(err.Error(), "already exists") {
			writeError(w, http.StatusConflict, "Article with this title already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to update article")
		return
	}

	// Return updated article response
	response := updatedArticle.ToArticleResponse()
	writeJSON(w, http.StatusOK, response)
}

// DeleteArticle handles article deletion
func (h *ArticleHandlers) DeleteArticle(w http.ResponseWriter, r *http.Request) {
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

	// Get slug from URL path
	vars := mux.Vars(r)
	slug := vars["slug"]
	if slug == "" {
		writeError(w, http.StatusBadRequest, "Missing article slug")
		return
	}

	// Get existing article to check authorization
	existingArticle, err := h.articleRepo.GetBySlug(slug)
	if err != nil {
		if containsString(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "Article not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get article")
		return
	}

	// Check if user is the author
	if existingArticle.AuthorID != userID {
		writeError(w, http.StatusForbidden, "You can only delete your own articles")
		return
	}

	// Delete article
	if err := h.articleRepo.Delete(existingArticle.ID); err != nil {
		if containsString(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "Article not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to delete article")
		return
	}

	// Return 204 No Content for successful deletion
	w.WriteHeader(http.StatusNoContent)
}

// ListArticles handles article listing with pagination
func (h *ArticleHandlers) ListArticles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse query parameters
	query := &entities.ArticleListQuery{
		Limit:  20, // Default limit
		Offset: 0,  // Default offset
	}

	// Parse limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			query.Limit = limit
		}
	}

	// Parse offset
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			query.Offset = offset
		}
	}

	// Parse author filter
	if author := r.URL.Query().Get("author"); author != "" {
		query.Author = author
	}

	// Get articles
	articles, totalCount, err := h.articleRepo.List(query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list articles")
		return
	}

	// Return articles response
	response := entities.ArticlesResponse{
		Articles:      articles,
		ArticlesCount: totalCount,
	}
	writeJSON(w, http.StatusOK, response)
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