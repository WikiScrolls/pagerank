package app

import (
	"github.com/WikiScrolls/pagerank/app/handler"
	"github.com/gin-gonic/gin"
)

func (a *App) Routes(router *gin.Engine) {
	api := router.Group("/api")
	{
		recommendation := api.Group("/recommendation")
		{
			h := handler.NewRecommendationHandler(&a.RecommendationService)
			recommendation.GET("/", h.GetRecommendations)
			recommendation.GET("/random", h.GetRandomArticles)
		}
	}
}
