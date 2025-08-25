package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/emotab87/vibe_coding/backend/internal/config"
	"github.com/emotab87/vibe_coding/backend/internal/server"
)

func main() {
	// Load configuration from environment variables
	cfg := config.LoadConfig()

	// Create and configure the server
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to create server: %v", err)
	}
	defer srv.Close()

	// Create HTTP server with configured settings
	httpServer := &http.Server{
		Addr:         cfg.ServerAddress(),
		Handler:      srv.Handler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("🚀 Server starting on %s", cfg.ServerAddress())
		log.Printf("📖 Environment: %s", cfg.Environment)
		log.Printf("🔧 Database: %s", cfg.DatabasePath)
		
		serverErrors <- httpServer.ListenAndServe()
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("❌ Server failed to start: %v", err)

	case sig := <-shutdown:
		log.Printf("🔄 Server shutting down due to signal: %v", sig)

		// Give outstanding requests 30 seconds to complete
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("⚠️  Graceful shutdown failed, forcing shutdown: %v", err)
			if err := httpServer.Close(); err != nil {
				log.Printf("❌ Force shutdown failed: %v", err)
			}
		}

		log.Println("✅ Server shutdown complete")
	}
}