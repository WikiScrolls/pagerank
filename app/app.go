package app

import (
	"github.com/WikiScrolls/pagerank/app/client"
	"github.com/WikiScrolls/pagerank/app/config"
	"github.com/WikiScrolls/pagerank/app/service"

	gorse "github.com/gorse-io/gorse-go"
)

type App struct {
	RecommendationService service.RecommendationService
}

func New(cfg *config.Config) (*App, error) {

	wikiClient := client.NewWikipediaClient()
	gorseClient := gorse.NewGorseClient("http://localhost:8088", "")

	return &App{
		RecommendationService: *service.NewRecommendationService(
			wikiClient,
			gorseClient,
		),
	}, nil
}
