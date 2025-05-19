package services

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/ty4g1/gamescout_backend/internal/config"
)

type Vectorizer struct {
	MicroserviceURL string
}

func NewVectorizer(cfg *config.Config) *Vectorizer {
	return &Vectorizer{
		MicroserviceURL: cfg.MicroserviceURL,
	}
}

func (v *Vectorizer) Vectorize(genres string, tags map[string]int, shortDesc string) []float32 {
	keys := make([]string, 0, len(tags))
	for k := range tags {
		keys = append(keys, k)
	}

	requestBody := &struct {
		Genres      string `json:"genres"`
		Tags        string `json:"tags"`
		Description string `json:"description"`
	}{
		Genres:      genres,
		Tags:        strings.Join(keys, " "),
		Description: shortDesc,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Error parsing into JSON: %v", err)
		return make([]float32, 300)
	}

	resp, err := http.Post(v.MicroserviceURL, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		log.Printf("Error getting feature vector: %v", err)
		return make([]float32, 300)
	}
	defer resp.Body.Close()

	var serviceResponse struct {
		Vector []float32 `json:"vector"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&serviceResponse); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return make([]float32, 300)
	}

	return serviceResponse.Vector
}
