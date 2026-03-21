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
