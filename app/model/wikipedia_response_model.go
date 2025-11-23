package model

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
