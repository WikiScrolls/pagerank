package client

import (
	"context"

	"google.golang.org/genai"
)

type GeminiClient struct {
	driver *genai.Client
}

func NewGeminiClient(apiKey string) (*GeminiClient, error) {
	driver, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, err
	}

	return &GeminiClient{driver: driver}, nil
}

func (g *GeminiClient) SummarizeWiki(article string) (string, error) {
	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText("For the following wikipedia summary I fetched from the API. You should create a summary (3 paragraph max) of the article to give a general overview of the topic for the user. Make it casual and fun to be highly engaging for users while still not missing out the important stuffs.", genai.RoleUser),
	}

	result, err := g.driver.Models.GenerateContent(
		context.Background(),
		"gemini-2.5-flash",
		genai.Text(article),
		config,
	)
	if err != nil {
		return "", err
	}

	return result.Text(), nil
}
