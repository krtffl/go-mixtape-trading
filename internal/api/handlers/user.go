package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/perezdid/go-mixtape-trading/internal/api/models"
	"github.com/perezdid/go-mixtape-trading/internal/config"
	"github.com/perezdid/go-mixtape-trading/internal/utils"
)

func UserInfo(w http.ResponseWriter, r *http.Request) {
	log.Printf("received %s request at %s from %s", r.Method, r.URL, r.RemoteAddr)

	accessToken, err := utils.GetCookie(r, "access_token")
	if err != nil {
		log.Printf("retrieving access token from cookie failed: %s", err)
		http.Error(w, "invalid access token", http.StatusUnauthorized)
		return
	}

	req, err := http.NewRequest("GET", config.UserInfoEndpoint, nil)
	if err != nil {
		log.Printf("creating user info request failed: %s", err)
		http.Error(w, "failed to create user profile request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("retrieving user information failed: %s", err)
		http.Error(w, "failed to retrieve user profile", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("retrieving user information failed - status not ok: %s", err)
		http.Error(w, "failed to retrieve user profile", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("reading user info response failed: %s", err)
		http.Error(w, "failed to read user profile response body", http.StatusInternalServerError)
		return
	}

	var userProfile models.UserProfile
	err = json.Unmarshal(body, &userProfile)
	if err != nil {
		log.Printf("parsing user info response failed: %s", err)
		http.Error(w, "failed to parse user profile response body", http.StatusInternalServerError)
		return
	}

	err = utils.SetCookie(w, "user_id", userProfile.ID)

	if err != nil {
		log.Printf("setting user id in cookie failed: %s", err)
		http.Error(w, "could not set user id", http.StatusUnauthorized)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
