package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ty4g1/gamescout_backend/internal/config"
	"github.com/ty4g1/gamescout_backend/internal/models"
	"github.com/ty4g1/gamescout_backend/internal/repository"
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

func (p *Populator) Populate(gr *repository.GameRepository, gmr *repository.GameMediaRepository) {
	// Each page returns 1000 game entries
	totalGames := p.Pages * 1000
	counter := 0
	for i := range p.Pages {
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

		gameEntries := make([]*models.Game, 0, 1000)
		gameMediaEntries := make([]*models.GameMedia, 0, 1000)

		for appIDKey, game := range games {
			if counter%1000 == 0 {
				fmt.Printf("Processed games: %d/%d\n", counter, totalGames)
			}
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

			// Get Steam API details
			respAPI, err := http.Get(fmt.Sprintf(p.SteamAPIDetailsFormat, appID))
			if err != nil {
				log.Printf("Error making Steam API details request for appID %d: %v\n", appID, err)
				continue
			}

			var gameDetailsApi models.SteamAPIDetailsWrapper
			if err := json.NewDecoder(respAPI.Body).Decode(&gameDetailsApi); err != nil {
				log.Printf("Error parsing Steam API details for appID %d: %v\n", appID, err)
				respAPI.Body.Close()
				continue
			}
			respAPI.Body.Close()

			// The Steam API response is nested, so we need to extract the actual data
			otherDetails, ok := gameDetailsApi[appIDKey]
			if !(ok && otherDetails.Success) {
				log.Printf("Steam API returned unsuccessful response for appID %v\n", appIDKey)
				continue
			}

			gameEntry, err := createGameEntry(game, gameDetailsSpy, otherDetails.Data, appIDKey, p.Vectorizer)
			if err != nil {
				log.Printf("Error creating game entry for appID %v: %v\n", appIDKey, err)
				continue
			}

			gameMediaEntry := createGameMediaEntry(otherDetails.Data, appID)

			gameEntries = append(gameEntries, gameEntry)
			gameMediaEntries = append(gameMediaEntries, gameMediaEntry)
			counter += 1
		}

		err = gr.BatchInsert(context.Background(), gameEntries)
		if err != nil {
			log.Printf("Error inserting entries into Games table %v\n", err)
			continue
		}
		err = gmr.BatchInsert(context.Background(), gameMediaEntries)
		if err != nil {
			log.Printf("Error inserting entries into Games_media table %v\n", err)
			continue
		}
	}
}

func createGameEntry(game models.SteamspyResponse, gameDetailsSpy models.SteamspyDetails, gameDetailsApi models.SteamAPIDetails, appIDKey string, vectorizer *Vectorizer) (*models.Game, error) {
	intPrice, _ := strconv.Atoi(game.Price)
	intInitialPrice, _ := strconv.Atoi(game.InitialPrice)
	intDiscount, _ := strconv.Atoi(game.Discount)

	platforms := make([]string, 0, len(gameDetailsApi.Platforms))
	for k, v := range gameDetailsApi.Platforms {
		if v {
			platforms = append(platforms, k)
		}
	}

	release_date, err := time.Parse("Jan 2, 2006", gameDetailsApi.ReleaseDate.Date)
	if err != nil {
		return nil, err
	}

	gameEntry := &models.Game{
		AppId:         game.AppID,
		Name:          game.Name,
		ShortDesc:     gameDetailsApi.ShortDesc,
		Price:         intPrice,
		InitialPrice:  intInitialPrice,
		Discount:      intDiscount,
		ReleaseDate:   release_date,
		Genres:        strings.Split(gameDetailsSpy.Genres, " "),
		Tags:          gameDetailsSpy.Tags,
		Positive:      game.Positive,
		Negative:      game.Negative,
		Platforms:     platforms,
		FeatureVector: vectorizer.Vectorize(gameDetailsSpy.Genres, gameDetailsSpy.Tags, gameDetailsApi.ShortDesc),
	}

	return gameEntry, nil
}

func createGameMediaEntry(gameDetailsAPI models.SteamAPIDetails, appId int) *models.GameMedia {
	gameMediaEntry := &models.GameMedia{
		AppID:         appId,
		ThumbnailURL:  gameDetailsAPI.Thumbnail,
		BackgroundURL: gameDetailsAPI.Background,
		Screenshots:   gameDetailsAPI.Screenshots,
		Movies:        gameDetailsAPI.Movies,
	}

	return gameMediaEntry
}
