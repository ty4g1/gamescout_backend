package models

type Screenshot struct {
	ID            int    `json:"id"`
	PathThumbnail string `json:"path_thumbnail"`
	PathFull      string `json:"path_full"`
}

type Video struct {
	Res480 string `json:"480"`
	ResMax string `json:"max"`
}

type Movie struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Thumbnail string `json:"thumbnail"`
	Webm      Video  `json:"webm"`
	Mp4       Video  `json:"mp4"`
	Highlight bool   `json:"highlight"`
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
	Movies      []Movie         `json:"movies"`
}
