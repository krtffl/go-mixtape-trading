package api

import (
	"net/http"

	"github.com/perezdid/go-mixtape-trading/internal/api/handlers"
)

func SetupRoutes() {
	http.HandleFunc("/", handlers.Status)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/callback", handlers.Callback)
	http.HandleFunc("/me", handlers.UserInfo)
	http.HandleFunc("/search", handlers.Search)
	http.HandleFunc("/playlist", handlers.CreatePlaylist)
}
