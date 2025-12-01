package handler

import (
	"github.com/WikiScrolls/pagerank/app/model"
	"github.com/WikiScrolls/pagerank/app/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	serv *service.UserService
}

func NewUserHandler(serv *service.UserService) *UserHandler {
	return &UserHandler{serv: serv}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(401, gin.H{
			"error": "Invalid user data from body.",
		})
	}

	err := h.serv.RegisterUser(c.Request.Context(), user.Id, user.Interests)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
	} else {
		c.JSON(200, gin.H{
			"message": "Successfully registered user to recommender.",
		})
	}
}
