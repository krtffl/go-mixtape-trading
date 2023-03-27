package api

import (
	"net/http"

	"github.com/perezdid/go-mixtape-trading/internal/api/handlers"
	"github.com/perezdid/go-mixtape-trading/internal/middleware"
)

func SetupRoutes() {
	http.Handle("/", http.HandlerFunc(handlers.Status))
	http.Handle("/login", middleware.GuardRoute(http.HandlerFunc(handlers.Login)))
	http.Handle("/callback", middleware.GuardRoute(http.HandlerFunc(handlers.Callback)))
	http.Handle("/me", middleware.GuardRoute(http.HandlerFunc(handlers.UserInfo)))
	http.Handle("/search", middleware.GuardRoute(http.HandlerFunc(handlers.Search)))
	http.Handle("/playlist", middleware.GuardRoute(http.HandlerFunc(handlers.CreatePlaylist)))
}
