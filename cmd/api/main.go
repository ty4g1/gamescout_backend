package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/ty4g1/gamescout_backend/internal/config"
	"github.com/ty4g1/gamescout_backend/internal/services"
)

func main() {
	// Load config
	cfg := config.NewConfig()

	// Connect to database
	conn, err := pgx.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	// Test database connection
	if err := conn.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v\n", err)
	}

	services.Populate()

	router := gin.Default()

	router.GET("/", ping)

	fmt.Println("Starting server...")

	router.Run(cfg.ServerAddress)
}

func ping(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "You have reached test server!"})
}
