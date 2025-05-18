package models

type Screenshot struct {
	ID            int    `json:"id"`
	PathThumbnail string `json:"path_thumbnail"`
	PathFull      string `json:"path_full"`
}

type SteamAPIDetailsWrapper map[string]struct {
	Success bool            `json:"success"`
	Data    SteamAPIDetails `json:"data"`
}

type SteamAPIDetails struct {
	ShortDesc   string `json:"short_description"`
	ReleaseDate struct {
		Date string `json:"date"`
	} `json:"release_date"`
	Platforms   map[string]bool `json:"platforms"`
	Thumbnail   string          `json:"header_image"`
	Background  string          `json:"background"`
	Screenshots []Screenshot    `json:"screenshots"`
}
