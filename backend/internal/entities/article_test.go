package entities

import (
	"testing"
)

func TestArticleCreateValidate(t *testing.T) {
	tests := []struct {
		name     string
		article  ArticleCreate
		wantErr  bool
		errorMsg string
	}{
		{
			name: "Valid article",
			article: ArticleCreate{
				Title:       "Test Article",
				Description: "Test description",
				Body:        "Test body content",
			},
			wantErr: false,
		},
		{
			name: "Empty title",
			article: ArticleCreate{
				Title:       "",
				Description: "Test description",
				Body:        "Test body content",
			},
			wantErr:  true,
			errorMsg: "title is required",
		},
		{
			name: "Whitespace only title",
			article: ArticleCreate{
				Title:       "   ",
				Description: "Test description",
				Body:        "Test body content",
			},
			wantErr:  true,
			errorMsg: "title cannot be empty",
		},
		{
			name: "Title too long",
			article: ArticleCreate{
				Title:       generateLongString(201),
				Description: "Test description",
				Body:        "Test body content",
			},
			wantErr:  true,
			errorMsg: "title must be less than 200 characters long",
		},
		{
			name: "Empty description",
			article: ArticleCreate{
				Title:       "Test Article",
				Description: "",
				Body:        "Test body content",
			},
			wantErr:  true,
			errorMsg: "description is required",
		},
		{
			name: "Whitespace only description",
			article: ArticleCreate{
				Title:       "Test Article",
				Description: "   ",
				Body:        "Test body content",
			},
			wantErr:  true,
			errorMsg: "description cannot be empty",
		},
		{
			name: "Description too long",
			article: ArticleCreate{
				Title:       "Test Article",
				Description: generateLongString(501),
				Body:        "Test body content",
			},
			wantErr:  true,
			errorMsg: "description must be less than 500 characters long",
		},
		{
			name: "Empty body",
			article: ArticleCreate{
				Title:       "Test Article",
				Description: "Test description",
				Body:        "",
			},
			wantErr:  true,
			errorMsg: "body is required",
		},
		{
			name: "Whitespace only body",
			article: ArticleCreate{
				Title:       "Test Article",
				Description: "Test description",
				Body:        "   ",
			},
			wantErr:  true,
			errorMsg: "body cannot be empty",
		},
		{
			name: "Body too long",
			article: ArticleCreate{
				Title:       "Test Article",
				Description: "Test description",
				Body:        generateLongString(10001),
			},
			wantErr:  true,
			errorMsg: "body must be less than 10000 characters long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.article.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleCreate.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && len(tt.errorMsg) > 0 {
				found := false
				for _, validationErr := range err.Errors {
					if validationErr.Message == tt.errorMsg {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error message '%s' not found in validation errors: %v", tt.errorMsg, err.Errors)
				}
			}
		})
	}
}

func TestArticleUpdateValidate(t *testing.T) {
	tests := []struct {
		name     string
		article  ArticleUpdate
		wantErr  bool
		errorMsg string
	}{
		{
			name: "Valid update - all fields",
			article: ArticleUpdate{
				Title:       stringPtr("Updated Title"),
				Description: stringPtr("Updated description"),
				Body:        stringPtr("Updated body content"),
			},
			wantErr: false,
		},
		{
			name: "Valid update - partial fields",
			article: ArticleUpdate{
				Title: stringPtr("Updated Title Only"),
			},
			wantErr: false,
		},
		{
			name: "Empty update",
			article: ArticleUpdate{},
			wantErr: false,
		},
		{
			name: "Empty title",
			article: ArticleUpdate{
				Title: stringPtr("   "),
			},
			wantErr:  true,
			errorMsg: "title cannot be empty",
		},
		{
			name: "Title too long",
			article: ArticleUpdate{
				Title: stringPtr(generateLongString(201)),
			},
			wantErr:  true,
			errorMsg: "title must be less than 200 characters long",
		},
		{
			name: "Empty description",
			article: ArticleUpdate{
				Description: stringPtr("   "),
			},
			wantErr:  true,
			errorMsg: "description cannot be empty",
		},
		{
			name: "Description too long",
			article: ArticleUpdate{
				Description: stringPtr(generateLongString(501)),
			},
			wantErr:  true,
			errorMsg: "description must be less than 500 characters long",
		},
		{
			name: "Empty body",
			article: ArticleUpdate{
				Body: stringPtr("   "),
			},
			wantErr:  true,
			errorMsg: "body cannot be empty",
		},
		{
			name: "Body too long",
			article: ArticleUpdate{
				Body: stringPtr(generateLongString(10001)),
			},
			wantErr:  true,
			errorMsg: "body must be less than 10000 characters long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.article.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleUpdate.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && len(tt.errorMsg) > 0 {
				found := false
				for _, validationErr := range err.Errors {
					if validationErr.Message == tt.errorMsg {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error message '%s' not found in validation errors: %v", tt.errorMsg, err.Errors)
				}
			}
		})
	}
}

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		expected string
	}{
		{
			name:     "Simple title",
			title:    "Hello World",
			expected: "hello-world",
		},
		{
			name:     "Title with special characters",
			title:    "Hello, World! How are you?",
			expected: "hello-world-how-are-you",
		},
		{
			name:     "Title with numbers",
			title:    "Article 123 Test",
			expected: "article-123-test",
		},
		{
			name:     "Title with multiple spaces",
			title:    "Hello    World   Test",
			expected: "hello-world-test",
		},
		{
			name:     "Title with leading/trailing spaces",
			title:    "   Hello World   ",
			expected: "hello-world",
		},
		{
			name:     "Empty title",
			title:    "",
			expected: "",
		},
		{
			name:     "Title with only special characters",
			title:    "!@#$%^&*()",
			expected: "",
		},
		{
			name:     "Long title",
			title:    "This is a very long title that should be truncated when converted to slug because it exceeds the maximum length limit",
			expected: "this-is-a-very-long-title-that-should-be-truncated-when-converted-to-slug-because-it-ex",
		},
		{
			name:     "Unicode characters",
			title:    "Hello 世界",
			expected: "hello",
		},
		{
			name:     "Title with hyphens",
			title:    "Pre-existing-hyphens",
			expected: "pre-existing-hyphens",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateSlug(tt.title)
			if result != tt.expected {
				t.Errorf("GenerateSlug(%s) = %s, want %s", tt.title, result, tt.expected)
			}
		})
	}
}

