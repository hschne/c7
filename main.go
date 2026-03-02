package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	baseURL = "https://context7.com/api/v2"
	version = "0.1.0"
)

// --- API types ---

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

// --- HTTP client ---

func apiKey() string {
	return os.Getenv("CONTEXT7_API_KEY")
}

func get(endpoint string, params url.Values) ([]byte, error) {
	u := baseURL + endpoint
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	if k := apiKey(); k != "" {
		req.Header.Set("Authorization", "Bearer "+k)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return body, nil
}

// --- Commands ---

func cmdSearch(args []string) {
	if len(args) < 1 {
		die("Usage: c7 search <library-name> [query]")
	}
	libraryName := args[0]
	query := libraryName
	if len(args) > 1 {
		query = strings.Join(args[1:], " ")
	}

	params := url.Values{
		"libraryName": {libraryName},
		"query":       {query},
	}

	body, err := get("/libs/search", params)
	must(err)

	var resp SearchResponse
	must(json.Unmarshal(body, &resp))
	libs := resp.Results

	if len(libs) == 0 {
		fmt.Println("No libraries found.")
		return
	}

	cacheSave(libraryName, libs[0])

	fmt.Printf("%-30s %-8s %s\n", "ID", "TRUST", "NAME")
	fmt.Println(strings.Repeat("─", 70))
	for _, l := range libs {
		fmt.Printf("%-30s %-8.1f %s\n", l.ID, l.TrustScore, l.Name)
		if l.Description != "" {
			wrapped := wrapText(l.Description, 60)
			for i, line := range wrapped {
				if i == 0 {
					fmt.Printf("%-30s          %s\n", "", line)
				} else {
					fmt.Printf("%-30s          %s\n", "", line)
				}
			}
		}
	}
}

func cmdDocs(args []string) {
	if len(args) < 2 {
		die("Usage: c7 docs <library-id> <query> [--tokens N] [--page N] [--topic TOPIC]")
	}
	libraryID := args[0]
	// collect remaining flags
	var queryParts []string
	tokens := "5000"
	page := "1"
	topic := ""

	i := 1
	for i < len(args) {
		switch args[i] {
		case "--tokens":
			i++
			if i < len(args) {
				tokens = args[i]
			}
		case "--page":
			i++
			if i < len(args) {
				page = args[i]
			}
		case "--topic":
			i++
			if i < len(args) {
				topic = args[i]
			}
		default:
			queryParts = append(queryParts, args[i])
		}
		i++
	}

	query := strings.Join(queryParts, " ")
	if query == "" {
		die("Please provide a query")
	}

	params := url.Values{
		"libraryId": {libraryID},
		"query":     {query},
		"tokens":    {tokens},
		"page":      {page},
	}
	if topic != "" {
		params.Set("topic", topic)
	}

	body, err := get("/context", params)
	must(err)

	// Try JSON array first, fall back to plain text
	var snippets []DocSnippet
	if json.Unmarshal(body, &snippets) == nil {
		for i, s := range snippets {
			if i > 0 {
				fmt.Println(strings.Repeat("─", 70))
			}
			fmt.Printf("## %s\n", s.Title)
			if s.Source != "" {
				fmt.Printf("Source: %s\n\n", s.Source)
			}
			fmt.Println(s.Content)
		}
	} else {
		// Plain text response
		fmt.Println(string(body))
	}
}

func cmdGet(args []string) {
	// Convenience: search then immediately fetch docs
	// c7 get rails "active record scopes"
	if len(args) < 2 {
		die("Usage: c7 get <library-name> <query> [--tokens N]")
	}
	libName := args[0]
	queryArgs := args[1:]

	// Extract --tokens before passing to docs
	tokens := "5000"
	var queryParts []string
	i := 0
	for i < len(queryArgs) {
		if queryArgs[i] == "--tokens" {
			i++
			if i < len(queryArgs) {
				tokens = queryArgs[i]
			}
		} else {
			queryParts = append(queryParts, queryArgs[i])
		}
		i++
	}
	query := strings.Join(queryParts, " ")

	// Step 1: resolve (check cache first)
	var best Library
	if entry, ok := cacheLookup(libName); ok {
		best = Library{ID: entry.ID, Name: entry.Name}
	} else {
		params := url.Values{
			"libraryName": {libName},
			"query":       {query},
		}
		body, err := get("/libs/search", params)
		must(err)

		var resp SearchResponse
		must(json.Unmarshal(body, &resp))

		if len(resp.Results) == 0 {
			die("No library found for: " + libName)
		}

		best = resp.Results[0]
		cacheSave(libName, best)
	}

	// Step 2: fetch docs
	docsParams := url.Values{
		"libraryId": {best.ID},
		"query":     {query},
		"tokens":    {tokens},
	}
	docsBody, err := get("/context", docsParams)
	must(err)

	var snippets []DocSnippet
	if json.Unmarshal(docsBody, &snippets) == nil {
		for i, s := range snippets {
			if i > 0 {
				fmt.Println(strings.Repeat("─", 70))
			}
			fmt.Printf("## %s\n", s.Title)
			if s.Source != "" {
				fmt.Printf("Source: %s\n\n", s.Source)
			}
			fmt.Println(s.Content)
		}
	} else {
		fmt.Println(string(docsBody))
	}
}

func cmdCache(args []string) {
	if len(args) < 1 {
		die("Usage: c7 cache clear")
	}
	switch args[0] {
	case "clear":
		must(cacheClear())
		fmt.Println("Cache cleared.")
	default:
		die("Usage: c7 cache clear")
	}
}

// --- Helpers ---

func wrapText(s string, width int) []string {
	words := strings.Fields(s)
	var lines []string
	var cur strings.Builder
	for _, w := range words {
		if cur.Len() > 0 && cur.Len()+1+len(w) > width {
			lines = append(lines, cur.String())
			cur.Reset()
		}
		if cur.Len() > 0 {
			cur.WriteByte(' ')
		}
		cur.WriteString(w)
	}
	if cur.Len() > 0 {
		lines = append(lines, cur.String())
	}
	return lines
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func die(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func usage() {
	fmt.Printf(`c7 - Context7 CLI (v%s)

Commands:
  c7 search <library-name> [query]
      Search for libraries by name. Optional query improves relevance ranking.
      Example: c7 search rails "active record"

  c7 docs <library-id> <query> [--tokens N] [--page N] [--topic TOPIC]
      Fetch docs for a known library ID.
      Example: c7 docs /rails/rails "how to use scopes"
      Example: c7 docs /vercel/next.js "middleware" --topic routing --page 2

  c7 get <library-name> <query> [--tokens N]
      One-shot: resolve library name then fetch docs. Easiest to use.
      Caches the resolved library ID for faster repeat lookups.
      Example: c7 get hotwire "form submission with turbo"
      Example: c7 get kamal "deploy with secrets" --tokens 8000

  c7 cache clear
      Remove all cached library lookups.

Environment:
  CONTEXT7_API_KEY   Optional API key for higher rate limits.
                     Get one at https://context7.com/dashboard

`, version)
	_ = strconv.Itoa(0) // keep import
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(0)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "search":
		cmdSearch(args)
	case "docs":
		cmdDocs(args)
	case "get":
		cmdGet(args)
	case "cache":
		cmdCache(args)
	case "version", "--version", "-v":
		fmt.Println("c7 version", version)
	case "help", "--help", "-h":
		usage()
	default:
		// If it looks like a library name, treat as `get`
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmd)
		usage()
		os.Exit(1)
	}
}
