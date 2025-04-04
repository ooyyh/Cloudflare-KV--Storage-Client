package main

import (
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "SRSC-KVN1-Plan",
		Width:  424,
		Height: 305,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},

		//BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		OnStartup: app.startup,
		Bind: []interface{}{
			&Request{},
			app,
			&FileManager{},
		},
		//Windows: &windows.Options{
		//	WebviewIsTransparent: true,
		//	WindowIsTranslucent:  true,
		//	BackdropType:         windows.Acrylic,
		//},

		Frameless: true,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
