package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/perezdid/go-mixtape-trading/handlers"
	"github.com/perezdid/go-mixtape-trading/utils"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = utils.SetEncryptionKeyEnvVar()
	if err != nil {
		log.Fatalf("Error setting encryption key: %v", err)
	}
	http.HandleFunc("/", handlers.Status)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/callback", handlers.Callback)
	http.HandleFunc("/me", handlers.UserInfo)
	http.HandleFunc("/search", handlers.Search)
	http.HandleFunc("/playlist", handlers.CreatePlaylist)

	http.ListenAndServe(":8080", nil)
}
