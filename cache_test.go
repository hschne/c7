package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// withCacheDir points the cache at a temp directory for the duration of a test.
// It sets XDG_CACHE_HOME and restores the original value on cleanup.
func withCacheDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	prev := os.Getenv("XDG_CACHE_HOME")
	t.Setenv("XDG_CACHE_HOME", dir)
	t.Cleanup(func() { os.Setenv("XDG_CACHE_HOME", prev) })
	return dir
}

func TestCacheLookupMiss(t *testing.T) {
	withCacheDir(t)

	_, ok := cacheLookup("nonexistent")
	if ok {
		t.Fatal("expected cache miss for empty cache")
	}
}

func TestCacheSaveAndLookup(t *testing.T) {
	withCacheDir(t)

	lib := Library{ID: "/rails/rails", Name: "Ruby on Rails"}
	cacheSave("rails", lib)

	entry, ok := cacheLookup("rails")
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
	withCacheDir(t)

	cacheSave("Rails", Library{ID: "/rails/rails", Name: "Ruby on Rails"})

	for _, key := range []string{"rails", "Rails", "RAILS"} {
		t.Run(key, func(t *testing.T) {
			_, ok := cacheLookup(key)
			if !ok {
				t.Errorf("expected cache hit for key %q", key)
			}
		})
	}
}

func TestCacheLookupExpired(t *testing.T) {
	dir := withCacheDir(t)

	// Write an entry with a timestamp older than the TTL.
	store := CacheStore{
		"old": {
			ID:   "/old/lib",
			Name: "Old Lib",
			TS:   time.Now().Add(-cacheTTL - time.Hour).Unix(),
		},
	}
	data, err := json.Marshal(store)
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

	_, ok := cacheLookup("old")
	if ok {
		t.Fatal("expected cache miss for expired entry")
	}
}

func TestCacheSaveMultipleKeys(t *testing.T) {
	withCacheDir(t)

	cacheSave("rails", Library{ID: "/rails/rails", Name: "Rails"})
	cacheSave("next.js", Library{ID: "/vercel/next.js", Name: "Next.js"})

	tests := []struct {
		key    string
		wantID string
	}{
		{"rails", "/rails/rails"},
		{"next.js", "/vercel/next.js"},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			entry, ok := cacheLookup(tt.key)
			if !ok {
				t.Fatalf("expected cache hit for %q", tt.key)
			}
			if entry.ID != tt.wantID {
				t.Errorf("ID = %q, want %q", entry.ID, tt.wantID)
			}
		})
	}
}

func TestCacheSaveOverwritesExisting(t *testing.T) {
	withCacheDir(t)

	cacheSave("rails", Library{ID: "/old/rails", Name: "Old"})
	cacheSave("rails", Library{ID: "/rails/rails", Name: "New"})

	entry, ok := cacheLookup("rails")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if entry.ID != "/rails/rails" {
		t.Errorf("ID = %q, want %q (overwrite failed)", entry.ID, "/rails/rails")
	}
}

func TestCacheClear(t *testing.T) {
	withCacheDir(t)

	cacheSave("rails", Library{ID: "/rails/rails", Name: "Rails"})

	if err := cacheClear(); err != nil {
		t.Fatalf("cacheClear failed: %v", err)
	}

	_, ok := cacheLookup("rails")
	if ok {
		t.Fatal("expected cache miss after clear")
	}
}

func TestCacheClearNoFile(t *testing.T) {
	withCacheDir(t)

	// Clearing when no cache file exists should not error.
	if err := cacheClear(); err != nil {
		t.Fatalf("cacheClear on empty cache failed: %v", err)
	}
}

func TestCacheLoadCorruptFile(t *testing.T) {
	dir := withCacheDir(t)

	p := filepath.Join(dir, "c7", "libs.json")
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("not json{{{"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, ok := cacheLookup("anything")
	if ok {
		t.Fatal("expected cache miss for corrupt file")
	}

	// Saving should overwrite the corrupt file gracefully.
	cacheSave("rails", Library{ID: "/rails/rails", Name: "Rails"})
	entry, ok := cacheLookup("rails")
	if !ok {
		t.Fatal("expected cache hit after saving over corrupt file")
	}
	if entry.ID != "/rails/rails" {
		t.Errorf("ID = %q, want %q", entry.ID, "/rails/rails")
	}
}
