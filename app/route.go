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
			recommendation.GET("/:userId", h.GetRecommendations) // Get recommendations for user ID
			recommendation.GET("/random", h.GetRandomArticles)   // Random articles
		}

		article := api.Group("/articles/:id")
		{
			h := handler.NewArticleHandler(&a.ArticleService)
			article.POST("/like", h.LikeArticle) // Like article id for user ?userId
			article.POST("/open", h.OpenArticle) // Open article id for user ?userId
		}

		user := api.Group("/user")
		{
			h := handler.NewUserHandler(&a.UserService)
			user.POST("/", h.RegisterUser) // Register user to recommender with format {"id": userId, "interests": ["interest1", "interest2", ...]}
		}
	}
}
