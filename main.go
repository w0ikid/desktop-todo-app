package main

import (
	"context"
	"embed"
	"fmt"

	adapter "github.com/w0ikid/dekstop-todo-app/internal/adapters/wails"
	"github.com/w0ikid/dekstop-todo-app/internal/app"
	db "github.com/w0ikid/dekstop-todo-app/internal/db/sqlc"
	"github.com/w0ikid/dekstop-todo-app/internal/infra/postgres"
	"github.com/w0ikid/dekstop-todo-app/internal/util"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

type App struct {
	ctx context.Context
}

//go:embed all:frontend/dist
var assets embed.FS

func NewApp() *App {
	return &App{}
}

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	// Load config
	cfg := util.InitConfig(util.CleanenvLoader{}, "config.yml")

	// Connect to DB
	conn, err := pgxpool.New(context.Background(), cfg.Database.DSN())
	if err != nil {
		panic("cannot connect to db: " + err.Error())
	}
	defer conn.Close()

	// Queries
	queries := db.New(conn)

	// Repository
	taskRepo := postgres.NewTaskRepository(queries, conn)

	// Use cases
	createTask := app.NewCreateTask(taskRepo)
	updateTask := app.NewUpdateTask(taskRepo)
	completeTask := app.NewCompleteTask(taskRepo)
	getTask := app.NewGetTask(taskRepo)
	listTasks := app.NewListTasks(taskRepo)
	getDashboard := app.NewGetDashboard(taskRepo)
	deleteTask := app.NewDeleteTask(taskRepo)

	// TaskHandler
	taskHandler := adapter.NewTaskHandler(
		createTask, updateTask, completeTask,
		getTask, listTasks, getDashboard, deleteTask,
	)

	appInstance := NewApp()

	// Run Wails
	err = wails.Run(&options.App{
		Title:  "dekstop-todo-app",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			appInstance.ctx = ctx
		},
		Bind: []interface{}{
			taskHandler, // биндим TaskHandler напрямую, чтобы фронтенд видел методы
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
