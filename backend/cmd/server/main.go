package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/ericsantos/pokercards/internal/handler"
	"github.com/ericsantos/pokercards/internal/session"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	store := session.NewStore()
	h := handler.New(store)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	slog.Info("server starting", "addr", addr)
	if err := http.ListenAndServe(addr, corsMiddleware(mux)); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}

// corsMiddleware adds CORS headers so the Vite dev server can reach the API.
// In production the nginx proxy handles this instead.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
