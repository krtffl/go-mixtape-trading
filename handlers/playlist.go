package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/perezdid/golang-mixtape-trading/config"
	"github.com/perezdid/golang-mixtape-trading/models"
)

func Playlist(w http.ResponseWriter, r *http.Request) {
	accessToken := r.URL.Query().Get("access_token")

	playlistRequest := models.PlaylistRequest{
		Name:        "My Playlist",
		Public:      true,
		Description: "A playlist created using the Spotify API",
	}

	reqBody, err := json.Marshal(playlistRequest)
	if err != nil {
		http.Error(w, "Failed to create playlist", http.StatusInternalServerError)
		return
	}

	req, _ := http.NewRequest("POST", fmt.Sprintf(config.CreateEndpoint, "me"), strings.NewReader(string(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to create playlist", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var playlistResponse models.PlaylistResponse
	err = json.NewDecoder(resp.Body).Decode(&playlistResponse)
	if err != nil {
		http.Error(w, "Failed to parse playlist creation response", http.StatusInternalServerError)
		return
	}

	playlistID := playlistResponse.ID
	fmt.Fprintf(w, "Created playlist with ID %s", playlistID)
}
