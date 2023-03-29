package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/perezdid/go-mixtape-trading/internal/api/models"
	"github.com/perezdid/go-mixtape-trading/internal/config"
	"github.com/perezdid/go-mixtape-trading/internal/utils"
)

func Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("received %s request at %s from %s", r.Method, r.URL, r.RemoteAddr)

	clientID := os.Getenv("SPOTIFY_CLIENT_ID")

	authURL := fmt.Sprintf("%s?response_type=code&client_id=%s&redirect_uri=%s&scope=%s",
		config.AuthEndpoint, clientID, fmt.Sprintf(config.RedirectURI, os.Getenv("BASE_URL")), config.Scope)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)

}

func Callback(w http.ResponseWriter, r *http.Request) {
	log.Printf("received %s request at %s from %s", r.Method, r.URL, r.RemoteAddr)

	code := r.URL.Query().Get("code")
	data := url.Values{}

	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", fmt.Sprintf(config.RedirectURI, os.Getenv("BASE_URL")))

	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	log.Printf("exchanging app key for access token")

	req, _ := http.NewRequest("POST", config.TokenEndpoint, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("exchanging app key for access token failed: %s", err)
		http.Error(w, "failed to retrieve access token", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	log.Printf("parsing access token from request")

	var tokenResponse models.TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		log.Printf("parsing access token from request failed: %s", err)
		http.Error(w, "failed to parse access token response", http.StatusInternalServerError)
		return
	}

	log.Printf("access token retrieved")

	accessToken := tokenResponse.AccessToken

	utils.SetCookie(w, "access_token", accessToken)
	http.Redirect(w, r, "/me", http.StatusTemporaryRedirect)
}
