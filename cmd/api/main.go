package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ty4g1/gamescout_backend/internal/config"
	"github.com/ty4g1/gamescout_backend/internal/repository"
	"github.com/ty4g1/gamescout_backend/internal/services"
)

func main() {
	// Load config
	cfg := config.NewConfig()

	// Load populator
	pop := services.NewPopulator(cfg)

	// Connect to database
	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbpool.Close()

	// Create repositories
	gr := repository.NewGamesRepository(dbpool)
	gmr := repository.NewGameMediaRepository(dbpool)

	// Populate database
	fmt.Println("Populating database...")
	pop.Populate(gr, gmr)

	router := gin.Default()

	router.GET("/", ping)

	fmt.Println("Starting server...")

	router.Run(cfg.ServerAddress)
}

func ping(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "You have reached test server!"})
}
