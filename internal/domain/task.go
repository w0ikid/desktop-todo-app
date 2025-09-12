package domain

import (
    "time"
)

type TaskStatus string

const (
    TaskStatusActive    TaskStatus = "active"
    TaskStatusCompleted TaskStatus = "completed"
)

type TaskPriority string

const (
    PriorityLow    TaskPriority = "low"
    PriorityMedium TaskPriority = "medium"
    PriorityHigh   TaskPriority = "high"
)

type Task struct {
    ID          string       `json:"id"`
    Title       string       `json:"title"`
    Status      TaskStatus   `json:"status"`
    CreatedAt   time.Time    `json:"createdAt"`
    DueDate     *time.Time   `json:"dueDate,omitempty"`
    Priority    TaskPriority `json:"priority,omitempty"`
}

// Interface for task repository
type TaskRepository interface {
    GetTaskByID(id string) (*Task, error)
    GetAllTasks() ([]*Task, error)
    SaveTask(task *Task) error
    DeleteTask(id string) error
    
    GetTasksByStatus(status TaskStatus) ([]*Task, error)
    GetSortedTasks(sortBy string, order string) ([]*Task, error)
}