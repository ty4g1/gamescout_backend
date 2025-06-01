package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	routes "github.com/ty4g1/gamescout_backend/internal/api"
	"github.com/ty4g1/gamescout_backend/internal/config"
	"github.com/ty4g1/gamescout_backend/internal/repository"
	"github.com/ty4g1/gamescout_backend/internal/services"
)

func main() {
	// Load config
	cfg := config.NewConfig()

	// Connect to database
	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbpool.Close()

	// Create repositories
	gr := repository.NewGamesRepository(dbpool)
	gmr := repository.NewGameMediaRepository(dbpool)
	ur := repository.NewUserRepository(dbpool)

	go func() {
		pop := services.NewPopulator(cfg)
		for {
			fmt.Println("Starting database population in background...")
			start := time.Now()
			pop.Populate(gr, gmr)
			fmt.Printf("Database population completed successfully in %v!\n", time.Since(start))
			fmt.Println("Next population will run in 24 hours...")
			time.Sleep(24 * time.Hour)
		}
	}()

	router := routes.SetupRouter(gr, gmr, ur)

	fmt.Println("Starting server...")

	router.Run(cfg.ServerAddress)
}
