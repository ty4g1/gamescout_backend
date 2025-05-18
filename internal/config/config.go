package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerAddress   string
	ApiKey          string
	DatabaseURL     string
	MicroserviceURL string
}

func NewConfig() *Config {
	serverAddr := os.Getenv("SERVER_ADDRESS")
	apiKey := os.Getenv("API_KEY")
	databaseUrl := os.Getenv("DATABASE_URL")
	microserviceUrl := os.Getenv("MICROSERVICE_URL")
	fmt.Println(serverAddr, databaseUrl)
	return &Config{
		ServerAddress:   serverAddr,
		ApiKey:          apiKey,
		DatabaseURL:     databaseUrl,
		MicroserviceURL: microserviceUrl,
	}
}
