package main

import (
	"context"
	"fmt"

	"github.com/w0ikid/dekstop-todo-app/internal/domain"
	"github.com/w0ikid/dekstop-todo-app/internal/service"
)

// App struct
type App struct {
	ctx context.Context
	taskService *service.TaskService
}

// NewApp creates a new App application struct
func NewApp(taskService *service.TaskService) *App {
	return &App{
		taskService: taskService,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// Methods to interact with the task service



// AddTaskFromFrontend adds a new task with the given title
func (a *App) AddTaskFromFrontend(title string) (*domain.Task, error) {
    return a.taskService.CreateTask(title)
}

// GetTasksFromFrontend retrieves all tasks
func (a *App) GetTasksFromFrontend() ([]*domain.Task, error) {
    return a.taskService.GetAllTasks()
}

// MarkTaskCompletedFromFrontend marks a task as completed
func (a *App) MarkTaskCompletedFromFrontend(id string) error {
    return a.taskService.MarkTaskAsCompleted(id)
}

// DeleteTaskFromFrontend deletes a task
func (a *App) DeleteTaskFromFrontend(id string) error {
    return a.taskService.DeleteTask(id)
}

