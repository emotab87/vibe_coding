package entities

import (
	"testing"
)

func TestCommentCreateValidate(t *testing.T) {
	tests := []struct {
		name     string
		comment  CommentCreate
		wantErr  bool
		errorMsg string
	}{
		{
			name: "Valid comment",
			comment: CommentCreate{
				Body: "This is a test comment",
			},
			wantErr: false,
		},
		{
			name: "Empty body",
			comment: CommentCreate{
				Body: "",
			},
			wantErr:  true,
			errorMsg: "body is required",
		},
		{
			name: "Whitespace only body",
			comment: CommentCreate{
				Body: "   ",
			},
			wantErr:  true,
			errorMsg: "body cannot be empty",
		},
		{
			name: "Body too long",
			comment: CommentCreate{
				Body: generateLongString(10001),
			},
			wantErr:  true,
			errorMsg: "body must be less than 10000 characters long",
		},
		{
			name: "Body at maximum length",
			comment: CommentCreate{
				Body: generateLongString(10000),
			},
			wantErr: false,
		},
		{
			name: "Body with newlines",
			comment: CommentCreate{
				Body: "This is a comment\nwith multiple\nlines",
			},
			wantErr: false,
		},
		{
			name: "Body with special characters",
			comment: CommentCreate{
				Body: "Comment with special chars: @#$%^&*()",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.comment.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CommentCreate.Validate() error = %v, wantErr %v", err, tt.wantErr)
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

func TestCommentToCommentResponse(t *testing.T) {
	comment := &Comment{
		ID:   1,
		Body: "Test comment",
		AuthorID: 1,
	}

	response := comment.ToCommentResponse()

	if response.Comment.ID != comment.ID {
		t.Errorf("Expected Comment.ID %d, got %d", comment.ID, response.Comment.ID)
	}
	if response.Comment.Body != comment.Body {
		t.Errorf("Expected Comment.Body %s, got %s", comment.Body, response.Comment.Body)
	}
	if response.Comment.AuthorID != comment.AuthorID {
		t.Errorf("Expected Comment.AuthorID %d, got %d", comment.AuthorID, response.Comment.AuthorID)
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