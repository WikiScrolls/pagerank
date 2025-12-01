package service

import (
	"context"
	"time"

	g "github.com/gorse-io/gorse-go"
)

type ArticleService struct {
	gorse *g.GorseClient
}

func NewArticleService(
	gorse *g.GorseClient,
) *RecommendationService {
	return &RecommendationService{gorse: gorse}
}

func (s *ArticleService) LikeArticle(ctx context.Context, userId string, itemId string) error {
	_, err := s.gorse.InsertFeedback(ctx, []g.Feedback{{
		FeedbackType: "like", UserId: userId, ItemId: itemId, Value: 1.0, Timestamp: time.Now(),
	}})

	return err
}

func (s *ArticleService) ViewArticle(ctx context.Context, userId string, itemId string) error {
	_, err := s.gorse.InsertFeedback(ctx, []g.Feedback{{
		FeedbackType: "open-article", UserId: userId, ItemId: itemId, Value: 1.0, Timestamp: time.Now(),
	}})

	return err
}
