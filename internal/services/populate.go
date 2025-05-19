package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ty4g1/gamescout_backend/internal/config"
	"github.com/ty4g1/gamescout_backend/internal/models"
)

type Populator struct {
	Pages                 int
	SteamSpyURLFormat     string
	SteamSpyDetailsFormat string
	SteamAPIDetailsFormat string
	Vectorizer            *Vectorizer
}

func NewPopulator(cfg *config.Config) *Populator {
	return &Populator{
		Pages:                 cfg.SteamSpyPages,
		SteamSpyURLFormat:     cfg.SteamSpyURLFormat,
		SteamSpyDetailsFormat: cfg.SteamSpyDetailsFormat,
		SteamAPIDetailsFormat: cfg.SteamAPIDetailsFormat,
		Vectorizer:            NewVectorizer(cfg),
	}
}

func (p *Populator) Populate() {
	// Number of pages (with 1000 games each) to retrieve
	const PAGES = 10

	for i := range PAGES {
		resp, err := http.Get(fmt.Sprintf(p.SteamSpyURLFormat, i))
		if err != nil {
			log.Printf("Error making request: %v\n", err)
			continue
		}

		var games map[string]models.SteamspyResponse
		if err := json.NewDecoder(resp.Body).Decode(&games); err != nil {
			log.Printf("Error reading and parsing JSON: %v\n", err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		for appIDKey, game := range games {
			appID := game.AppID // Use the AppID from the struct

			// Respect rate limits (1 request per second)
			time.Sleep(1 * time.Second)

			// Get Steam Spy details
			respSpy, err := http.Get(fmt.Sprintf(p.SteamSpyDetailsFormat, appID))
			if err != nil {
				log.Printf("Error making Steam Spy details request for appID %d: %v\n", appID, err)
				continue
			}

			var gameDetailsSpy models.SteamspyDetails
			if err := json.NewDecoder(respSpy.Body).Decode(&gameDetailsSpy); err != nil {
				log.Printf("Error parsing Steam Spy details for appID %d: %v\n", appID, err)
				respSpy.Body.Close()
				continue
			}
			respSpy.Body.Close()
			fmt.Printf("Steam Spy details: %+v\n", gameDetailsSpy)

			// Get Steam API details
			respAPI, err := http.Get(fmt.Sprintf(p.SteamAPIDetailsFormat, appID))
			if err != nil {
				log.Printf("Error making Steam API details request for appID %d: %v\n", appID, err)
				continue
			}

			var gameDetailsAPI models.SteamAPIDetailsWrapper
			if err := json.NewDecoder(respAPI.Body).Decode(&gameDetailsAPI); err != nil {
				log.Printf("Error parsing Steam API details for appID %d: %v\n", appID, err)
				respAPI.Body.Close()
				continue
			}
			respAPI.Body.Close()

			fmt.Println(createGameEntry(game, gameDetailsSpy, gameDetailsAPI, appIDKey, p.Vectorizer))

			break // Only process one game for testing
		}

		// Temporary, for testing with just 1 page
		break
	}
}

func createGameEntry(game models.SteamspyResponse, gameDetailsSpy models.SteamspyDetails, gameDetailsAPI models.SteamAPIDetailsWrapper, appIDKey string, vectorizer *Vectorizer) *models.Game {
	// The Steam API response is nested, so we need to extract the actual data
	otherDetails, ok := gameDetailsAPI[appIDKey]
	if !(ok && otherDetails.Success) {
		log.Printf("Steam API returned unsuccessful response for appID %v\n", appIDKey)
		return nil
	}

	intPrice, _ := strconv.Atoi(game.Price)
	intInitialPrice, _ := strconv.Atoi(game.InitialPrice)
	intDiscount, _ := strconv.Atoi(game.Discount)

	platforms := make([]string, 0, len(otherDetails.Data.Platforms))
	for k, v := range otherDetails.Data.Platforms {
		if v {
			platforms = append(platforms, k)
		}
	}

	gameEntry := &models.Game{
		AppId:         game.AppID,
		Name:          game.Name,
		ShortDesc:     otherDetails.Data.ShortDesc,
		Price:         intPrice,
		InitialPrice:  intInitialPrice,
		Discount:      intDiscount,
		ReleaseDate:   otherDetails.Data.ReleaseDate.Date,
		Genres:        strings.Split(gameDetailsSpy.Genres, " "),
		Tags:          gameDetailsSpy.Tags,
		Positive:      game.Positive,
		Negative:      game.Negative,
		Platforms:     platforms,
		FeatureVector: vectorizer.Vectorize(gameDetailsSpy.Genres, gameDetailsSpy.Tags, otherDetails.Data.ShortDesc),
	}

	return gameEntry
}
