package repositories

import (
	"testing"
	"time"

	"github.com/emotab87/vibe_coding/backend/internal/database"
	"github.com/emotab87/vibe_coding/backend/internal/entities"
)

func TestCommentRepository_Create(t *testing.T) {
	// Setup test database
	db, err := database.NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate("../../../migrations"); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Create repositories
	userRepo := NewUserRepository(db)
	articleRepo := NewArticleRepository(db, userRepo)
	commentRepo := NewCommentRepository(db, userRepo)

	// Create test user
	userReg := &entities.UserRegistration{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	user, err := userRepo.Create(userReg)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create test article
	articleCreate := &entities.ArticleCreate{
		Title:       "Test Article",
		Description: "Test description",
		Body:        "Test body",
	}
	article, err := articleRepo.Create(user.ID, articleCreate)
	if err != nil {
		t.Fatalf("Failed to create test article: %v", err)
	}

	// Test comment creation
	commentCreate := &entities.CommentCreate{
		Body: "This is a test comment",
	}

	comment, err := commentRepo.Create(user.ID, article.ID, commentCreate)
	if err != nil {
		t.Fatalf("Failed to create comment: %v", err)
	}

	// Verify comment
	if comment.Body != commentCreate.Body {
		t.Errorf("Expected comment body %s, got %s", commentCreate.Body, comment.Body)
	}
	if comment.AuthorID != user.ID {
		t.Errorf("Expected author ID %d, got %d", user.ID, comment.AuthorID)
	}
	if comment.ArticleID != article.ID {
		t.Errorf("Expected article ID %d, got %d", article.ID, comment.ArticleID)
	}
	if comment.Author == nil {
		t.Error("Expected comment to have author information")
	} else if comment.Author.Username != user.Username {
		t.Errorf("Expected author username %s, got %s", user.Username, comment.Author.Username)
	}
}

func TestCommentRepository_GetByArticleSlug(t *testing.T) {
	// Setup test database
	db, err := database.NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate("../../../migrations"); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Create repositories
	userRepo := NewUserRepository(db)
	articleRepo := NewArticleRepository(db, userRepo)
	commentRepo := NewCommentRepository(db, userRepo)

	// Create test user
	userReg := &entities.UserRegistration{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	user, err := userRepo.Create(userReg)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create test article
	articleCreate := &entities.ArticleCreate{
		Title:       "Test Article",
		Description: "Test description",
		Body:        "Test body",
	}
	article, err := articleRepo.Create(user.ID, articleCreate)
	if err != nil {
		t.Fatalf("Failed to create test article: %v", err)
	}

	// Create multiple test comments
	comments := []string{
		"First comment",
		"Second comment",
		"Third comment",
	}

	for _, body := range comments {
		commentCreate := &entities.CommentCreate{Body: body}
		_, err := commentRepo.Create(user.ID, article.ID, commentCreate)
		if err != nil {
			t.Fatalf("Failed to create comment: %v", err)
		}
		// Add small delay to ensure different timestamps
		time.Sleep(time.Millisecond)
	}

	// Get comments by article slug
	retrievedComments, err := commentRepo.GetByArticleSlug(article.Slug)
	if err != nil {
		t.Fatalf("Failed to get comments: %v", err)
	}

	// Verify number of comments
	if len(retrievedComments) != len(comments) {
		t.Errorf("Expected %d comments, got %d", len(comments), len(retrievedComments))
	}

	// Verify comments are ordered by created_at ASC
	for i := 0; i < len(retrievedComments)-1; i++ {
		if retrievedComments[i].CreatedAt.After(retrievedComments[i+1].CreatedAt) {
			t.Error("Comments should be ordered by created_at ASC")
		}
	}

	// Verify comment content
	for i, comment := range retrievedComments {
		if comment.Body != comments[i] {
			t.Errorf("Expected comment body %s, got %s", comments[i], comment.Body)
		}
		if comment.Author == nil {
			t.Error("Expected comment to have author information")
		}
	}
}

func TestCommentRepository_GetByID(t *testing.T) {
	// Setup test database
	db, err := database.NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate("../../../migrations"); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Create repositories
	userRepo := NewUserRepository(db)
	articleRepo := NewArticleRepository(db, userRepo)
	commentRepo := NewCommentRepository(db, userRepo)

	// Create test data
	userReg := &entities.UserRegistration{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	user, _ := userRepo.Create(userReg)

	articleCreate := &entities.ArticleCreate{
		Title:       "Test Article",
		Description: "Test description",
		Body:        "Test body",
	}
	article, _ := articleRepo.Create(user.ID, articleCreate)

	commentCreate := &entities.CommentCreate{
		Body: "Test comment",
	}
	createdComment, _ := commentRepo.Create(user.ID, article.ID, commentCreate)

	// Test GetByID
	retrievedComment, err := commentRepo.GetByID(createdComment.ID)
	if err != nil {
		t.Fatalf("Failed to get comment by ID: %v", err)
	}

	if retrievedComment.ID != createdComment.ID {
		t.Errorf("Expected comment ID %d, got %d", createdComment.ID, retrievedComment.ID)
	}
	if retrievedComment.Body != createdComment.Body {
		t.Errorf("Expected comment body %s, got %s", createdComment.Body, retrievedComment.Body)
	}

	// Test non-existent comment
	_, err = commentRepo.GetByID(9999)
	if err == nil {
		t.Error("Expected error for non-existent comment")
	}
}

func TestCommentRepository_Delete(t *testing.T) {
	// Setup test database
	db, err := database.NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate("../../../migrations"); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Create repositories
	userRepo := NewUserRepository(db)
	articleRepo := NewArticleRepository(db, userRepo)
	commentRepo := NewCommentRepository(db, userRepo)

	// Create test data
	userReg := &entities.UserRegistration{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	user, _ := userRepo.Create(userReg)

	articleCreate := &entities.ArticleCreate{
		Title:       "Test Article",
		Description: "Test description",
		Body:        "Test body",
	}
	article, _ := articleRepo.Create(user.ID, articleCreate)

	commentCreate := &entities.CommentCreate{
		Body: "Test comment",
	}
	comment, _ := commentRepo.Create(user.ID, article.ID, commentCreate)

	// Test successful deletion
	err = commentRepo.Delete(comment.ID)
	if err != nil {
		t.Fatalf("Failed to delete comment: %v", err)
	}

	// Verify comment is deleted
	_, err = commentRepo.GetByID(comment.ID)
	if err == nil {
		t.Error("Expected error when getting deleted comment")
	}

	// Test deleting non-existent comment
	err = commentRepo.Delete(9999)
	if err == nil {
		t.Error("Expected error when deleting non-existent comment")
	}
}

func TestCommentRepository_IsAuthor(t *testing.T) {
	// Setup test database
	db, err := database.NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate("../../../migrations"); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Create repositories
	userRepo := NewUserRepository(db)
	articleRepo := NewArticleRepository(db, userRepo)
	commentRepo := NewCommentRepository(db, userRepo)

	// Create test users
	user1Reg := &entities.UserRegistration{
		Username: "user1",
		Email:    "user1@example.com",
		Password: "password123",
	}
	user1, _ := userRepo.Create(user1Reg)

	user2Reg := &entities.UserRegistration{
		Username: "user2",
		Email:    "user2@example.com",
		Password: "password123",
	}
	user2, _ := userRepo.Create(user2Reg)

	// Create test article
	articleCreate := &entities.ArticleCreate{
		Title:       "Test Article",
		Description: "Test description",
		Body:        "Test body",
	}
	article, _ := articleRepo.Create(user1.ID, articleCreate)

	// Create test comment by user1
	commentCreate := &entities.CommentCreate{
		Body: "Test comment",
	}
	comment, _ := commentRepo.Create(user1.ID, article.ID, commentCreate)

	// Test author check
	isAuthor, err := commentRepo.IsAuthor(comment.ID, user1.ID)
	if err != nil {
		t.Fatalf("Failed to check author: %v", err)
	}
	if !isAuthor {
		t.Error("Expected user1 to be the author of the comment")
	}

	// Test non-author check
	isAuthor, err = commentRepo.IsAuthor(comment.ID, user2.ID)
	if err != nil {
		t.Fatalf("Failed to check author: %v", err)
	}
	if isAuthor {
		t.Error("Expected user2 to not be the author of the comment")
	}

	// Test non-existent comment
	isAuthor, err = commentRepo.IsAuthor(9999, user1.ID)
	if err != nil {
		t.Fatalf("Failed to check author for non-existent comment: %v", err)
	}
	if isAuthor {
		t.Error("Expected non-existent comment to return false for author check")
	}
}