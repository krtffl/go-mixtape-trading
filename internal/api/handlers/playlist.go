package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/perezdid/go-mixtape-trading/internal/api/models"
	"github.com/perezdid/go-mixtape-trading/internal/config"
	"github.com/perezdid/go-mixtape-trading/internal/utils"
)

func CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	log.Printf("received %s request at %s from %s", r.Method, r.URL, r.RemoteAddr)

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	accessToken, err := utils.GetCookie(r, "access_token")
	if err != nil {
		log.Printf("retrieving access token from cookie failed: %s", err)
		http.Error(w, "invalid access token", http.StatusUnauthorized)
		return
	}

	log.Printf("parsing create playlist request body")

	var requestBody models.PlaylistRequestBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Printf("parsing create playlist request body failed: %s", err)
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	log.Printf("request body - PlaylistFor: %s", requestBody.PlaylistFor)
	log.Printf("request body - Tracks: %s", requestBody.Tracks)

	log.Printf("retrieving user ID from cookie")

	userID, err := utils.GetCookie(r, "user_id")
	if err != nil {
		log.Printf("retrieving user ID from cookie failed: %s", err)
		http.Error(w, "invalid access token", http.StatusUnauthorized)
		return
	}

	client := &http.Client{}

	log.Printf("creating playlist for user %s", userID)

	playlist, err := createPlaylist(w, client, requestBody.PlaylistFor, userID, accessToken)
	if err != nil {
		log.Printf("creating playlist failed: %s", err)
		http.Error(w, "failed to create playlist", http.StatusInternalServerError)
		return
	}

	log.Printf("generating track recommendations")

	tracks, err := getRecommendations(requestBody.Tracks, accessToken)
	if err != nil {
		log.Printf("generating track recommendations failed: %s", err)
		http.Error(w, "failed to generate track recommendations", http.StatusInternalServerError)
		return
	}

	log.Printf("obtained track recommendations: %s", tracks)
	log.Printf("addings recommended tracks into playlist %s", playlist.ID)

	err = addTracksToPlaylist(w, client, tracks, playlist.ID, accessToken)
	if err != nil {
		log.Printf("addings recommended tracks into playlit failed: %s", err)
		http.Error(w, "failed to add tracks to new playlist", http.StatusInternalServerError)
		return
	}

	log.Printf("tracks added to playlist %s for user %s", playlist.ID, userID)

	jsonResponse, err := json.Marshal(playlist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

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

func createPlaylist(w http.ResponseWriter, client *http.Client, name string, userID string, accessToken string) (*models.PlaylistCreationResponse, error) {
	playlistName := fmt.Sprintf("This is %s", name)
	playlistData := map[string]string{
		"name":        playlistName,
		"description": "i have used the go mixtape trading application to create this playlist for you",
	}

	playlistBody, _ := json.Marshal(playlistData)
	playlistReq, err := http.NewRequest("POST", fmt.Sprintf(config.CreatePlaylistEndpoint, userID), bytes.NewBuffer(playlistBody))

	if err != nil {
		http.Error(w, "Failed to create playlist request", http.StatusInternalServerError)
		return nil, err
	}

	playlistReq.Header.Set("Content-Type", "application/json")
	playlistReq.Header.Set("Authorization", "Bearer "+accessToken)

	playlistResp, err := client.Do(playlistReq)
	if err != nil {
		http.Error(w, "Failed to create playlist", http.StatusInternalServerError)
		return nil, err
	}

	defer playlistResp.Body.Close()

	if playlistResp.StatusCode != http.StatusCreated {
		http.Error(w, "Failed to create playlist", http.StatusInternalServerError)
		return nil, err
	}

	playlistBody, err = io.ReadAll(playlistResp.Body)
	if err != nil {
		http.Error(w, "Failed to read playlist response body", http.StatusInternalServerError)
		return nil, err
	}

	var playlistResponse models.PlaylistResponse
	err = json.Unmarshal(playlistBody, &playlistResponse)
	if err != nil {
		http.Error(w, "Failed to parse playlist response body", http.StatusInternalServerError)
		return nil, err
	}

	response := &models.PlaylistCreationResponse{
		ID:   playlistResponse.ID,
		Link: playlistResponse.ExternalURLs.Spotify,
	}

	return response, nil
}

func addTracksToPlaylist(w http.ResponseWriter, client *http.Client, tracks []string, playlistID string, accessToken string) error {
	trackData := map[string][]string{
		"uris": tracks,
	}

	addTracksBody, _ := json.Marshal(trackData)
	addTracksReq, err := http.NewRequest("POST", fmt.Sprintf(config.AddTracksEndpoint, playlistID), bytes.NewBuffer(addTracksBody))
	if err != nil {
		http.Error(w, "failed to create add tracks request", http.StatusInternalServerError)
		return err
	}

	addTracksReq.Header.Set("Content-Type", "application/json")
	addTracksReq.Header.Set("Authorization", "Bearer "+accessToken)

	addTracksResp, err := client.Do(addTracksReq)
	if err != nil {
		http.Error(w, "failed to add tracks to playlist", http.StatusInternalServerError)
		return err
	}

	defer addTracksResp.Body.Close()

	if addTracksResp.StatusCode != http.StatusCreated {
		http.Error(w, "failed to add tracks to playlist", http.StatusInternalServerError)
		return err
	}

	return nil
}
