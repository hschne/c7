package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const cacheTTL = 7 * 24 * time.Hour

type CacheEntry struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	TS   int64  `json:"ts"`
}

// CacheStore is the on-disk format: map of lowercased library name → entry.
type CacheStore map[string]CacheEntry

func cachePath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "c7", "libs.json"), nil
}

func cacheLoad() CacheStore {
	p, err := cachePath()
	if err != nil {
		return nil
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return nil
	}
	var store CacheStore
	if json.Unmarshal(data, &store) != nil {
		return nil
	}
	return store
}

func cacheLookup(key string) (CacheEntry, bool) {
	store := cacheLoad()
	if store == nil {
		return CacheEntry{}, false
	}
	entry, ok := store[strings.ToLower(key)]
	if !ok {
		return CacheEntry{}, false
	}
	if time.Since(time.Unix(entry.TS, 0)) > cacheTTL {
		return CacheEntry{}, false
	}
	return entry, true
}

func cacheSave(key string, lib Library) {
	p, err := cachePath()
	if err != nil {
		return
	}

	store := cacheLoad()
	if store == nil {
		store = make(CacheStore)
	}

	store[strings.ToLower(key)] = CacheEntry{
		ID:   lib.ID,
		Name: lib.Name,
		TS:   time.Now().Unix(),
	}

	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return
	}

	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return
	}
	_ = os.WriteFile(p, data, 0o644)
}

func cacheClear() error {
	p, err := cachePath()
	if err != nil {
		return err
	}
	err = os.Remove(p)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
