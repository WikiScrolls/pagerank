package handler

import (
	"github.com/WikiScrolls/pagerank/app/service"
)

type ArticleHandler struct {
	serv *service.RecommendationService
}

func NewArticleHandler(serv *service.ArticleService) *RecommendationHandler {
	return &RecommendationHandler{serv: serv}
}
