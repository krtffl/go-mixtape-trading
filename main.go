package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/perezdid/go-mixtape-trading/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/", handlers.Status)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/callback", handlers.Callback)
	http.HandleFunc("/me", handlers.UserInfo)
	http.HandleFunc("/search", handlers.Search)
	http.HandleFunc("/playlist", handlers.CreatePlaylist)

	http.ListenAndServe(":8080", nil)
}
