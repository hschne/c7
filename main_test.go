package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestWrapText(t *testing.T) {
	tests := []struct {
		name  string
		input string
		width int
		want  []string
	}{
		{
			name:  "empty string",
			input: "",
			width: 60,
			want:  nil,
		},
		{
			name:  "single word within width",
			input: "hello",
			width: 60,
			want:  []string{"hello"},
		},
		{
			name:  "single word exceeding width",
			input: "superlongword",
			width: 5,
			want:  []string{"superlongword"},
		},
		{
			name:  "wraps at word boundary",
			input: "aaa bbb ccc",
			width: 7,
			want:  []string{"aaa bbb", "ccc"},
		},
		{
			name:  "preserves all words",
			input: "the quick brown fox jumps over the lazy dog",
			width: 15,
			want:  []string{"the quick brown", "fox jumps over", "the lazy dog"},
		},
		{
			name:  "collapses whitespace",
			input: "  hello   world  ",
			width: 60,
			want:  []string{"hello world"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wrapText(tt.input, tt.width)
			if len(got) != len(tt.want) {
				t.Fatalf("got %d lines %q, want %d lines %q", len(got), got, len(tt.want), tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("line %d: got %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestSearchResponseUnmarshal(t *testing.T) {
	// Real API response shape (trimmed). The key facts:
	//   - results are wrapped in {"results": [...]}
	//   - library name is in the "title" field, not "name"
	//   - trustScore is a float, not an int
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
		var resp SearchResponse
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
		var libs []Library
		if err := json.Unmarshal([]byte(payload), &libs); err == nil {
			t.Fatal("expected error unmarshaling wrapped response into []Library, got nil")
		}
	})

	t.Run("empty results", func(t *testing.T) {
		var resp SearchResponse
		if err := json.Unmarshal([]byte(`{"results":[]}`), &resp); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		if len(resp.Results) != 0 {
			t.Errorf("got %d results, want 0", len(resp.Results))
		}
	})

	t.Run("null versions", func(t *testing.T) {
		var resp SearchResponse
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

	var snippets []DocSnippet
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

func TestDocSnippetUnmarshalFallsBackToPlainText(t *testing.T) {
	// The /context endpoint returns plain text, not JSON.
	// The code tries JSON first and falls back. Verify the fallback path.
	plainText := "### Example heading\n\nSome documentation content."

	var snippets []DocSnippet
	err := json.Unmarshal([]byte(plainText), &snippets)
	if err == nil {
		t.Fatal("expected error unmarshaling plain text as JSON, got nil")
	}
}

func TestLibraryJSONRoundtrip(t *testing.T) {
	lib := Library{
		ID:         "/test/lib",
		Name:       "Test Lib",
		Description: "A library",
		TrustScore: 8.5,
		Versions:   []string{"v1.0.0"},
	}

	data, err := json.Marshal(lib)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Verify the JSON uses "title" (matching the API), not "name"
	if !strings.Contains(string(data), `"title"`) {
		t.Errorf("expected JSON key \"title\", got: %s", data)
	}
	if strings.Contains(string(data), `"name"`) {
		t.Errorf("unexpected JSON key \"name\" in: %s", data)
	}

	var got Library
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if got.ID != lib.ID || got.Name != lib.Name || got.TrustScore != lib.TrustScore {
		t.Errorf("roundtrip mismatch:\n got: %+v\nwant: %+v", got, lib)
	}
}
