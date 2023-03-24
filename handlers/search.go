package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/perezdid/go-mixtape-trading/config"
	"github.com/perezdid/go-mixtape-trading/models"
)

func Search(w http.ResponseWriter, r *http.Request) {
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

	query := r.URL.Query().Get("track")
	if query == "" {
		http.Error(w, "Track name not found", http.StatusBadRequest)
		return
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s?q=%s&type=track&limit=%s", config.SearchEndpoint, query, "5"), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to search tracks", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var searchResponse models.SearchResponse
	err = json.NewDecoder(resp.Body).Decode(&searchResponse)
	if err != nil {
		http.Error(w, "Failed to parse track search response", http.StatusInternalServerError)
		return
	}

	tracks := []map[string]string{}
	for _, item := range searchResponse.Tracks.Items {
		track := map[string]string{
			"name":        item.Name,
			"artist":      item.Artists[0].Name,
			"album":       item.Album.Name,
			"image":       item.Album.Images[0].URL,
			"preview_url": item.PreviewURL,
			"uri":         item.URI,
		}
		tracks = append(tracks, track)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tracks)
}
