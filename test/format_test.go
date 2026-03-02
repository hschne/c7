package test

import (
	"testing"

	"github.com/hschne/c7/internal"
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
			got := internal.WrapText(tt.input, tt.width)
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
