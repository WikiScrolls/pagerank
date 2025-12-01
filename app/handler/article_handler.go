package handler

import (
	"github.com/WikiScrolls/pagerank/app/service"
	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	serv *service.ArticleService
}

func NewArticleHandler(serv *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{serv: serv}
}

func (h *ArticleHandler) LikeArticle(c *gin.Context) {
	itemId := c.Param("id")
	userId := c.Query("userId")

	err := h.serv.LikeArticle(c.Request.Context(), userId, itemId)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
	} else {
		c.JSON(200, gin.H{
			"message": "Liked Successfully",
		})
	}
}

func (h *ArticleHandler) OpenArticle(c *gin.Context) {
	itemId := c.Param("id")
	userId := c.Query("userId")

	err := h.serv.OpenArticle(c.Request.Context(), userId, itemId)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
	} else {
		c.JSON(200, gin.H{
			"message": "Opened Successfully",
		})
	}
}
