package config

const (
	RedirectURI = "http://localhost:8080/callback"
	Scope       = "playlist-modify-public"
	Cookie      = "access_token"
)

var (
	AuthEndpoint   = "https://accounts.spotify.com/authorize"
	TokenEndpoint  = "https://accounts.spotify.com/api/token"
	SearchEndpoint = "https://api.spotify.com/v1/search"
	CreateEndpoint = "https://api.spotify.com/v1/users/%s/playlists"
)
