package http

import (
	"encoding/json"
	"net/http"

	"github.com/afonsodemori/afonsodev-api/internal/config"
	"github.com/afonsodemori/afonsodev-api/internal/contact"
)

// NewRouter creates the application HTTP handler with global middleware and routes.
func NewRouter(allowedOrigins []string, contactHandler *contact.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", redirectHandler)
	mux.HandleFunc("/health.json", healthHandler)
	// TODO: Remove /send-email route
	mux.HandleFunc("/send-email", contactHandler.HandleSendEmail)
	mux.HandleFunc("/contact", contactHandler.HandleSendEmail)

	return corsMiddleware(allowedOrigins, mux)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	health := struct {
		Status  string `json:"status"`
		Commit  string `json:"commit"`
		BuiltAt string `json:"built_at"`
	}{
		Status:  "up",
		Commit:  config.Commit,
		BuiltAt: config.BuiltAt,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(health); err != nil {
		http.Error(w, "failed to encode health status", http.StatusInternalServerError)
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://afonso.dev", http.StatusTemporaryRedirect)
}

func corsMiddleware(allowedOrigins []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if len(allowedOrigins) > 0 {
			for _, ao := range allowedOrigins {
				if origin == ao {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
