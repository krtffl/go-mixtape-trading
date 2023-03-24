package models

type PlaylistRequest struct {
	Name        string `json:"name"`
	Public      bool   `json:"public"`
	Description string `json:"description"`
}
