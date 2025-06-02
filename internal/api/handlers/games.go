package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ty4g1/gamescout_backend/internal/models"
	"github.com/ty4g1/gamescout_backend/internal/repository"
	"github.com/ty4g1/gamescout_backend/internal/utils"
)

type GameHandler struct {
	gr  *repository.GameRepository
	gmr *repository.GameMediaRepository
	ur  *repository.UserRepository
}

func NewGameHandler(gr *repository.GameRepository, gmr *repository.GameMediaRepository, ur *repository.UserRepository) *GameHandler {
	return &GameHandler{
		gr:  gr,
		gmr: gmr,
		ur:  ur,
	}
}

func (gh *GameHandler) GetRandomGames(c *gin.Context) {
	// Parse limit
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	priceRange := &models.PriceRange{Min: 0, Max: 10000} // Always create with defaults

	// Override min if provided
	if minStr := c.Query("min_price"); minStr != "" {
		if min, err := strconv.Atoi(minStr); err == nil && min >= 0 {
			priceRange.Min = min
		}
	}

	// Override max if provided
	if maxStr := c.Query("max_price"); maxStr != "" {
		if max, err := strconv.Atoi(maxStr); err == nil && max >= 0 {
			priceRange.Max = max
		}
	}

	// Parse release date
	var releaseDate *models.ReleaseDate
	if dateStr := c.Query("release_date"); dateStr != "" {
		if date, err := time.Parse("2006-01-02", dateStr); err == nil {
			isBefore := c.Query("before") == "true"
			releaseDate = &models.ReleaseDate{Date: date, IsBefore: isBefore}
		}
	}

	// Parse arrays (handle empty strings as nil)
	var tags, genres, platforms []string
	if tagsStr := c.Query("tags"); tagsStr != "" {
		tags = strings.Split(tagsStr, ",")
		// Trim whitespace and URL decode
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
			if decoded, err := url.QueryUnescape(tags[i]); err == nil {
				tags[i] = decoded
			}
		}
	}

	if genresStr := c.Query("genres"); genresStr != "" {
		genres = strings.Split(genresStr, ",")
		for i := range genres {
			genres[i] = strings.TrimSpace(genres[i])
			if decoded, err := url.QueryUnescape(genres[i]); err == nil {
				genres[i] = decoded
			}
		}
	}

	if platformsStr := c.Query("platforms"); platformsStr != "" {
		platforms = strings.Split(platformsStr, ",")
		for i := range platforms {
			platforms[i] = strings.TrimSpace(platforms[i])
			if decoded, err := url.QueryUnescape(platforms[i]); err == nil {
				platforms[i] = decoded
			}
		}
	}

	games, err := gh.gr.GetRandom(c.Request.Context(), limit, priceRange, releaseDate, tags, genres, platforms)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get games"})
		return
	}

	type GameWithMedia struct {
		models.Game
		Media *models.GameMedia `json:"media,omitempty"`
	}

	var response []GameWithMedia
	for _, game := range games {
		gameWithMedia := GameWithMedia{Game: game}

		if media, err := gh.gmr.GetByAppID(c.Request.Context(), game.AppId); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get game media"})
			return
		} else {
			gameWithMedia.Media = media
		}

		response = append(response, gameWithMedia)
	}

	c.JSON(http.StatusOK, gin.H{"games": response, "count": len(games)})
}

func (gh *GameHandler) GetRecommendations(c *gin.Context) {
	// Parse id
	var id string
	if idStr := c.Query("id"); idStr != "" {
		id = idStr
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Parse limit
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	priceRange := &models.PriceRange{Min: 0, Max: 10000} // Always create with defaults

	// Override min if provided
	if minStr := c.Query("min_price"); minStr != "" {
		if min, err := strconv.Atoi(minStr); err == nil && min >= 0 {
			priceRange.Min = min
		}
	}

	// Override max if provided
	if maxStr := c.Query("max_price"); maxStr != "" {
		if max, err := strconv.Atoi(maxStr); err == nil && max >= 0 {
			priceRange.Max = max
		}
	}

	// Parse release date
	var releaseDate *models.ReleaseDate
	if dateStr := c.Query("release_date"); dateStr != "" {
		if date, err := time.Parse("2006-01-02", dateStr); err == nil {
			isBefore := c.Query("before") == "true"
			releaseDate = &models.ReleaseDate{Date: date, IsBefore: isBefore}
		}
	}

	// Parse arrays (handle empty strings as nil)
	var tags, genres, platforms []string
	if tagsStr := c.Query("tags"); tagsStr != "" {
		tags = strings.Split(tagsStr, ",")
		// Trim whitespace
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
	}

	if genresStr := c.Query("genres"); genresStr != "" {
		genres = strings.Split(genresStr, ",")
		for i := range genres {
			genres[i] = strings.TrimSpace(genres[i])
		}
	}

	if platformsStr := c.Query("platforms"); platformsStr != "" {
		platforms = strings.Split(platformsStr, ",")
		for i := range platforms {
			platforms[i] = strings.TrimSpace(platforms[i])
		}
	}

	games, err := gh.gr.GetAll(c.Request.Context(), priceRange, releaseDate, tags, genres, platforms)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get games"})
		return
	}

	preferenceVector, err := gh.ur.GetUserPreference(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get user preferences: %v", err)})
		return
	}

	swipeHistory, err := gh.ur.GetUserSwipes(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user swipe history"})
		return
	}

	swipeSet := make(map[int]bool, 100)
	for _, appId := range swipeHistory[len(swipeHistory)-100:] {
		swipeSet[appId] = true
	}

	filteredGames := utils.Filter(games, func(g models.Game) bool {
		return !swipeSet[g.AppId]
	})

	sort.Slice(filteredGames, func(i, j int) bool {
		// Comparison logic based on the key
		s1, err := utils.ComputeSimilarity(preferenceVector, filteredGames[i].FeatureVector)
		if err != nil {
			return false
		}
		s2, err := utils.ComputeSimilarity(preferenceVector, filteredGames[j].FeatureVector)
		if err != nil {
			return true
		}
		return s1 > s2
	})

	if len(filteredGames) > limit {
		games = filteredGames[:limit]
	} else {
		games = filteredGames
	}

	type GameWithMedia struct {
		models.Game
		Media *models.GameMedia `json:"media,omitempty"`
	}

	var response []GameWithMedia
	for _, game := range games {
		gameWithMedia := GameWithMedia{Game: game}

		if media, err := gh.gmr.GetByAppID(c.Request.Context(), game.AppId); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get game media"})
			return
		} else {
			gameWithMedia.Media = media
		}

		response = append(response, gameWithMedia)
	}

	c.JSON(http.StatusOK, gin.H{"games": response, "count": len(games)})
}

func (gh *GameHandler) GetTags(c *gin.Context) {
	tags, err := gh.gr.GetAllTags(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tags"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tags": tags, "count": len(tags)})
}

func (gh *GameHandler) GetGenres(c *gin.Context) {
	genres, err := gh.gr.GetAllGenres(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get genres"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"genres": genres, "count": len(genres)})
}
