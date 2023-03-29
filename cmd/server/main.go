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
		log.Fatal("error loading environment variables")
	}

	err = utils.SetEncryptionKeyEnvVar()
	if err != nil {
		log.Fatalf("%v", err)
	}

	api.SetupRoutes()

	log.Printf("listening on port 8080!")
	http.ListenAndServe(":8080", nil)

}
