package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const BaseURL = "https://context7.com/api/v2"

type Library struct {
	ID          string   `json:"id"`
	Name        string   `json:"title"`
	Description string   `json:"description"`
	TrustScore  float64  `json:"trustScore"`
	Versions    []string `json:"versions"`
}

type SearchResponse struct {
	Results []Library `json:"results"`
}

type DocSnippet struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Source  string `json:"source"`
}

type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL:    BaseURL,
		APIKey:     os.Getenv("CONTEXT7_API_KEY"),
		HTTPClient: http.DefaultClient,
	}
}

func (c *Client) get(endpoint string, params url.Values) ([]byte, error) {
	u := c.BaseURL + endpoint
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return body, nil
}

func (c *Client) Search(libraryName, query string) ([]Library, error) {
	params := url.Values{
		"libraryName": {libraryName},
		"query":       {query},
	}
	body, err := c.get("/libs/search", params)
	if err != nil {
		return nil, err
	}
	var resp SearchResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing search response: %w", err)
	}
	return resp.Results, nil
}

func (c *Client) FetchDocs(libraryID, query, tokens, page, topic string) ([]byte, error) {
	params := url.Values{
		"libraryId": {libraryID},
		"query":     {query},
		"tokens":    {tokens},
		"page":      {page},
	}
	if topic != "" {
		params.Set("topic", topic)
	}
	return c.get("/context", params)
}
