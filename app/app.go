package app

import (
	"github.com/WikiScrolls/pagerank/app/client"
	"github.com/WikiScrolls/pagerank/app/config"
	"github.com/WikiScrolls/pagerank/app/service"

	gorse "github.com/gorse-io/gorse-go"
)

type App struct {
	RecommendationService service.RecommendationService
	ArticleService        service.ArticleService
	UserService           service.UserService
}

func New(cfg *config.Config) (*App, error) {

	wikiClient := client.NewWikipediaClient()
	gorseClient := gorse.NewGorseClient(cfg.GorseURL, cfg.GorseKey)
	geminiClient, err := client.NewGeminiClient(cfg.APIKey)
	if err != nil {
		return nil, err
	}

	return &App{
		RecommendationService: *service.NewRecommendationService(
			wikiClient,
			gorseClient,
		),
		ArticleService: *service.NewArticleService(
			gorseClient,
			wikiClient,
			geminiClient,
		),
		UserService: *service.NewUserService(
			gorseClient,
		),
	}, nil
}
