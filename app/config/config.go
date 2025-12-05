package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GorseURL string
	GorseKey string
	AppPort  string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	return &Config{
		GorseURL: os.Getenv("GORSE_URL"),
		GorseKey: os.Getenv("GORSE_KEY"),
		AppPort:  os.Getenv("APP_PORT"),
	}, nil
}
