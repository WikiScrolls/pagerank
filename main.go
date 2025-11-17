package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
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
	fmt.Println("Connection established.")
}

func cleanTitle(title string) string {
	clean := strings.Trim(title, "'")
	clean = strings.ReplaceAll(clean, " ", "_")
	if strings.Contains(clean, "(") {
		clean += ")"
	}
	return clean
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

	var retTitles []string
	for _, t := range titlesArr {
		title, ok := t.(string)
		if !ok || title == "" {
			continue
		}
		retTitles = append(retTitles, cleanTitle(title))
	}

	return retTitles, nil
}

type Article struct {
	Id           string `json:"id"`
	Title        string `json:"title"`
	WikipediaUrl string `json:"wikipediaUrl"`
	Content      string `json:"content"`
	Thumbnail    string `json:"thumbnail"`
}

type WikipediaResponse struct {
	Continue struct {
		Excontinue  int    `json:"excontinue"`
		Grncontinue string `json:"grncontinue"`
		Continue    string `json:"continue"`
	} `json:"continue"`

	Query struct {
		Pages map[string]Page `json:"pages"`
	} `json:"query"`
}

type Page struct {
	PageID           int    `json:"pageid"`
	Ns               int    `json:"ns"`
	Title            string `json:"title"`
	Extract          string `json:"extract,omitempty"`
	ContentModel     string `json:"contentmodel"`
	PageLanguage     string `json:"pagelanguage"`
	PageLanguageHTML string `json:"pagelanguagehtmlcode"`
	PageLanguageDir  string `json:"pagelanguagedir"`
	Touched          string `json:"touched"`
	LastRevID        int    `json:"lastrevid"`
	Length           int    `json:"length"`
	FullURL          string `json:"fullurl"`
	EditURL          string `json:"editurl"`
	CanonicalURL     string `json:"canonicalurl"`
	Thumbnail        struct {
		Source string `json:"source"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"thumbnail"`
}

func unMarshalArticles(jsonRaw []byte) ([]Article, error) {
	wikiResponse := WikipediaResponse{}
	json.Unmarshal(jsonRaw, &wikiResponse)

	var articles []Article
	for _, page := range wikiResponse.Query.Pages {
		articles = append(articles, Article{
			strconv.Itoa(page.PageID), page.Title, page.FullURL, page.Extract, page.Thumbnail.Source,
		})
	}

	return articles, nil
}

func hydrateWikipedia(titles []string) ([]Article, error) {
	var titleQuery string = titles[0]
	for i := 1; i < len(titles); i++ {
		titleQuery += "|" + titles[i]
	}

	params := url.Values{
		"action":      {"query"},
		"format":      {"json"},
		"prop":        {"extracts|info|pageimages"},
		"inprop":      {"url"},
		"exintro":     {"1"},
		"exlimit":     {"max"},
		"exsentences": {"10"},
		"explaintext": {"1"},
		"origin":      {"*"},
		"variant":     {"en"},
		"piprop":      {"thumbnail"},
		"pithumbsize": {"800"},
		"titles":      {titleQuery},
	}

	reqURL := "https://en.wikipedia.org/w/api.php?" + params.Encode()
	fmt.Printf("URL: %s\n", reqURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "WikiScrolls/1.0 (; nadzhiff@gmail.com) Go-http-client/1.1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic("Failed to fetch http")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	os.WriteFile("res.json", body, 0644)

	// return nil, nil
	return unMarshalArticles(body)
}

func main() {
	connectDB()
	defer driver.Close(context.Background())

	router := gin.Default()

	router.GET("/recommendations", func(c *gin.Context) {
		ctx := context.Background()
		titles, err := getRecommendations(ctx)
		articles, err := hydrateWikipedia(titles)
		if err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"message": "OK",
			"data": gin.H{
				"articles": articles,
			},
		})
	})

	router.Run(":8080")
}
