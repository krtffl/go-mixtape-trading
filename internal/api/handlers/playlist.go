package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/perezdid/go-mixtape-trading/internal/api/models"
	"github.com/perezdid/go-mixtape-trading/internal/config"
	"github.com/perezdid/go-mixtape-trading/internal/utils"
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

	tracks, err := getRecommendations(requestBody.Tracks, accessToken)
	if err != nil {
		http.Error(w, "Failed to generate track recommendations", http.StatusInternalServerError)
		return
	}

	trackData := map[string][]string{
		"uris": tracks,
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

func getRecommendations(songURIs []string, accessToken string) ([]string, error) {
	limit := 50 - len(songURIs)

	songsIDs, err := getTrackIDsFromURIs(songURIs)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}

	seedURIs := strings.Join(songsIDs, ",")
	requestURL := fmt.Sprintf("%s?limit=%d&seed_tracks=%s&min_energy=0.5&min_popularity=40", config.RecommendationsEndpoint, limit, seedURIs)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var recommendationResponse models.RecommendationResponse
	err = json.Unmarshal(body, &recommendationResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %v", err)
	}

	for _, track := range recommendationResponse.Tracks {
		songURIs = append(songURIs, track.URI)
	}

	return songURIs, nil
}

func getTrackIDsFromURIs(trackURIs []string) ([]string, error) {
	const uriPrefix = "spotify:track:"
	var trackIDs []string
	for _, trackURI := range trackURIs {
		if !strings.HasPrefix(trackURI, uriPrefix) {
			return nil, fmt.Errorf("%s is not a valid track URI", trackURI)
		}
		trackIDs = append(trackIDs, strings.TrimPrefix(trackURI, uriPrefix))
	}
	return trackIDs, nil
}
