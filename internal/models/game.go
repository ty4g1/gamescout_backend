package models

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

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

// Only create games using this factory method, otherwise feature vector is nil
func NewGame(appId int, name string, shortDesc string, price int, initialPrice int,
	discount int, releaseDate string, genres []string, tags map[string]int,
	positive int, negative int, platforms []string) *Game {

	keys := make([]string, 0, len(tags))
	for k := range tags {
		keys = append(keys, k)
	}

	requestBody := &struct {
		Genres      string `json:"genres"`
		Tags        string `json:"tags"`
		Description string `json:"description"`
	}{
		Genres:      strings.Join(genres, " "),
		Tags:        strings.Join(keys, " "),
		Description: shortDesc,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatalf("Error parsing into JSON: %v", err)
	}

	resp, err := http.Post("https://microservice:8000/vectorize", "application/json", bytes.NewReader(jsonData))
	if err != nil {
		log.Fatalf("Error getting feature vector: %v", err)
	}
	defer resp.Body.Close()

	var featureVector []float32
	if err := json.NewDecoder(resp.Body).Decode(&featureVector); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	game := &Game{
		AppId:         appId,
		Name:          name,
		ShortDesc:     shortDesc,
		Price:         price,
		InitialPrice:  initialPrice,
		Discount:      discount,
		ReleaseDate:   releaseDate,
		Genres:        genres,
		Tags:          tags,
		Positive:      positive,
		Negative:      negative,
		Platforms:     platforms,
		FeatureVector: featureVector,
	}

	return game
}
