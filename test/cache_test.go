package test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hschne/c7/internal"
)

func withTempCache(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", dir)
	return dir
}

func TestCacheLookupMiss(t *testing.T) {
	withTempCache(t)

	_, ok := internal.CacheLookup("nonexistent")
	if ok {
		t.Fatal("expected cache miss for empty cache")
	}
}

func TestCacheSaveAndLookup(t *testing.T) {
	withTempCache(t)

	internal.CacheSave("rails", "/rails/rails", "Ruby on Rails")

	entry, ok := internal.CacheLookup("rails")
	if !ok {
		t.Fatal("expected cache hit after save")
	}
	if entry.ID != "/rails/rails" {
		t.Errorf("ID = %q, want %q", entry.ID, "/rails/rails")
	}
	if entry.Name != "Ruby on Rails" {
		t.Errorf("Name = %q, want %q", entry.Name, "Ruby on Rails")
	}
}

func TestCacheLookupCaseInsensitive(t *testing.T) {
	withTempCache(t)

	internal.CacheSave("Rails", "/rails/rails", "Ruby on Rails")

	for _, key := range []string{"rails", "Rails", "RAILS"} {
		t.Run(key, func(t *testing.T) {
			if _, ok := internal.CacheLookup(key); !ok {
				t.Errorf("expected cache hit for key %q", key)
			}
		})
	}
}

func TestCacheLookupExpired(t *testing.T) {
	dir := withTempCache(t)

	s := internal.CacheStore{
		"old": {
			ID:   "/old/lib",
			Name: "Old Lib",
			TS:   time.Now().Add(-internal.CacheTTL - time.Hour).Unix(),
		},
	}
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(dir, "c7", "libs.json")
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, data, 0o644); err != nil {
		t.Fatal(err)
	}

	if _, ok := internal.CacheLookup("old"); ok {
		t.Fatal("expected cache miss for expired entry")
	}
}

func TestCacheSaveMultipleKeys(t *testing.T) {
	withTempCache(t)

	internal.CacheSave("rails", "/rails/rails", "Rails")
	internal.CacheSave("next.js", "/vercel/next.js", "Next.js")

	tests := []struct {
		key    string
		wantID string
	}{
		{"rails", "/rails/rails"},
		{"next.js", "/vercel/next.js"},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			entry, ok := internal.CacheLookup(tt.key)
			if !ok {
				t.Fatalf("expected cache hit for %q", tt.key)
			}
			if entry.ID != tt.wantID {
				t.Errorf("ID = %q, want %q", entry.ID, tt.wantID)
			}
		})
	}
}

func TestCacheSaveOverwrites(t *testing.T) {
	withTempCache(t)

	internal.CacheSave("rails", "/old/rails", "Old")
	internal.CacheSave("rails", "/rails/rails", "New")

	entry, ok := internal.CacheLookup("rails")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if entry.ID != "/rails/rails" {
		t.Errorf("ID = %q, want %q", entry.ID, "/rails/rails")
	}
}

func TestCacheClear(t *testing.T) {
	withTempCache(t)

	internal.CacheSave("rails", "/rails/rails", "Rails")
	if err := internal.CacheClear(); err != nil {
		t.Fatalf("CacheClear failed: %v", err)
	}
	if _, ok := internal.CacheLookup("rails"); ok {
		t.Fatal("expected cache miss after clear")
	}
}

func TestCacheClearNoFile(t *testing.T) {
	withTempCache(t)

	if err := internal.CacheClear(); err != nil {
		t.Fatalf("CacheClear on empty cache: %v", err)
	}
}

func TestCacheLoadCorruptFile(t *testing.T) {
	dir := withTempCache(t)

	p := filepath.Join(dir, "c7", "libs.json")
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("not json{{{"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, ok := internal.CacheLookup("anything"); ok {
		t.Fatal("expected cache miss for corrupt file")
	}

	internal.CacheSave("rails", "/rails/rails", "Rails")
	entry, ok := internal.CacheLookup("rails")
	if !ok {
		t.Fatal("expected cache hit after saving over corrupt file")
	}
	if entry.ID != "/rails/rails" {
		t.Errorf("ID = %q, want %q", entry.ID, "/rails/rails")
	}
}
