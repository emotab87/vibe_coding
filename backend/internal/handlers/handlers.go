package handlers

import (
	"encoding/json"
	"net/http"
)

// Temporary stub handlers - to be implemented in future issues

// User authentication handlers
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "User registration not yet implemented")
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "User login not yet implemented")
}

func GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "Get current user not yet implemented")
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "Update user not yet implemented")
}

// Article handlers
func ListArticlesHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "List articles not yet implemented")
}

func GetArticleHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "Get article not yet implemented")
}

func CreateArticleHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "Create article not yet implemented")
}

func UpdateArticleHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "Update article not yet implemented")
}

func DeleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "Delete article not yet implemented")
}

// Comment handlers
func ListCommentsHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "List comments not yet implemented")
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "Create comment not yet implemented")
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "Delete comment not yet implemented")
}

// Profile handlers
func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "Get profile not yet implemented")
}

// Helper function to return "not implemented" responses
func writeNotImplemented(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)

	response := map[string]string{
		"error":   "Not implemented",
		"message": message,
	}

	json.NewEncoder(w).Encode(response)
}