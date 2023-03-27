package models

type PlaylistRequestBody struct {
	PlaylistFor string   `json:"playlistFor"`
	Tracks      []string `json:"tracks"`
}
