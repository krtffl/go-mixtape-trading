package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/perezdid/go-mixtape-trading/internal/api"
	"github.com/perezdid/go-mixtape-trading/internal/utils"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = utils.SetEncryptionKeyEnvVar()
	if err != nil {
		log.Fatalf("Error setting encryption key: %v", err)
	}

	api.SetupRoutes()

	http.ListenAndServe(":8080", nil)
}
