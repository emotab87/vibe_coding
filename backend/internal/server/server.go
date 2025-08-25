package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/emotab87/vibe_coding/backend/internal/config"
	"github.com/emotab87/vibe_coding/backend/internal/database"
	"github.com/emotab87/vibe_coding/backend/internal/handlers"
	"github.com/emotab87/vibe_coding/backend/internal/middleware"
	"github.com/emotab87/vibe_coding/backend/internal/repositories"
	"github.com/emotab87/vibe_coding/backend/internal/services"
)

// Server represents our application server
type Server struct {
	config      *config.Config
	router      *mux.Router
	handler     http.Handler
	db          *database.DB
	userRepo    repositories.UserRepository
	articleRepo repositories.ArticleRepository
	jwtService  services.JWTService
	authHandlers *handlers.AuthHandlers
	articleHandlers *handlers.ArticleHandlers
}

// NewServer creates a new server instance with all routes and middleware configured
func NewServer(cfg *config.Config) (*Server, error) {
	// Initialize database
	db, err := database.NewDB(cfg.DatabasePath)
	if err != nil {
		return nil, err
	}

	// Run migrations
	if err := db.Migrate("./migrations"); err != nil {
		return nil, err
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	articleRepo := repositories.NewArticleRepository(db, userRepo)

	// Initialize services
	jwtService := services.NewJWTService(cfg.JWTSecret, 24) // 24 hours token expiry

	// Initialize handlers
	authHandlers := handlers.NewAuthHandlers(userRepo, jwtService)
	articleHandlers := handlers.NewArticleHandlers(articleRepo)

	s := &Server{
		config:       cfg,
		router:       mux.NewRouter(),
		db:           db,
		userRepo:     userRepo,
		articleRepo:  articleRepo,
		jwtService:   jwtService,
		authHandlers: authHandlers,
		articleHandlers: articleHandlers,
	}

	s.setupRoutes()
	s.setupMiddleware()

	return s, nil
}

// Handler returns the configured HTTP handler
func (s *Server) Handler() http.Handler {
	return s.handler
}

// Close closes the server and its dependencies
func (s *Server) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// setupRoutes configures all application routes
func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.HandleFunc("/health", handlers.HealthCheckHandler).Methods("GET")

	// API routes under /api prefix
	api := s.router.PathPrefix("/api").Subrouter()

	// Authentication routes
	api.HandleFunc("/users", s.authHandlers.RegisterUser).Methods("POST")
	api.HandleFunc("/users/login", s.authHandlers.LoginUser).Methods("POST")

	// Protected routes (require authentication)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware(s.config.JWTSecret))

	protected.HandleFunc("/user", s.authHandlers.GetCurrentUser).Methods("GET")
	protected.HandleFunc("/user", s.authHandlers.UpdateUser).Methods("PUT")

	// Articles routes
	api.HandleFunc("/articles", s.articleHandlers.ListArticles).Methods("GET")
	api.HandleFunc("/articles/{slug}", s.articleHandlers.GetArticle).Methods("GET")

	// Protected article routes
	protected.HandleFunc("/articles", s.articleHandlers.CreateArticle).Methods("POST")
	protected.HandleFunc("/articles/{slug}", s.articleHandlers.UpdateArticle).Methods("PUT")
	protected.HandleFunc("/articles/{slug}", s.articleHandlers.DeleteArticle).Methods("DELETE")

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