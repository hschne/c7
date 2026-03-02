package internal

import (
	"encoding/json"
	"fmt"
	"strings"
)

// WrapText wraps s at word boundaries to fit within width columns.
func WrapText(s string, width int) []string {
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

// PrintDocs outputs documentation from a raw API response body.
// It tries to parse as JSON snippets first, falling back to plain text.
func PrintDocs(body []byte) {
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
		fmt.Println(string(body))
	}
}
