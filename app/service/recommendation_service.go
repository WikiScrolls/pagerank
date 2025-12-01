package service

import (
	"context"
	"strconv"
	"time"

	"github.com/WikiScrolls/pagerank/app/client"
	"github.com/WikiScrolls/pagerank/app/model"

	g "github.com/gorse-io/gorse-go"
)

type RecommendationService struct {
	wiki  *client.WikipediaClient
	gorse *g.GorseClient
}

func NewRecommendationService(
	wiki *client.WikipediaClient, gorse *g.GorseClient,
) *RecommendationService {
	return &RecommendationService{wiki: wiki, gorse: gorse}
}

func (s *RecommendationService) GetRecommendations(ctx context.Context, chainLength int, userId string) ([]model.Article, error) {
	gorseIds, err := s.gorse.GetRecommend(ctx, userId, "article", chainLength, 0)

	if err == nil && len(gorseIds) > 0 {
		resp, err := s.wiki.FetchByIDs(ctx, gorseIds)
		if err == nil {
			articles := wikipediaResponseToArticles(resp)

			if len(articles) < chainLength {
				fill, _ := s.GetRandomArticles(ctx, chainLength-len(articles))
				articles = append(articles, fill...)
			}
			return articles, nil
		}
	}

	return s.GetRandomArticles(ctx, chainLength)
}

func (s *RecommendationService) GetRandomArticles(ctx context.Context, articleCount int) ([]model.Article, error) {
	wikiResponse, err := s.wiki.GetRandomArticles(ctx, articleCount)
	if err != nil {
		return nil, err
	}

	articles := wikipediaResponseToArticles(wikiResponse)

	for _, article := range articles {
		s.gorse.InsertItem(ctx, g.Item{
			ItemId:     article.Id,
			IsHidden:   false,
			Labels:     []string{"wikipedia", "article"},
			Categories: []string{"article"},
			Comment:    article.Title,
			Timestamp:  time.Now(),
		})
	}

	return articles, err
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
