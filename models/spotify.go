package models

type Track struct {
	Name       string   `json:"name"`
	Artists    []Artist `json:"artists"`
	Album      Album    `json:"album"`
	PreviewURL string   `json:"preview_url"`
	URI        string   `json:"uri"`
}

type Artist struct {
	Name string `json:"name"`
}

type Album struct {
	Name   string  `json:"name"`
	Images []Image `json:"images"`
}

type Image struct {
	URL string `json:"url"`
}
