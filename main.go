package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/w0ikid/dekstop-todo-app/internal/repository/inmemory"
	"github.com/w0ikid/dekstop-todo-app/internal/service"

	_ "github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	repo := inmemory.NewInMemoryTaskRepository()
	service := service.NewTaskService(repo)

	// Create an instance of the app structure
	app := NewApp(service)

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "dekstop-todo-app",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
