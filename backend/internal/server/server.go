package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/emotab87/vibe_coding/backend/internal/config"
	"github.com/emotab87/vibe_coding/backend/internal/handlers"
	"github.com/emotab87/vibe_coding/backend/internal/middleware"
)

// Server represents our application server
type Server struct {
	config  *config.Config
	router  *mux.Router
	handler http.Handler
}

// NewServer creates a new server instance with all routes and middleware configured
func NewServer(cfg *config.Config) *Server {
	s := &Server{
		config: cfg,
		router: mux.NewRouter(),
	}

	s.setupRoutes()
	s.setupMiddleware()

	return s
}

// Handler returns the configured HTTP handler
func (s *Server) Handler() http.Handler {
	return s.handler
}

// setupRoutes configures all application routes
func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.HandleFunc("/health", handlers.HealthCheckHandler).Methods("GET")

	// API routes under /api prefix
	api := s.router.PathPrefix("/api").Subrouter()

	// Authentication routes
	api.HandleFunc("/users", handlers.RegisterUserHandler).Methods("POST")
	api.HandleFunc("/users/login", handlers.LoginUserHandler).Methods("POST")

	// Protected routes (require authentication)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware(s.config.JWTSecret))

	protected.HandleFunc("/user", handlers.GetCurrentUserHandler).Methods("GET")
	protected.HandleFunc("/user", handlers.UpdateUserHandler).Methods("PUT")

	// Articles routes
	api.HandleFunc("/articles", handlers.ListArticlesHandler).Methods("GET")
	api.HandleFunc("/articles/{slug}", handlers.GetArticleHandler).Methods("GET")

	// Protected article routes
	protected.HandleFunc("/articles", handlers.CreateArticleHandler).Methods("POST")
	protected.HandleFunc("/articles/{slug}", handlers.UpdateArticleHandler).Methods("PUT")
	protected.HandleFunc("/articles/{slug}", handlers.DeleteArticleHandler).Methods("DELETE")

	// Comments routes
	api.HandleFunc("/articles/{slug}/comments", handlers.ListCommentsHandler).Methods("GET")
	protected.HandleFunc("/articles/{slug}/comments", handlers.CreateCommentHandler).Methods("POST")
	protected.HandleFunc("/articles/{slug}/comments/{id}", handlers.DeleteCommentHandler).Methods("DELETE")

	// Profile routes
	api.HandleFunc("/profiles/{username}", handlers.GetProfileHandler).Methods("GET")

	if s.config.IsDevelopment() {
		log.Printf("🛣️  Routes configured for development environment")
	}
}

// setupMiddleware configures all middleware for the server
func (s *Server) setupMiddleware() {
	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins: parseCORSOrigins(s.config.CORSOrigins),
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
		Debug:            s.config.DebugCORS,
	})

	// Apply middleware stack
	handler := s.router
	handler = middleware.LoggingMiddleware(handler)
	handler = middleware.RecoveryMiddleware(handler)
	handler = c.Handler(handler)

	s.handler = handler

	if s.config.IsDevelopment() {
		log.Printf("🛡️  Middleware configured for development")
		log.Printf("🌐 CORS origins: %s", s.config.CORSOrigins)
	}
}

// parseCORSOrigins parses CORS origins from environment variable
func parseCORSOrigins(origins string) []string {
	if origins == "" {
		return []string{"*"}
	}

	// Split by comma and trim whitespace
	result := make([]string, 0)
	for _, origin := range strings.Split(origins, ",") {
		if trimmed := strings.TrimSpace(origin); trimmed != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return []string{"*"}
	}

	return result
}