package config

const (
	RedirectURI             = "http://localhost:8080/callback"
	Scope                   = "playlist-modify-public"
	CreatePlaylistEndpoint  = "https://api.spotify.com/v1/users/%s/playlists"
	AddTracksEndpoint       = "https://api.spotify.com/v1/playlists/%s/tracks"
	SearchEndpoint          = "https://api.spotify.com/v1/search"
	CreateEndpoint          = "https://api.spotify.com/v1/users/%s/playlists"
	AuthEndpoint            = "https://accounts.spotify.com/authorize"
	UserInfoEndpoint        = "https://api.spotify.com/v1/me"
	TokenEndpoint           = "https://accounts.spotify.com/api/token"
	RecommendationsEndpoint = "https://api.spotify.com/v1/recommendations"
)
