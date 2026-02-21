package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Krystal-Search",
		Width:  800,
		Height: 500,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Frameless:        true,
		AlwaysOnTop:      true,
		BackgroundColour: &options.RGBA{R: 10, G: 10, B: 10, A: 255},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId:               "krystal-search-unique-id",
			OnSecondInstanceLaunch: app.OnSecondInstance,
		},
		Linux: &linux.Options{
			WindowIsTranslucent: true,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
