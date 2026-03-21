// Package services
package services

import (
	"strings"

	"changeme/internal/models"
)

func Search(data []models.SearchResult, query string) []models.SearchResult {
	var filtered []models.SearchResult
	for _, app := range data {
		if strings.Contains(strings.ToLower(app.Title), query) {
			filtered = append(filtered, app)
		}
	}
	return filtered
}
