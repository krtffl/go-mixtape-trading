package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/perezdid/go-mixtape-trading/config"
	"github.com/perezdid/go-mixtape-trading/models"
	"github.com/perezdid/go-mixtape-trading/utils"
)

func CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	accessToken, err := utils.GetCookie(r, "access_token")
	if err != nil {
		http.Error(w, "Invalid access token", http.StatusUnauthorized)
		return
	}

	var requestBody models.PlaylistRequestBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	userID, err := utils.GetCookie(r, "user_id")
	if err != nil {
		http.Error(w, "Invalid access token", http.StatusUnauthorized)
		return
	}

	client := &http.Client{}
	playlistName := fmt.Sprintf("This is %s", requestBody.PlaylistFor)
	playlistData := map[string]string{
		"name":        playlistName,
		"description": "i have used the go mixtape trading application to create this playlist for you",
	}
	playlistBody, _ := json.Marshal(playlistData)
	playlistReq, err := http.NewRequest("POST", fmt.Sprintf(config.CreatePlaylistEndpoint, userID), bytes.NewBuffer(playlistBody))
	if err != nil {
		http.Error(w, "Failed to create playlist request", http.StatusInternalServerError)
		return
	}
	playlistReq.Header.Set("Content-Type", "application/json")
	playlistReq.Header.Set("Authorization", "Bearer "+accessToken)

	playlistResp, err := client.Do(playlistReq)
	if err != nil {
		http.Error(w, "Failed to create playlist", http.StatusInternalServerError)
		return
	}
	defer playlistResp.Body.Close()

	if playlistResp.StatusCode != http.StatusCreated {
		http.Error(w, "Failed to create playlist", http.StatusInternalServerError)
		return
	}

	playlistBody, err = io.ReadAll(playlistResp.Body)
	if err != nil {
		http.Error(w, "Failed to read playlist response body", http.StatusInternalServerError)
		return
	}

	var playlistResponse models.PlaylistResponse
	err = json.Unmarshal(playlistBody, &playlistResponse)
	if err != nil {
		http.Error(w, "Failed to parse playlist response body", http.StatusInternalServerError)
		return
	}

	trackData := map[string][]string{
		"uris": requestBody.Tracks,
	}

	addTracksBody, _ := json.Marshal(trackData)
	addTracksReq, err := http.NewRequest("POST", fmt.Sprintf(config.AddTracksEndpoint, playlistResponse.ID), bytes.NewBuffer(addTracksBody))
	if err != nil {
		http.Error(w, "Failed to create add tracks request", http.StatusInternalServerError)
		return
	}
	addTracksReq.Header.Set("Content-Type", "application/json")
	addTracksReq.Header.Set("Authorization", "Bearer "+accessToken)

	addTracksResp, err := client.Do(addTracksReq)
	if err != nil {
		http.Error(w, "Failed to add tracks to playlist", http.StatusInternalServerError)
		return
	}
	defer addTracksResp.Body.Close()

	if addTracksResp.StatusCode != http.StatusCreated {
		http.Error(w, "Failed to add tracks to playlist", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