func TestEnsureUniqueSlug(t *testing.T) {
	tests := []struct {
		name          string
		baseSlug      string
		existingSlugs []string
		expected      string
	}{
		{
			name:          "Unique slug",
			baseSlug:      "hello-world",
			existingSlugs: []string{"other-slug", "another-slug"},
			expected:      "hello-world",
		},
		{
			name:          "Duplicate slug - first suffix",
			baseSlug:      "hello-world",
			existingSlugs: []string{"hello-world"},
			expected:      "hello-world-1",
		},
		{
			name:          "Duplicate slug - multiple suffixes",
			baseSlug:      "hello-world",
			existingSlugs: []string{"hello-world", "hello-world-1", "hello-world-2"},
			expected:      "hello-world-3",
		},
		{
			name:          "Empty base slug",
			baseSlug:      "",
			existingSlugs: []string{"some-slug"},
			expected:      "",
		},
		{
			name:          "No existing slugs",
			baseSlug:      "hello-world",
			existingSlugs: []string{},
			expected:      "hello-world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EnsureUniqueSlug(tt.baseSlug, tt.existingSlugs)
			if result != tt.expected {
				t.Errorf("EnsureUniqueSlug(%s, %v) = %s, want %s", tt.baseSlug, tt.existingSlugs, result, tt.expected)
			}
		})
	}
}

func TestIsValidSlug(t *testing.T) {
	tests := []struct {
		name     string
		slug     string
		expected bool
	}{
		{
			name:     "Valid slug",
			slug:     "hello-world",
			expected: true,
		},
		{
			name:     "Valid slug with numbers",
			slug:     "article-123",
			expected: true,
		},
		{
			name:     "Empty slug",
			slug:     "",
			expected: false,
		},
		{
			name:     "Slug with uppercase",
			slug:     "Hello-World",
			expected: false,
		},
		{
			name:     "Slug with spaces",
			slug:     "hello world",
			expected: false,
		},
		{
			name:     "Slug starting with hyphen",
			slug:     "-hello-world",
			expected: false,
		},
		{
			name:     "Slug ending with hyphen",
			slug:     "hello-world-",
			expected: false,
		},
		{
			name:     "Slug with consecutive hyphens",
			slug:     "hello--world",
			expected: false,
		},
		{
			name:     "Slug too long",
			slug:     generateLongString(101),
			expected: false,
		},
		{
			name:     "Slug with special characters",
			slug:     "hello@world",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidSlug(tt.slug)
			if result != tt.expected {
				t.Errorf("IsValidSlug(%s) = %v, want %v", tt.slug, result, tt.expected)
			}
		})
	}
}

func TestArticleToArticleResponse(t *testing.T) {
	article := &Article{
		ID:          1,
		Slug:        "test-article",
		Title:       "Test Article",
		Description: "Test description",
		Body:        "Test body",
		AuthorID:    1,
	}

	response := article.ToArticleResponse()

	if response.Article.ID != article.ID {
		t.Errorf("Expected Article.ID %d, got %d", article.ID, response.Article.ID)
	}
	if response.Article.Slug != article.Slug {
		t.Errorf("Expected Article.Slug %s, got %s", article.Slug, response.Article.Slug)
	}
	if response.Article.Title != article.Title {
		t.Errorf("Expected Article.Title %s, got %s", article.Title, response.Article.Title)
	}
}

func TestIntToString(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{
			name:     "Zero",
			input:    0,
			expected: "0",
		},
		{
			name:     "Positive number",
			input:    123,
			expected: "123",
		},
		{
			name:     "Negative number",
			input:    -456,
			expected: "-456",
		},
		{
			name:     "Single digit",
			input:    7,
			expected: "7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := intToString(tt.input)
			if result != tt.expected {
				t.Errorf("intToString(%d) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

// Helper function to generate a long string of specified length
func generateLongString(length int) string {
	if length <= 0 {
		return ""
	}
	result := make([]byte, length)
	for i := range result {
		result[i] = 'a'
	}
	return string(result)
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}