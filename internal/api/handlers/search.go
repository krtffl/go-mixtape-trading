package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/perezdid/go-mixtape-trading/internal/api/models"
	"github.com/perezdid/go-mixtape-trading/internal/config"
	"github.com/perezdid/go-mixtape-trading/internal/utils"
)

func Search(w http.ResponseWriter, r *http.Request) {
	log.Printf("received %s request at %s from %s", r.Method, r.URL, r.RemoteAddr)

	accessToken, err := utils.GetCookie(r, "access_token")
	if err != nil {
		log.Printf("retrieving access token from cookie failed: %s", err)
		http.Error(w, "invalid access token", http.StatusUnauthorized)
		return
	}

	log.Printf("retrieving track name from query: %s", r.URL.Query())

	query := r.URL.Query().Get("track")
	if query == "" {
		log.Printf("retrieving track name from query failed: %s", err)
		http.Error(w, "missing track name", http.StatusBadRequest)
		return
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s?q=%s&type=track&limit=%s", config.SearchEndpoint, query, "5"), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	log.Printf("request to spotify search api: %s", req.URL)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("request to spotify search api failed: %s", err)
		http.Error(w, "could not find spotify tracks", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	var searchResponse models.SearchResponse
	err = json.NewDecoder(resp.Body).Decode(&searchResponse)

	if err != nil {
		log.Printf("parsing track search response failed: %s", err)
		http.Error(w, "failed to parse track search response", http.StatusInternalServerError)
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

	log.Printf("response: %s", tracks)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tracks)
}
