package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"whendo/app"
	"whendo/internal/database"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	db, err := database.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	application := app.NewApp(db)

	err = wails.Run(&options.App{
		Title:  "WhenDo",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:        application.Startup,
		OnShutdown:       application.Shutdown,
		OnDomReady:       application.DomReady,
		Bind: []interface{}{
			application,
		},
	})
	if err != nil {
		panic(err)
	}
}
