package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/perezdid/go-mixtape-trading/internal/config"
)

func Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("received %s request at %s from %s", r.Method, r.URL, r.RemoteAddr)

	clientID := os.Getenv("SPOTIFY_CLIENT_ID")

	authURL := fmt.Sprintf("%s?response_type=code&client_id=%s&redirect_uri=%s&scope=%s",
		config.AuthEndpoint, clientID, fmt.Sprintf(config.RedirectURI, os.Getenv("BASE_URL")), config.Scope)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)

}
