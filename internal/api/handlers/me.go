package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/perezdid/go-mixtape-trading/internal/api/models"
	"github.com/perezdid/go-mixtape-trading/internal/config"
	"github.com/perezdid/go-mixtape-trading/internal/utils"
)

func UserInfo(w http.ResponseWriter, r *http.Request) {
	accessToken, err := utils.GetCookie(r, "access_token")
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

	err = utils.SetCookie(w, "user_id", userProfile.ID)
	if err != nil {
		http.Error(w, "Could not set user config", http.StatusUnauthorized)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
