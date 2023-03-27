package models

type PlaylistResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Public bool `json:"public"`
}

type PlaylistCreationResponse struct {
	ID   string `json:"id"`
	Link string `json:"link"`
}

type SearchResponse struct {
	Tracks struct {
		Items []Track `json:"items"`
	} `json:"tracks"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type UserProfile struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Images      []struct {
		URL string `json:"url"`
	} `json:"images"`
	Product string `json:"product"`
}

type RecommendationResponse struct {
	Tracks []struct {
		URI string `json:"uri"`
	} `json:"tracks"`
}
