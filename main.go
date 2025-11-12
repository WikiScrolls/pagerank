package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var driver neo4j.DriverWithContext

func connectDB() {
	ctx := context.Background()
	_ = godotenv.Load()

	dbUri := os.Getenv("NEO4J_URI")
	dbUser := os.Getenv("NEO4J_USER")
	dbPassword := os.Getenv("NEO4J_PASSWORD")

	var err error
	driver, err = neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""),
	)
	if err != nil {
		panic(err)
	}

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("✅ Connection established.")
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
	result, err := neo4j.ExecuteQuery(
		ctx,
		driver,
		`
		MATCH (n)
		WITH n, rand() AS r
		ORDER BY r
		LIMIT 1
		MATCH p = (n)-[*1..3]-(m)
		WITH p, rand() AS r2
		ORDER BY r2
		LIMIT 1
		RETURN [x IN nodes(p) | x.title] AS pathTitles
		`,
		nil,
		neo4j.EagerResultTransformer,
	)
	if err != nil {
		return nil, err
	}

	if len(result.Records) == 0 {
		return []string{}, nil
	}

	titles, _ := result.Records[0].Get("pathTitles")
	titlesArr, ok := titles.([]interface{})
	if !ok {
		return []string{}, nil
	}

	var links []string
	for _, t := range titlesArr {
		title, ok := t.(string)
		if !ok || title == "" {
			continue
		}
		links = append(links, formatToWikiLink(title))
	}

	return links, nil
}

func main() {
	connectDB()
	defer driver.Close(context.Background())

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
