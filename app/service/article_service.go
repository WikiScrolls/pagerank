package service

import (
	"context"
	"time"

	"github.com/WikiScrolls/pagerank/app/client"
	"github.com/WikiScrolls/pagerank/app/model"
	g "github.com/gorse-io/gorse-go"
)

type ArticleService struct {
	gorse *g.GorseClient
	wiki  *client.WikipediaClient
}

func NewArticleService(
	gorse *g.GorseClient,
	wiki *client.WikipediaClient,

) *ArticleService {
	return &ArticleService{gorse: gorse, wiki: wiki}
}

func (s *ArticleService) LikeArticle(ctx context.Context, userId string, itemId string) error {
	_, err := s.gorse.InsertFeedback(ctx, []g.Feedback{{
		FeedbackType: "like", UserId: userId, ItemId: itemId, Value: 1.0, Timestamp: time.Now(),
	}})

	return err
}

func (s *ArticleService) OpenArticle(ctx context.Context, userId string, itemId string) error {
	_, err := s.gorse.InsertFeedback(ctx, []g.Feedback{{
		FeedbackType: "open_article", UserId: userId, ItemId: itemId, Value: 1.0, Timestamp: time.Now(),
	}})

	return err
}

func (s *ArticleService) SearchArticles(ctx context.Context, search string) ([]model.Article, error) {
	searchResponse, err := s.wiki.FetchBySearch(ctx, search)
	if err != nil {
		return nil, err
	}

	var titles []string
	for _, article := range searchResponse.Query.Search {
		titles = append(titles, article.Title)
	}

	wikiResponse, err := s.wiki.FetchByTitles(ctx, titles)
	if err != nil {
		return nil, err
	}

	return wikipediaResponseToArticles(wikiResponse), nil
}
