package service

import (
	"context"
	"strconv"

	"github.com/WikiScrolls/pagerank/app/client"
	"github.com/WikiScrolls/pagerank/app/model"
	"github.com/WikiScrolls/pagerank/app/repository"
)

type RecommendationService struct {
	repo repository.RecommendationRepository
	wiki *client.WikipediaClient
}

func NewRecommendationService(
	repo repository.RecommendationRepository, wiki *client.WikipediaClient,
) *RecommendationService {
	return &RecommendationService{repo: repo, wiki: wiki}
}

func (s *RecommendationService) GetRecommendations(ctx context.Context, chainLength int) ([]model.Article, error) {
	titles, err := s.repo.GetRecommendationTitles(ctx, chainLength)
	if err != nil {
		return nil, err
	}

	wikipediaResponse, err := s.wiki.FetchByTitles(ctx, titles)
	if err != nil {
		return nil, err
	}

	articles := wikipediaResponseToArticles(wikipediaResponse)

	if len(articles) == chainLength {
		return articles, nil
	}

	fillerWikipedia, err := s.wiki.GetRandomArticles(ctx, chainLength-len(articles))
	if err != nil {
		return articles, err
	}

	articles = append(articles, wikipediaResponseToArticles(fillerWikipedia)...)

	return articles, nil
}

func (s *RecommendationService) GetRandomArticles(ctx context.Context, articleCount int) ([]model.Article, error) {
	wikiResponse, err := s.wiki.GetRandomArticles(ctx, articleCount)
	if err != nil {
		return nil, err
	}
	return wikipediaResponseToArticles(wikiResponse), nil
}

func wikipediaResponseToArticles(wikipediaResponse *model.WikipediaResponse) []model.Article {
	var articles []model.Article
	for _, page := range wikipediaResponse.Query.Pages {
		if page.Extract == "" {
			continue
		}
		articles = append(articles, model.Article{
			Id:           strconv.Itoa(page.PageID),
			Title:        page.Title,
			WikipediaUrl: page.FullURL,
			Content:      page.Extract,
			Thumbnail:    page.Thumbnail.Source,
		})
	}
	return articles
}
