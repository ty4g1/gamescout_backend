package models

type GameMedia struct {
	AppID         int
	ThumbnailURL  string
	BackgroundURL string
	Screenshots   []Screenshot
	Movies        []Movie
}
