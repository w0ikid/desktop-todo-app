package wails

import (
	"context"
	"time"
	"github.com/w0ikid/dekstop-todo-app/internal/app"
)

// TaskHandler - адаптер для Wails frontend binding
type TaskHandler struct {
	createTask   app.CreateTask
	updateTask   app.UpdateTask
	completeTask app.CompleteTask
	getTask      app.GetTask
	listTasks    app.ListTasks
	getDashboard app.GetDashboard
	deleteTask   app.DeleteTask
}

func NewTaskHandler(
	createTask app.CreateTask,
	updateTask app.UpdateTask,
	completeTask app.CompleteTask,
	getTask app.GetTask,
	listTasks app.ListTasks,
	getDashboard app.GetDashboard,
	deleteTask app.DeleteTask,
) *TaskHandler {
	return &TaskHandler{
		createTask:   createTask,
		updateTask:   updateTask,
		completeTask: completeTask,
		getTask:      getTask,
		listTasks:    listTasks,
		getDashboard: getDashboard,
		deleteTask:   deleteTask,
	}
}

// Методы для Wails binding - они автоматически будут доступны во frontend

func (h *TaskHandler) CreateTask(title, priority string, dueDate *time.Time) (app.CreateTaskOutput, error) {
	return h.createTask.Execute(context.Background(), app.CreateTaskInput{
		Title:    title,
		Priority: priority,
		DueDate:  dueDate,
	})
}

func (h *TaskHandler) UpdateTask(id string, title, status, priority *string, dueDate *time.Time) error {
	return h.updateTask.Execute(context.Background(), app.UpdateTaskInput{
		ID:       id,
		Title:    title,
		Status:   status,
		Priority: priority,
		DueDate:  dueDate,
	})
}

func (h *TaskHandler) CompleteTask(id string) error {
	return h.completeTask.Execute(context.Background(), app.CompleteTaskInput{ID: id})
}

func (h *TaskHandler) GetTask(id string) (app.GetTaskOutput, error) {
	return h.getTask.Execute(context.Background(), app.GetTaskInput{ID: id})
}

func (h *TaskHandler) ListTasks(status, priority, filter *string) (app.ListTasksOutput, error) {
	return h.listTasks.Execute(context.Background(), app.ListTasksInput{
		Status:   status,
		Priority: priority,
		Filter:   filter,
	})
}

func (h *TaskHandler) GetDashboard() (app.GetDashboardOutput, error) {
	return h.getDashboard.Execute(context.Background())
}

func (h *TaskHandler) DeleteTask(id string) error {
	return h.deleteTask.Execute(context.Background(), app.DeleteTaskInput{ID: id})
}