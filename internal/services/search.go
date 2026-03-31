// Package services
package services

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	"changeme/internal/models"
	"changeme/internal/utility"
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

func FileSearch(query string) []models.FileEntry {
	fmt.Println("query received:", strconv.Quote(query))
	bin := "plocate"
	if _, err := exec.LookPath(bin); err != nil {
		bin = "locate"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []string{"--limit", "100", "--basename"}

	args = append(args, query)

	cmd := exec.CommandContext(ctx, bin, args...)
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	var results []models.FileEntry
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		path := scanner.Text()
		if path == "" {
			continue
		}

		base := filepath.Base(path)
		contains := strings.Contains(base, query)
		if !contains {
			continue
		}

		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		iconType := map[bool]string{true: "folder", false: ""}[info.IsDir()]
		results = append(results, models.FileEntry{
			Name:  filepath.Base(path),
			Path:  path,
			IsDir: info.IsDir(),
			Icon:  utility.IconToDataURL(utility.ResolveIcon(iconType)),
		})
	}

	strip := func(s string) string {
		var b strings.Builder
		for _, r := range s {
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				b.WriteRune(r)
			}
		}
		return b.String()
	}

	slices.SortStableFunc(results, func(a, b models.FileEntry) int {
		aName := strip(strings.ToLower(strings.TrimSuffix(filepath.Base(a.Path), filepath.Ext(a.Path))))
		bName := strip(strings.ToLower(strings.TrimSuffix(filepath.Base(b.Path), filepath.Ext(b.Path))))
		qStripped := strip(strings.ToLower(query))

		aExact := aName == qStripped
		bExact := bName == qStripped

		switch {
		case aExact && !bExact:
			return -1
		case !aExact && bExact:
			return 1
		default:
			return 0
		}
	})

	return results
}
