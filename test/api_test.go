package test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/hschne/c7/internal"
)

func TestSearchResponseUnmarshal(t *testing.T) {
	const payload = `{
		"results": [
			{
				"id": "/rails/rails",
				"title": "Ruby on Rails",
				"description": "A web framework",
				"trustScore": 7.6,
				"versions": ["v8.0.0", "v7.2.2"]
			},
			{
				"id": "/dry-rb/dry-rails",
				"title": "Dry Rails",
				"description": "Dry-rb integration for Rails",
				"trustScore": 9.2,
				"versions": []
			}
		]
	}`

	t.Run("into SearchResponse", func(t *testing.T) {
		var resp internal.SearchResponse
		if err := json.Unmarshal([]byte(payload), &resp); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		if got := len(resp.Results); got != 2 {
			t.Fatalf("got %d results, want 2", got)
		}
		first := resp.Results[0]
		if first.ID != "/rails/rails" {
			t.Errorf("ID = %q, want %q", first.ID, "/rails/rails")
		}
		if first.Name != "Ruby on Rails" {
			t.Errorf("Name = %q, want %q", first.Name, "Ruby on Rails")
		}
		if first.TrustScore != 7.6 {
			t.Errorf("TrustScore = %v, want 7.6", first.TrustScore)
		}
		if got := len(first.Versions); got != 2 {
			t.Errorf("got %d versions, want 2", got)
		}
	})

	t.Run("bare array fails", func(t *testing.T) {
		var libs []internal.Library
		if err := json.Unmarshal([]byte(payload), &libs); err == nil {
			t.Fatal("expected error unmarshaling wrapped response into []Library, got nil")
		}
	})

	t.Run("empty results", func(t *testing.T) {
		var resp internal.SearchResponse
		if err := json.Unmarshal([]byte(`{"results":[]}`), &resp); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		if len(resp.Results) != 0 {
			t.Errorf("got %d results, want 0", len(resp.Results))
		}
	})

	t.Run("null versions", func(t *testing.T) {
		var resp internal.SearchResponse
		if err := json.Unmarshal([]byte(`{"results":[{"id":"/x","title":"X","trustScore":1.0,"versions":null}]}`), &resp); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		if resp.Results[0].Versions != nil {
			t.Errorf("expected nil versions, got %v", resp.Results[0].Versions)
		}
	})
}

func TestDocSnippetUnmarshal(t *testing.T) {
	const payload = `[
		{"title": "Scopes", "content": "Use scope to define queries.", "source": "https://example.com/scopes"},
		{"title": "Validations", "content": "Use validates to check data.", "source": ""}
	]`

	var snippets []internal.DocSnippet
	if err := json.Unmarshal([]byte(payload), &snippets); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if got := len(snippets); got != 2 {
		t.Fatalf("got %d snippets, want 2", got)
	}
	if snippets[0].Title != "Scopes" {
		t.Errorf("Title = %q, want %q", snippets[0].Title, "Scopes")
	}
	if snippets[1].Source != "" {
		t.Errorf("Source = %q, want empty", snippets[1].Source)
	}
}

func TestDocSnippetPlainTextFallback(t *testing.T) {
	plainText := "### Example heading\n\nSome documentation content."
	var snippets []internal.DocSnippet
	if err := json.Unmarshal([]byte(plainText), &snippets); err == nil {
		t.Fatal("expected error unmarshaling plain text as JSON, got nil")
	}
}

func TestLibraryJSONRoundtrip(t *testing.T) {
	lib := internal.Library{
		ID:          "/test/lib",
		Name:        "Test Lib",
		Description: "A library",
		TrustScore:  8.5,
		Versions:    []string{"v1.0.0"},
	}

	data, err := json.Marshal(lib)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	if !strings.Contains(string(data), `"title"`) {
		t.Errorf("expected JSON key \"title\", got: %s", data)
	}
	if strings.Contains(string(data), `"name"`) {
		t.Errorf("unexpected JSON key \"name\" in: %s", data)
	}

	var got internal.Library
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if got.ID != lib.ID || got.Name != lib.Name || got.TrustScore != lib.TrustScore {
		t.Errorf("roundtrip mismatch:\n got: %+v\nwant: %+v", got, lib)
	}
}
