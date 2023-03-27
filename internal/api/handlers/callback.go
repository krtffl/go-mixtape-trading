package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/perezdid/go-mixtape-trading/internal/api/models"
	"github.com/perezdid/go-mixtape-trading/internal/config"
	"github.com/perezdid/go-mixtape-trading/internal/utils"
)

func Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	data := url.Values{}

	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", config.RedirectURI)

	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	req, _ := http.NewRequest("POST", config.TokenEndpoint, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to retrieve access token", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var tokenResponse models.TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		http.Error(w, "Failed to parse access token response", http.StatusInternalServerError)
		return
	}

	accessToken := tokenResponse.AccessToken

	utils.SetCookie(w, "access_token", accessToken)
	http.Redirect(w, r, "/me", http.StatusTemporaryRedirect)
}
