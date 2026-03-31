// Package models
package models

type SearchResult struct {
	Title string `json:"Title"`
	Path  string `json:"Path"`
	Icon  string `json:"Icon"`
}

const (
	Dir           = ".krystalSearch/data"
	CacheFileName = "data.json"
)

type FileEntry struct {
	Name  string `json:"Name"`
	Path  string `json:"Path"`
	IsDir bool   `json:"IsDir"`
	Icon  string `json:"Icon"`
}
