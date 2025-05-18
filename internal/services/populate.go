package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ty4g1/gamescout_backend/internal/models"
)

func Populate() {
	// Number of pages (with 1000 games each) to retrieve
	const PAGES = 10
	const STEAM_SPY_URL_FORMAT = "https://steamspy.com/api.php?request=all&page=%d"
	const STEAM_SPY_DETAILS_URL_FORMAT = "https://steamspy.com/api.php?request=appdetails&appid=%d"
	const STEAM_API_DETAILS_URL_FORMAT = "https://store.steampowered.com/api/appdetails?cc=us&appids=%d"

	for i := range PAGES {
		resp, err := http.Get(fmt.Sprintf(STEAM_SPY_URL_FORMAT, i))
		if err != nil {
			log.Fatalf("Error making request: %v\n", err)
		}

		var games map[string]models.SteamspyResponse
		if err := json.NewDecoder(resp.Body).Decode(&games); err != nil {
			log.Fatalf("Error reading and parsing JSON: %v\n", err)
		}
		resp.Body.Close()

		for appIDKey, game := range games {
			appID := game.AppID // Use the AppID from the struct
			fmt.Printf("Got the following game: %s (AppID: %d)\n", game.Name, appID)

			// Respect rate limits (1 request per second)
			time.Sleep(1 * time.Second)

			// Get Steam Spy details
			respSpy, err := http.Get(fmt.Sprintf(STEAM_SPY_DETAILS_URL_FORMAT, appID))
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
			respAPI, err := http.Get(fmt.Sprintf(STEAM_API_DETAILS_URL_FORMAT, appID))
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

			// The Steam API response is nested, so we need to extract the actual data
			if details, ok := gameDetailsAPI[appIDKey]; ok && details.Success {
				fmt.Printf("Steam API details: %+v\n", details.Data)
			} else {
				log.Printf("Steam API returned unsuccessful response for appID %d\n", appID)
			}

			break // Only process one game for testing
		}

		// Temporary, for testing with just 1 page
		break
	}
}
