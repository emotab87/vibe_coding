package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
)

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// RecoveryMiddleware recovers from panics and returns a 500 error
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with stack trace
				log.Printf("🚨 PANIC: %v\n%s", err, debug.Stack())

				// Return 500 error to client
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				response := ErrorResponse{
					Error: "Internal server error",
				}

				if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
					// If JSON encoding fails, fall back to plain text
					w.Header().Set("Content-Type", "text/plain")
					w.Write([]byte("Internal server error"))
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}