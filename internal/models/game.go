package models

import "time"

type Game struct {
	AppId         int
	Name          string
	ShortDesc     string
	Price         int
	InitialPrice  int
	Discount      int
	ReleaseDate   time.Time
	Genres        []string
	Tags          map[string]int
	Positive      int
	Negative      int
	Platforms     []string
	FeatureVector []float64
}
