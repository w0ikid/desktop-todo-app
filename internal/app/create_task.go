package app

import (
	"context"
	"fmt"
	"time"

	"github.com/w0ikid/dekstop-todo-app/internal/domain"
)

type CreateTask struct {
	repo domain.TaskRepository
}

func NewCreateTask(repo domain.TaskRepository) CreateTask {
	return CreateTask{repo: repo}
}

type CreateTaskInput struct {
	Title    string     `json:"title"`
	Priority string     `json:"priority"`
	DueDate  *time.Time `json:"due_date,omitempty"`
}

type CreateTaskOutput struct {
	ID string `json:"id"`
}

func (uc CreateTask) Execute(ctx context.Context, in CreateTaskInput) (CreateTaskOutput, error) {
	priority := domain.Priority(in.Priority)
	if priority == "" {
		priority = domain.PriorityMedium
	}

	task, err := domain.NewTask(in.Title, priority, in.DueDate)
	if err != nil {
		return CreateTaskOutput{}, fmt.Errorf("create task: %w", err)
	}

	if err := uc.repo.Save(ctx, task); err != nil {
		return CreateTaskOutput{}, fmt.Errorf("save task: %w", err)
	}

	return CreateTaskOutput{ID: task.ID}, nil
}