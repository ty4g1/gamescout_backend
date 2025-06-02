package services

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
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

func (v *Vectorizer) Vectorize(genres string, tags map[string]int, shortDesc string) []float64 {
	VECTOR_DIM, err := strconv.Atoi(os.Getenv("VECTOR_DIM"))
	if err != nil {
		VECTOR_DIM = 384
	}

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
		log.Printf("Error parsing into JSON: %v\n", err)
		return make([]float64, VECTOR_DIM)
	}

	resp, err := http.Post(v.MicroserviceURL, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		log.Printf("Error getting feature vector: %v\n", err)
		return make([]float64, VECTOR_DIM)
	}
	defer resp.Body.Close()

	var serviceResponse struct {
		Vector []float64 `json:"vector"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&serviceResponse); err != nil {
		log.Printf("Error parsing JSON: %v\n", err)
		return make([]float64, VECTOR_DIM)
	}

	return serviceResponse.Vector
}
