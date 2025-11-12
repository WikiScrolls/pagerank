package main

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var demoRecommendations = [][]string{
	{"'Uwe_Ludwig'", "'Dentistry_in_Saint_Lucia'", "'Japan'", "'Kyrgyzstan_Sweden_relations'"},
	{"'Laurier_Regnier'", "'The_New_Barbarians_(band'", "'Sahare,_Surkhet'", "'The_Rolling_Stones:_Voodoo_Lounge_Live'"},
	{"'Homoserine_Lactone'", "'Beijing_Municipal_No._2_Prison'", "'November_criminal'", "'List_of_prisons_in_Anhui'"},
}

func formatToWikiLink(title string) string {
	clean := strings.Trim(title, "'")
	clean = strings.ReplaceAll(clean, " ", "_")
	link := "https://en.wikipedia.org/wiki/" + clean
	if strings.Contains(link, "(") {
		link += ")"
	}
	return link
}

func getRecommendations(ctx context.Context) ([]string, error) {

	titlesArr := demoRecommendations[rand.Intn(len(demoRecommendations))]

	var links []string
	for _, title := range titlesArr {
		links = append(links, formatToWikiLink(title))
	}

	return links, nil
}

func main() {
	rand.Seed(time.Now().Unix())

	router := gin.Default()

	router.GET("/recommendations", func(c *gin.Context) {
		ctx := context.Background()
		links, err := getRecommendations(ctx)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"recommendations": links})
	})

	router.Run(":8080")
}
