package main

import (
	"context"
	"embed"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"localProxy/service"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure

	// Create application with options
	err := wails.Run(&options.App{
		Title:     service.AppName,
		Width:     1024,
		Height:    768,
		MinWidth:  800,
		MinHeight: 600,
		Frameless: false,
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: false,
				HideTitle:                  false,
				FullSizeContent:            false,
			},
			About: &mac.AboutInfo{
				Title:   service.AppName,
				Message: "© 2023 Jerry Bendy",
				Icon:    nil,
			},
		},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup: func(ctx context.Context) {
			service.Init(ctx)
		},
		Bind: []interface{}{
			&service.Hosts{},
			&service.Service{},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
