package main

import (
	"fmt"
	"log"

	"github.com/WikiScrolls/pagerank/app"
	"github.com/WikiScrolls/pagerank/app/config"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	app, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	app.Routes(router)

	router.Run(":" + cfg.AppPort)
	fmt.Println("Pagerank Running on Port :" + cfg.AppPort)
}
