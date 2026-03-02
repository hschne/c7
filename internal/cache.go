package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const CacheTTL = 7 * 24 * time.Hour

type CacheEntry struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	TS   int64  `json:"ts"`
}

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
	var s CacheStore
	if json.Unmarshal(data, &s) != nil {
		return nil
	}
	return s
}

func CacheLookup(key string) (CacheEntry, bool) {
	s := cacheLoad()
	if s == nil {
		return CacheEntry{}, false
	}
	e, ok := s[strings.ToLower(key)]
	if !ok {
		return CacheEntry{}, false
	}
	if time.Since(time.Unix(e.TS, 0)) > CacheTTL {
		return CacheEntry{}, false
	}
	return e, true
}

func CacheSave(key, id, name string) {
	p, err := cachePath()
	if err != nil {
		return
	}

	s := cacheLoad()
	if s == nil {
		s = make(CacheStore)
	}

	s[strings.ToLower(key)] = CacheEntry{
		ID:   id,
		Name: name,
		TS:   time.Now().Unix(),
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return
	}
	_ = os.WriteFile(p, data, 0o644)
}

func CacheClear() error {
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
