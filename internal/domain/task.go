package domain

import (
	"errors"
	"strings"
	"time"
)

type TaskStatus string

const (
	StatusActive    TaskStatus = "active"
	StatusCompleted TaskStatus = "completed"
)

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// Domain errors - sentinel errors для четкой сигнализации
var (
	ErrTaskNotFound    = errors.New("task not found")
	ErrInvalidTitle    = errors.New("invalid title")
	ErrInvalidStatus   = errors.New("invalid status")
	ErrInvalidPriority = errors.New("invalid priority")
)

type Task struct {
	ID        string
	Title     string
	Status    TaskStatus
	CreatedAt time.Time
	DueDate   *time.Time
	Priority  Priority
}

// Фабрика для создания новой задачи
func NewTask(title string, priority Priority, dueDate *time.Time) (*Task, error) {
	task := &Task{
		ID:        generateID(), // перенесем в domain
		Title:     title,
		Status:    StatusActive,
		CreatedAt: time.Now(),
		DueDate:   dueDate,
		Priority:  priority,
	}

	if err := task.IsValid(); err != nil {
		return nil, err
	}

	return task, nil
}

func (t *Task) IsValid() error {
	if strings.TrimSpace(t.Title) == "" {
		return ErrInvalidTitle
	}

	if len(t.Title) > 255 {
		return ErrInvalidTitle
	}

	if t.Status != StatusActive && t.Status != StatusCompleted {
		return ErrInvalidStatus
	}

	if t.Priority != PriorityLow && t.Priority != PriorityMedium && t.Priority != PriorityHigh {
		return ErrInvalidPriority
	}

	return nil
}

func (t *Task) Complete() error {
	if t.Status == StatusCompleted {
		return errors.New("task already completed")
	}
	t.Status = StatusCompleted
	return nil
}

func (t *Task) IsOverdue() bool {
	if t.DueDate == nil || t.Status == StatusCompleted {
		return false
	}
	return t.DueDate.Before(time.Now())
}

func generateID() string {
	return "task_" + time.Now().Format("20060102150405000")
}