package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/perezdid/go-mixtape-trading/config"
	"github.com/perezdid/go-mixtape-trading/models"
)

func UserInfo(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("mixtape_trading")
	if err != nil {
		http.Error(w, "Access token not found", http.StatusUnauthorized)
		return
	}

	accessToken, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid access token", http.StatusUnauthorized)
		return
	}

	req, err := http.NewRequest("GET", config.UserInfoEndpoint, nil)
	if err != nil {
		http.Error(w, "Failed to create user profile request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to retrieve user profile", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to retrieve user profile", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read user profile response body", http.StatusInternalServerError)
		return
	}

	var userProfile models.UserProfile
	err = json.Unmarshal(body, &userProfile)
	if err != nil {
		http.Error(w, "Failed to parse user profile response body", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "user_id",
		Value: url.QueryEscape(userProfile.ID),
	})

	http.Redirect(w, r, "/", http.StatusFound)
}
