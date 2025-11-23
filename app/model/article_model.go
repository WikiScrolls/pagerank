package model

type Article struct {
	Id           string `json:"id"`
	Title        string `json:"title"`
	WikipediaUrl string `json:"wikipediaUrl"`
	Content      string `json:"content"`
	Thumbnail    string `json:"thumbnail"`
}
