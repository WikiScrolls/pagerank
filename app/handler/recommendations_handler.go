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
	articles, err := h.serv.GetRecommendations(c.Request.Context(), 3)

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
