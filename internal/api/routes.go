package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ty4g1/gamescout_backend/internal/api/handlers"
	"github.com/ty4g1/gamescout_backend/internal/repository"
)

func SetupRouter(gr *repository.GameRepository, gmr *repository.GameMediaRepository) *gin.Engine {
	router := gin.Default()

	gameHandler := handlers.NewGameHandler(gr, gmr)

	router.GET("/", ping)
	router.GET("/games/random", gameHandler.GetRandomGames)

	return router
}

func ping(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "You have reached test server!"})
}
