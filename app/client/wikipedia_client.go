package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/WikiScrolls/pagerank/app/model"
)

type WikipediaClient struct {
	httpClient *http.Client
}

func NewWikipediaClient() *WikipediaClient {
	return &WikipediaClient{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (w *WikipediaClient) GetRandomArticles(ctx context.Context, articleCount int) (*model.WikipediaResponse, error) {
	params := url.Values{
		"action":       {"query"},
		"format":       {"json"},
		"generator":    {"random"},
		"grnnamespace": {"0"},
		"grnlimit":     {strconv.Itoa(articleCount)},
		"prop":         {"extracts|info|pageimages"},
		"exintro":      {"1"},
		"exlimit":      {"max"},
		"exsentences":  {"10"},
		"explaintext":  {"1"},
		"origin":       {"*"},
		"variant":      {"en"},
		"piprop":       {"thumbnail"},
		"pithumbsize":  {"800"},
	}

	body, err := w.fetchWikipedia(ctx, params.Encode())
	if err != nil {
		return nil, err
	}

	return unMarshalWikipediaResponse(body)
}

func (w *WikipediaClient) FetchByTitles(ctx context.Context, titles []string) (*model.WikipediaResponse, error) {
	if len(titles) == 0 {
		return nil, errors.New("no titles given")
	}

	var titleParam string = titles[0]
	for i := 1; i < len(titles); i++ {
		titleParam += "|" + titles[i]
	}

	params := url.Values{
		"action":      {"query"},
		"format":      {"json"},
		"prop":        {"extracts|info|pageimages"},
		"exintro":     {"1"},
		"exlimit":     {"max"},
		"exsentences": {"10"},
		"explaintext": {"1"},
		"origin":      {"*"},
		"variant":     {"en"},
		"piprop":      {"thumbnail"},
		"pithumbsize": {"800"},
		"titles":      {titleParam},
	}

	body, err := w.fetchWikipedia(ctx, params.Encode())
	if err != nil {
		return nil, err
	}

	return unMarshalWikipediaResponse(body)
}

func (w *WikipediaClient) FetchBySearch(ctx context.Context, search string) (*model.WikipediaSearch, error) {
	params := url.Values{
		"action":   {"query"},
		"format":   {"json"},
		"list":     {"search"},
		"srprop":   {""},
		"srsearch": {search},
	}

	body, err := w.fetchWikipedia(ctx, params.Encode())
	if err != nil {
		return nil, err
	}

	return unMarshalWikipediaSearch(body)
}

func (w *WikipediaClient) FetchByIDs(ctx context.Context, ids []string) (*model.WikipediaResponse, error) {
	if len(ids) == 0 {
		return nil, errors.New("no ids given")
	}

	var idParam string = ids[0]
	for i := 1; i < len(ids); i++ {
		idParam += "|" + ids[i]
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
		"pageids":     {idParam},
	}

	body, err := w.fetchWikipedia(ctx, params.Encode())
	if err != nil {
		return nil, err
	}

	return unMarshalWikipediaResponse(body)
}

func (w *WikipediaClient) fetchWikipedia(ctx context.Context, params string) ([]byte, error) {
	reqURL := "https://en.wikipedia.org/w/api.php?" + params

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "WikiScrolls/1.0 (; nadzhiff@gmail.com) Go-http-client/1.1")

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func unMarshalWikipediaResponse(jsonRaw []byte) (*model.WikipediaResponse, error) {
	wikiResponse := model.WikipediaResponse{}
	json.Unmarshal(jsonRaw, &wikiResponse)
	return &wikiResponse, nil
}

func unMarshalWikipediaSearch(jsonRaw []byte) (*model.WikipediaSearch, error) {
	wikiSearch := model.WikipediaSearch{}
	json.Unmarshal(jsonRaw, &wikiSearch)
	return &wikiSearch, nil
}
