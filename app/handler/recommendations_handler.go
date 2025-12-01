package handler

import (
	"github.com/WikiScrolls/pagerank/app/service"
	"github.com/gin-gonic/gin"
)

type RecommendationHandler struct {
	serv *service.RecommendationService
}

func NewRecommendationHandler(serv *service.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{serv: serv}
}

func (h *RecommendationHandler) GetRecommendations(c *gin.Context) {
	userId := c.Param("userId")

	articles, err := h.serv.GetRecommendations(c.Request.Context(), 10, userId)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
	} else {
		c.JSON(200, gin.H{
			"data": articles,
		})
	}
}

func (h *RecommendationHandler) GetRandomArticles(c *gin.Context) {
	articles, err := h.serv.GetRandomArticles(c.Request.Context(), 20)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
	} else {
		c.JSON(200, gin.H{
			"data": articles,
		})
	}
}
