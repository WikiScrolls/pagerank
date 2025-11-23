package app

import (
	"context"

	"github.com/WikiScrolls/pagerank/app/client"
	"github.com/WikiScrolls/pagerank/app/config"
	"github.com/WikiScrolls/pagerank/app/database"
	"github.com/WikiScrolls/pagerank/app/repository"
	"github.com/WikiScrolls/pagerank/app/service"
)

type App struct {
	RecommendationService service.RecommendationService
}

func New(cfg *config.Config) (*App, error) {
	ctx := context.Background()

	neo4jDatabase, err := database.NewNeo4jClient(
		ctx,
		cfg.Neo4jUri,
		cfg.Neo4jUser,
		cfg.Neo4jPassword,
	)
	if err != nil {
		return nil, err
	}

	recommendationRepository := repository.NewNeo4jRecommendationRepository(neo4jDatabase)
	wikiClient := client.NewWikipediaClient()

	return &App{
		RecommendationService: *service.NewRecommendationService(
			recommendationRepository,
			wikiClient,
		),
	}, nil
}
