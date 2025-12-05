package model

type WikipediaSearch struct {
	Query struct {
		Search []struct {
			Title  string `json:"title"`
			PageID int    `json:"pageId"`
		}
	} `json:"query"`
}
