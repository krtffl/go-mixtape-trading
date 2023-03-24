package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/perezdid/golang-mixtape-trading/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/", handlers.Status)
	http.HandleFunc("/search", handlers.Search)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/callback", handlers.Callback)
	http.HandleFunc("/playlist", handlers.Playlist)

	http.ListenAndServe(":8080", nil)
}
