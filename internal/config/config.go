package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	ServerAddress         string
	ApiKey                string
	DatabaseURL           string
	MicroserviceURL       string
	SteamSpyPages         int
	SteamSpyURLFormat     string
	SteamSpyDetailsFormat string
	SteamAPIDetailsFormat string
}

func NewConfig() *Config {
	serverAddr := os.Getenv("SERVER_ADDRESS")
	apiKey := os.Getenv("API_KEY")
	databaseUrl := os.Getenv("DATABASE_URL")
	microserviceUrl := os.Getenv("MICROSERVICE_URL")
	steamSpyPages, err := strconv.Atoi(os.Getenv("STEAM_SPY_PAGES"))
	if err != nil {
		log.Printf("Error getting page limit from env, defaulting to 5: %v\n", err)
		steamSpyPages = 5
	}
	steamSpyUrlFormat := os.Getenv("STEAM_SPY_URL_FORMAT")
	steamSpyDetailsFormat := os.Getenv("STEAM_SPY_DETAILS_FORMAT")
	steamApiDetailsFormat := os.Getenv("STEAM_API_DETAILS_FORMAT")
	return &Config{
		ServerAddress:         serverAddr,
		ApiKey:                apiKey,
		DatabaseURL:           databaseUrl,
		MicroserviceURL:       microserviceUrl,
		SteamSpyPages:         steamSpyPages,
		SteamSpyURLFormat:     steamSpyUrlFormat,
		SteamSpyDetailsFormat: steamSpyDetailsFormat,
		SteamAPIDetailsFormat: steamApiDetailsFormat,
	}
}
