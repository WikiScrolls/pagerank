package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Neo4jUri      string
	Neo4jUser     string
	Neo4jPassword string
	GorseURL      string
	GorseKey      string
	AppPort       string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	return &Config{
		Neo4jUri:      os.Getenv("NEO4J_URI"),
		Neo4jUser:     os.Getenv("NEO4J_USER"),
		Neo4jPassword: os.Getenv("NEO4J_PASSWORD"),
		GorseURL:      os.Getenv("GORSE_URL"),
		GorseKey:      os.Getenv("GORSE_KEY"),
		AppPort:       os.Getenv("AppPort"),
	}, nil
}
