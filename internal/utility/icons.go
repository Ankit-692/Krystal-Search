package utility

import (
	"bufio"
	"changeme/internal/models"
	"encoding/base64"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var iconBaseDirs = []string{
	filepath.Join(os.Getenv("HOME"), ".local/share/icons"),
	"/usr/local/share/icons",
	"/usr/share/icons",
	"/var/lib/icons",
}

var sizes = []string{"scalable", "48x48", "48", "64x64", "64"}
var extensions = []string{".png", ".svg"}

// Get the current icon theme safely across different Desktop Environments
func getCurrentIconTheme() string {
	// 1. Try GNOME/Cinnamon/MATE (Ubuntu, Linux Mint)
	if out, err := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "icon-theme").Output(); err == nil {
		theme := strings.Trim(strings.TrimSpace(string(out)), "'\"")
		if theme != "" {
			return theme
		}
	}

	// 2. Try XFCE fallback (Xubuntu)
	if out, err := exec.Command("xfconf-query", "-c", "xsettings", "-p", "/Net/IconThemeName").Output(); err == nil {
		theme := strings.TrimSpace(string(out))
		if theme != "" {
			return theme
		}
	}

	// 3. Try KDE fallback (Kubuntu - reads from config file directly)
	kdeConfig := filepath.Join(os.Getenv("HOME"), ".config", "kdeglobals")
	if file, err := os.Open(kdeConfig); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "Theme=") {
				return strings.TrimSpace(strings.Split(line, "=")[1])
			}
		}
	}

	return "hicolor"
}

// Tracks theme inheritance recursively so we don't drop icons missing from the primary theme
func getThemeFallbacks(theme string) []string {
	themes := []string{theme}
	current := theme

	for i := 0; i < 3; i++ {
		parent := parseInheritedTheme(current)
		if parent == "" || parent == "hicolor" {
			break
		}
		themes = append(themes, parent)
		current = parent
	}

	if theme != "hicolor" {
		themes = append(themes, "hicolor")
	}
	return themes
}

func parseInheritedTheme(theme string) string {
	for _, base := range iconBaseDirs {
		themePath := filepath.Join(base, theme, "index.theme")
		file, err := os.Open(themePath)
		if err != nil {
			continue
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "Inherits=") {
				val := strings.TrimSpace(strings.Split(line, "=")[1])
				return strings.Split(val, ",")[0]
			}
		}
	}
	return ""
}

func ResolveIcon(iconName string) string {
	defaultIcon := filepath.Join(os.Getenv("HOME"), models.Dir, "images", "app.png")
	if iconName == "" {
		return defaultIcon
	}

	if filepath.IsAbs(iconName) {
		if _, err := os.Stat(iconName); err == nil {
			return iconName
		}
		for _, ext := range extensions {
			if _, err := os.Stat(iconName + ext); err == nil {
				return iconName + ext
			}
		}
	}

	activeTheme := getCurrentIconTheme()
	themesToSearch := getThemeFallbacks(activeTheme)
	categories := []string{"apps", "applications"}

	for _, theme := range themesToSearch {
		for _, base := range iconBaseDirs {
			for _, size := range sizes {
				for _, cat := range categories {
					// Variation 1: theme/size/category (e.g., Mint-L-Brown/128x128/apps)
					candidate1 := filepath.Join(base, theme, size, cat)
					// Variation 2: theme/category/size (e.g., Mint-L-Brown/apps/48)
					candidate2 := filepath.Join(base, theme, cat, size)

					for _, ext := range extensions {
						if _, err := os.Stat(filepath.Join(candidate1, iconName+ext)); err == nil {
							return filepath.Join(candidate1, iconName+ext)
						}
						if _, err := os.Stat(filepath.Join(candidate2, iconName+ext)); err == nil {
							return filepath.Join(candidate2, iconName+ext)
						}
					}
				}
			}
		}
	}

	pixmapDir := "/usr/share/pixmaps"
	for _, ext := range extensions {
		candidate := filepath.Join(pixmapDir, iconName+ext)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	return defaultIcon
}

func ResolveFolderIcon(iconName string) string {
	defaultFolderIconPath := filepath.Join(os.Getenv("HOME"), models.Dir, "images", "folder.png")

	if iconName == "" {
		return defaultFolderIconPath
	}

	activeTheme := getCurrentIconTheme()
	themesToSearch := getThemeFallbacks(activeTheme)
	categories := []string{"places"}

	for _, theme := range themesToSearch {
		for _, base := range iconBaseDirs {
			for _, size := range sizes {
				for _, cat := range categories {
					// Variation 1: theme/size/category
					candidate1 := filepath.Join(base, theme, size, cat)
					// Variation 2: theme/category/size (Matches your Mint-L-Brown/places/48 path)
					candidate2 := filepath.Join(base, theme, cat, size)

					for _, ext := range extensions {
						if _, err := os.Stat(filepath.Join(candidate1, iconName+ext)); err == nil {
							return filepath.Join(candidate1, iconName+ext)
						}
						if _, err := os.Stat(filepath.Join(candidate2, iconName+ext)); err == nil {
							return filepath.Join(candidate2, iconName+ext)
						}
					}
				}
			}
		}
	}

	return defaultFolderIconPath
}

// 3. Resolve File/Mimetype Icons (e.g., "text-plain", "application-pdf")
func ResolveFileIcon(mimetype string) string {
	defaultFileIconPath := filepath.Join(os.Getenv("HOME"), models.Dir, "images", "file.png")

	if mimetype == "" {
		return defaultFileIconPath
	}

	sanitizedMime := strings.ReplaceAll(mimetype, "/", "-")
	activeTheme := getCurrentIconTheme()
	themesToSearch := getThemeFallbacks(activeTheme)
	categories := []string{"mimetypes"}

	for _, theme := range themesToSearch {
		for _, base := range iconBaseDirs {
			for _, size := range sizes {
				for _, cat := range categories {
					// Variation 1: theme/size/category
					candidate1 := filepath.Join(base, theme, size, cat)
					// Variation 2: theme/category/size
					candidate2 := filepath.Join(base, theme, cat, size)

					for _, ext := range extensions {
						if _, err := os.Stat(filepath.Join(candidate1, sanitizedMime+ext)); err == nil {
							return filepath.Join(candidate1, sanitizedMime+ext)
						}
						if _, err := os.Stat(filepath.Join(candidate2, sanitizedMime+ext)); err == nil {
							return filepath.Join(candidate2, sanitizedMime+ext)
						}
					}
				}
			}
		}
	}
	return defaultFileIconPath
}

func IconToDataURL(iconPath string) string {
	if iconPath == "" {
		return ""
	}

	data, err := os.ReadFile(iconPath)
	if err != nil {
		return ""
	}

	ext := strings.ToLower(filepath.Ext(iconPath))
	var mime string
	switch ext {
	case ".png":
		mime = "image/png"
	case ".svg":
		mime = "image/svg+xml"
	default:
		return ""
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	return "data:" + mime + ";base64," + encoded
}
