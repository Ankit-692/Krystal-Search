package main

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx      context.Context
	visible  bool
	AppCache []SearchResult
}

type SearchResult struct {
	Title string
	Path  string
}

func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	paths := []string{"/usr/share/applications", filepath.Join(os.Getenv("HOME"), ".local/share/applications")}
	for _, dir := range paths {
		files, _ := os.ReadDir(dir)
		for _, f := range files {
			if filepath.Ext(f.Name()) == ".desktop" {
				name := a.parseDesktopFile(filepath.Join(dir, f.Name()))
				if name != "" {
					a.AppCache = append(a.AppCache, SearchResult{Title: name, Path: f.Name()})
				}
			}
		}
	}

	runtime.EventsOn(a.ctx, "wails:window-blur", func(optionalData ...interface{}) {
		runtime.WindowHide(a.ctx)
	})
}

func (a *App) parseDesktopFile(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Name=") {
			return strings.TrimPrefix(line, "Name=")
		}
	}
	return ""
}

func (a *App) Search(query string) []SearchResult {
	query = strings.ToLower(query)
	var filtered []SearchResult
	for _, app := range a.AppCache {
		if strings.Contains(strings.ToLower(app.Title), query) {
			filtered = append(filtered, app)
		}
	}
	return filtered
}

func (a *App) Launch(item SearchResult) {
	go func() {
		var cmd *exec.Cmd
		if strings.HasSuffix(item.Path, ".desktop") {
			cmd = exec.Command("gtk-launch", item.Path)
		} else {
			cmd = exec.Command("xdg-open", item.Path)
		}

		_ = cmd.Start()
	}()
}

func (a *App) OnSecondInstance(secondInstanceData options.SecondInstanceData) {
	a.ToggleWindow()
}

func (a *App) ToggleWindow() {
	if a.visible {
		runtime.WindowHide(a.ctx)
		a.visible = false
	} else {
		runtime.WindowSetAlwaysOnTop(a.ctx, true)
		runtime.WindowCenter(a.ctx)
		runtime.WindowShow(a.ctx)
		runtime.EventsEmit(a.ctx, "focus-search")
		a.visible = true
	}
}
