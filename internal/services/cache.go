// Package services
package services

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"changeme/internal/models"
	"changeme/internal/utility"
)

func BuildCache() []models.SearchResult {
	var results []models.SearchResult
	paths := []string{
		"/usr/share/applications",
		filepath.Join(os.Getenv("HOME"), ".local/share/applications"),
	}
	for _, dir := range paths {
		files, _ := os.ReadDir(dir)
		for _, f := range files {
			if filepath.Ext(f.Name()) == ".desktop" {
				name, iconName := parseDesktopFile(filepath.Join(dir, f.Name()))
				if name != "" {
					iconPath := utility.ResolveIcon(iconName)
					results = append(results, models.SearchResult{Title: name, Path: f.Name(), Icon: utility.IconToDataURL(iconPath)})
				}
			}
		}
	}
	SaveCache(results)

	return results
}

func SaveCache(cache []models.SearchResult) {
	dataFilePath := filepath.Join(os.Getenv("HOME"), models.Dir)
	fullPath := filepath.Join(dataFilePath, "data.json")
	data, _ := json.Marshal(cache)
	_ = os.WriteFile(fullPath, data, 0o755)
}

func LoadCache() ([]models.SearchResult, error) {
	var data []models.SearchResult
	dataFilePath := filepath.Join(os.Getenv("HOME"), models.Dir)
	fullPath := filepath.Join(dataFilePath, "data.json")
	if _, err := os.Stat(dataFilePath); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(dataFilePath, 0o755)
		return nil, err
	}

	if _, err := os.Stat(filepath.Join(dataFilePath, "data.json")); errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(fullPath)
		if err != nil {
			return nil, err
		}
		data = BuildCache()
		return data, nil
	}

	file, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func parseDesktopFile(path string) (name, icon string) {
	file, err := os.Open(path)
	if err != nil {
		return "", ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "Name=") && name == "":
			name = strings.TrimPrefix(line, "Name=")
		case strings.HasPrefix(line, "Icon=") && icon == "":
			icon = strings.TrimPrefix(line, "Icon=")
		}
		if name != "" && icon != "" {
			break
		}
	}
	return name, icon
}
