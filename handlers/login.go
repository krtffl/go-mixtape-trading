package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/perezdid/golang-mixtape-trading/config"
)

func Login(w http.ResponseWriter, r *http.Request) {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")

	authURL := fmt.Sprintf("%s?response_type=code&client_id=%s&redirect_uri=%s&scope=%s",
		config.AuthEndpoint, clientID, config.RedirectURI, config.Scope)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)

}
