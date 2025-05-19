package models

type Game struct {
	AppId         int
	Name          string
	ShortDesc     string
	Price         int
	InitialPrice  int
	Discount      int
	ReleaseDate   string
	Genres        []string
	Tags          map[string]int
	Positive      int
	Negative      int
	Platforms     []string
	FeatureVector []float32
}
