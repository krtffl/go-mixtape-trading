package middleware

import (
	"net/http"
	"os"
)

type Middleware func(http.Handler) http.Handler

func GuardRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")

		if apiKey == "" || apiKey != os.Getenv("API_KEY") {
			http.Error(w, "you're not authorized to perform this action", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
