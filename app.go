package main

import (
	"context"
	"strings"

	"changeme/internal/models"
	"changeme/internal/services"

	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx      context.Context
	visible  bool
	AppCache []models.SearchResult
}

func NewApp() *App {
	return &App{}
}

func (a *App) loadCache() error {
	data, err := services.LoadCache()
	if err != nil {
		return err
	}
	a.AppCache = data
	return nil
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	err := a.loadCache()
	if err != nil {
		println(err.Error())
	}
	runtime.EventsOn(a.ctx, "wails:window-blur", func(optionalData ...interface{}) {
		runtime.WindowHide(a.ctx)
	})
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
