package routes

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ty4g1/gamescout_backend/internal/api/handlers"
	"github.com/ty4g1/gamescout_backend/internal/repository"
)

func SetupRouter(gr *repository.GameRepository, gmr *repository.GameMediaRepository, ur *repository.UserRepository) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://165.22.59.11:4173", // preview/production
			"http://localhost:5173",    // dev server (if you ever run npm run dev)
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	gameHandler := handlers.NewGameHandler(gr, gmr, ur)
	userHandler := handlers.NewUserHandler(ur, gr)

	router.GET("/health", healthCheck)
	router.GET("/games/random", gameHandler.GetRandomGames)
	router.GET("/games/recommend", gameHandler.GetRecommendations)
	router.GET("/games/tags", gameHandler.GetTags)
	router.GET("/games/genres", gameHandler.GetGenres)

	router.POST("/users/add", userHandler.AddUser)
	router.PATCH("/users/preferences", userHandler.UpdatePreference)

	return router
}

func healthCheck(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Service is healthy!"})
}
