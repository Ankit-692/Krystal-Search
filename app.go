package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"changeme/internal/models"
	"changeme/internal/services"

	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx            context.Context
	AppCache       []models.SearchResult
	watcher        *services.WatcherService
	frontendAssets embed.FS
}

func NewApp(assets embed.FS) *App {
	return &App{
		frontendAssets: assets,
	}
}

func (a *App) loadCache() error {
	data, err := services.LoadCache()
	if err != nil {
		return err
	}
	a.AppCache = data
	return nil
}

func (a *App) rebuildCache() {
	log.Println("[app] rebuilding cache...")
	a.AppCache = services.BuildCache()
	runtime.EventsEmit(a.ctx, "cache:updated")
}

func (a *App) initializeLocalFolder() error {
	targetDir := filepath.Join(os.Getenv("HOME"), models.Dir, "images")

	err := os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		return err
	}

	images := []string{"app.png", "file.png", "folder.png"}

	for _, filename := range images {
		destinationPath := filepath.Join(targetDir, filename)
		if _, err := os.Stat(destinationPath); err == nil {
			continue
		}

		var imageBytes []byte

		embeddedPath := filepath.Join("frontend", "dist", "assets", "images", filename)

		imageBytes, err = a.frontendAssets.ReadFile(embeddedPath)
		if err != nil {
			devPath := filepath.Join("frontend", "src", "assets", "images", filename)
			imageBytes, err = os.ReadFile(devPath)
			if err != nil {
				return err
			}
		}

		err = os.WriteFile(destinationPath, imageBytes, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	err := a.loadCache()
	if err != nil {
		println(err.Error())
	}

	err = a.initializeLocalFolder()
	if err != nil {
		fmt.Printf("Error initializing assets: %v\n", err)
	}

	ws, err := services.NewWatcherService()
	if err != nil {
		log.Println("[app] watcher init error:", err)
	} else {
		ws.OnAppsChanged = func() {
			a.rebuildCache()
		}
		ws.OnIconChanged = func(newTheme string) {
			log.Printf("[app] icon theme changed to %q, rebuilding cache", newTheme)
			a.rebuildCache()
			runtime.EventsEmit(a.ctx, "icon:changed", newTheme)
		}
		ws.Start()
		a.watcher = ws
	}

	runtime.EventsOn(a.ctx, "wails:window-blur", func(optionalData ...interface{}) {
		runtime.WindowHide(a.ctx)
	})
}

func (a *App) shutdown(ctx context.Context) {
	if a.watcher != nil {
		a.watcher.Stop()
	}
}

func (a *App) Search(query string) []models.SearchResult {
	query = strings.ToLower(query)
	result := services.Search(a.AppCache, query)
	return result
}

func (a *App) FileSearch(query string) []models.FileEntry {
	result := services.FileSearch(query)
	return result
}

func (a *App) Launch(item models.SearchResult) {
	services.Launch(item)
}

func (a *App) RunCommand(command string, password string) string {
	return services.RunCommand(command, password)
}

func (a *App) OnSecondInstance(secondInstanceData options.SecondInstanceData) {
	a.showWindow()
}

func (a *App) showWindow() {
	runtime.WindowSetAlwaysOnTop(a.ctx, true)
	runtime.WindowCenter(a.ctx)
	runtime.WindowShow(a.ctx)

	go func() {
		time.Sleep(150 * time.Millisecond)
		_ = exec.Command("wmctrl", "-a", "Krystal-Search").Run()
		time.Sleep(50 * time.Millisecond)
		runtime.EventsEmit(a.ctx, "focus-search")
	}()
}
