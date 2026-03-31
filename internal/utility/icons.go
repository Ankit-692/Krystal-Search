// Package services
package utility

import (
	"encoding/base64"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var iconCategories = []string{"apps", "applications"}

func getCurrentIconTheme() string {
	out, err := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "icon-theme").Output()
	if err != nil {
		return "hicolor"
	}

	return strings.Trim(strings.TrimSpace(string(out)), "'")
}

func ResolveIcon(iconName string) string {
	theme := getCurrentIconTheme()
	if iconName == "" {
		return ""
	}

	if filepath.IsAbs(iconName) {
		if _, err := os.Stat(iconName); err == nil {
			return iconName
		}
		return ""
	}

	iconBaseDirs := []string{
		filepath.Join(os.Getenv("HOME"), ".local/share/icons"),
		"/usr/local/share/icons",
		"/usr/share/icons",
	}

	type probe struct {
		theme string
		size  string
		exts  []string
	}

	probes := []probe{
		{theme, "128x128", []string{".png", ".svg", ".xpm"}},
		{"hicolor", "128x128", []string{".png", ".svg", ".xpm"}},
		{theme, "scalable", []string{".svg"}},
		{"hicolor", "scalable", []string{".svg"}},
	}

	if theme == "hicolor" {
		probes = []probe{
			{"hicolor", "128x128", []string{".png", ".svg", ".xpm"}},
			{"hicolor", "scalable", []string{".svg"}},
		}
	}

	if iconName == "folder" {
		iconCategories = []string{"places"}
	}

	for _, p := range probes {
		for _, base := range iconBaseDirs {
			for _, cat := range iconCategories {
				dir := filepath.Join(base, p.theme, p.size, cat)
				for _, ext := range p.exts {
					candidate := filepath.Join(dir, iconName+ext)
					if _, err := os.Stat(candidate); err == nil {
						return candidate
					}
				}
			}
		}
	}

	return ""
}

func ResolveFolderIcon() {
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
	case ".xpm":
		return ""
	default:
		return ""
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	return "data:" + mime + ";base64," + encoded
}
