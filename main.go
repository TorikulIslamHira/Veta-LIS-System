package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"veta-lis/internal/app"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create application with options
	application := app.NewApp("config.yaml")

	// Create application with wails
	err := wails.Run(&options.App{
		Title:      "Veta LIS - Veterinary Laboratory Information System",
		Width:      1400,
		Height:     900,
		MinWidth:   1000,
		MinHeight:  600,
		OnStartup:  application.Startup,
		OnShutdown: application.Shutdown,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Bind: []interface{}{
			application,
		},
	})

	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
