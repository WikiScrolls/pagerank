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

	"github.com/gin-gonic/gin"
)

type Article struct {
	Id           string `json:"id"`
	Title        string `json:"title"`
	WikipediaUrl string `json:"wikipediaUrl"`
	Content      string `json:"content"`
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
}

func unMarshalArticles(jsonRaw []byte) ([]Article, error) {
	wikiResponse := WikipediaResponse{}
	json.Unmarshal(jsonRaw, &wikiResponse)

	var articles []Article
	for _, page := range wikiResponse.Query.Pages {
		articles = append(articles, Article{
			strconv.Itoa(page.PageID), page.Title, page.FullURL, page.Extract,
		})
	}

	return articles, nil
}

func getArticles(ctx context.Context) ([]Article, error) {
	params := url.Values{
		"action":       {"query"},
		"format":       {"json"},
		"generator":    {"random"},
		"grnnamespace": {"0"},
		"grnlimit":     {"20"},
		"prop":         {"extracts|info"},
		"inprop":       {"url"},
		"exintro":      {"1"},
		"exlimit":      {"max"},
		"exsentences":  {"10"},
		"explaintext":  {"1"},
		"origin":       {"*"},
		"variant":      {"en"},
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
	router := gin.Default()

	router.GET("/recommendations", func(c *gin.Context) {
		ctx := context.Background()
		articles, err := getArticles(ctx)
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
