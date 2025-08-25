package entities

import (
	"regexp"
	"strings"
	"time"
	"unicode"
)

// Article represents an article in the system
type Article struct {
	ID          int64     `json:"id"`
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	AuthorID    int64     `json:"-"`
	Author      *User     `json:"author,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	
	// Additional fields for future features
	FavoritesCount int  `json:"favoritesCount"`
	Favorited      bool `json:"favorited"`
}

// ArticleCreate represents article creation request
type ArticleCreate struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Body        string `json:"body"`
}

// ArticleUpdate represents article update request
type ArticleUpdate struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Body        *string `json:"body,omitempty"`
}

// ArticleResponse represents single article API response
type ArticleResponse struct {
	Article Article `json:"article"`
}

// ArticlesResponse represents multiple articles API response
type ArticlesResponse struct {
	Articles      []Article `json:"articles"`
	ArticlesCount int       `json:"articlesCount"`
}

// ArticleListQuery represents query parameters for article listing
type ArticleListQuery struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Author string `json:"author"`
}

// Validate validates article creation data
func (ac *ArticleCreate) Validate() *ValidationErrors {
	var errors []ValidationError

	// Title validation
	if ac.Title == "" {
		errors = append(errors, ValidationError{
			Field:   "title",
			Message: "title is required",
		})
	} else if len(strings.TrimSpace(ac.Title)) < 1 {
		errors = append(errors, ValidationError{
			Field:   "title",
			Message: "title cannot be empty",
		})
	} else if len(ac.Title) > 200 {
		errors = append(errors, ValidationError{
			Field:   "title",
			Message: "title must be less than 200 characters long",
		})
	}

	// Description validation
	if ac.Description == "" {
		errors = append(errors, ValidationError{
			Field:   "description",
			Message: "description is required",
		})
	} else if len(strings.TrimSpace(ac.Description)) < 1 {
		errors = append(errors, ValidationError{
			Field:   "description",
			Message: "description cannot be empty",
		})
	} else if len(ac.Description) > 500 {
		errors = append(errors, ValidationError{
			Field:   "description",
			Message: "description must be less than 500 characters long",
		})
	}

	// Body validation
	if ac.Body == "" {
		errors = append(errors, ValidationError{
			Field:   "body",
			Message: "body is required",
		})
	} else if len(strings.TrimSpace(ac.Body)) < 1 {
		errors = append(errors, ValidationError{
			Field:   "body",
			Message: "body cannot be empty",
		})
	} else if len(ac.Body) > 10000 {
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

// Validate validates article update data
func (au *ArticleUpdate) Validate() *ValidationErrors {
	var errors []ValidationError

	// Title validation (if provided)
	if au.Title != nil {
		title := strings.TrimSpace(*au.Title)
		if title == "" {
			errors = append(errors, ValidationError{
				Field:   "title",
				Message: "title cannot be empty",
			})
		} else if len(*au.Title) > 200 {
			errors = append(errors, ValidationError{
				Field:   "title",
				Message: "title must be less than 200 characters long",
			})
		}
	}

	// Description validation (if provided)
	if au.Description != nil {
		description := strings.TrimSpace(*au.Description)
		if description == "" {
			errors = append(errors, ValidationError{
				Field:   "description",
				Message: "description cannot be empty",
			})
		} else if len(*au.Description) > 500 {
			errors = append(errors, ValidationError{
				Field:   "description",
				Message: "description must be less than 500 characters long",
			})
		}
	}

	// Body validation (if provided)
	if au.Body != nil {
		body := strings.TrimSpace(*au.Body)
		if body == "" {
			errors = append(errors, ValidationError{
				Field:   "body",
				Message: "body cannot be empty",
			})
		} else if len(*au.Body) > 10000 {
			errors = append(errors, ValidationError{
				Field:   "body",
				Message: "body must be less than 10000 characters long",
			})
		}
	}

	if len(errors) > 0 {
		return &ValidationErrors{Errors: errors}
	}
	return nil
}

// ToArticleResponse converts Article to ArticleResponse
func (a *Article) ToArticleResponse() ArticleResponse {
	return ArticleResponse{
		Article: *a,
	}
}

// GenerateSlug generates a URL-friendly slug from title
func GenerateSlug(title string) string {
	if title == "" {
		return ""
	}

	// Convert to lowercase
	slug := strings.ToLower(title)
	
	// Replace spaces and multiple whitespace with single hyphen
	re := regexp.MustCompile(`\s+`)
	slug = re.ReplaceAllString(slug, "-")
	
	// Remove non-alphanumeric characters except hyphens
	var result strings.Builder
	for _, r := range slug {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '-' {
			result.WriteRune(r)
		}
	}
	slug = result.String()
	
	// Remove leading/trailing hyphens
	slug = strings.Trim(slug, "-")
	
	// Replace multiple consecutive hyphens with single hyphen
	re = regexp.MustCompile(`-+`)
	slug = re.ReplaceAllString(slug, "-")
	
	// Limit length to 100 characters
	if len(slug) > 100 {
		slug = slug[:100]
		slug = strings.TrimRight(slug, "-")
	}
	
	return slug
}

// EnsureUniqueSlug ensures slug uniqueness by appending suffix if needed
func EnsureUniqueSlug(baseSlug string, existingSlugs []string) string {
	if baseSlug == "" {
		return ""
	}

	// Check if base slug is unique
	slug := baseSlug
	exists := false
	for _, existing := range existingSlugs {
		if existing == slug {
			exists = true
			break
		}
	}
	
	if !exists {
		return slug
	}
	
	// Try with numeric suffixes
	counter := 1
	for {
		candidateSlug := slug + "-" + intToString(counter)
		exists = false
		
		for _, existing := range existingSlugs {
			if existing == candidateSlug {
				exists = true
				break
			}
		}
		
		if !exists {
			return candidateSlug
		}
		
		counter++
		// Prevent infinite loop
		if counter > 1000 {
			return slug + "-" + intToString(int(time.Now().Unix()))
		}
	}
}

// IsValidSlug checks if a slug is valid format
func IsValidSlug(slug string) bool {
	if slug == "" {
		return false
	}
	
	if len(slug) > 100 {
		return false
	}
	
	// Check for valid characters (lowercase letters, numbers, hyphens)
	re := regexp.MustCompile(`^[a-z0-9-]+$`)
	if !re.MatchString(slug) {
		return false
	}
	
	// Cannot start or end with hyphen
	if strings.HasPrefix(slug, "-") || strings.HasSuffix(slug, "-") {
		return false
	}
	
	// Cannot have consecutive hyphens
	if strings.Contains(slug, "--") {
		return false
	}
	
	return true
}

// Helper function to convert int to string
func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	
	var digits []rune
	negative := n < 0
	if negative {
		n = -n
	}
	
	for n > 0 {
		digit := n % 10
		digits = append([]rune{rune('0'+digit)}, digits...)
		n /= 10
	}
	
	result := string(digits)
	if negative {
		result = "-" + result
	}
	
	return result
}