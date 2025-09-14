package app

import (
	"context"
	"fmt"
	"time"

	"github.com/w0ikid/dekstop-todo-app/internal/domain"
)

type UpdateTask struct {
	repo domain.TaskRepository
}

func NewUpdateTask(repo domain.TaskRepository) UpdateTask {
	return UpdateTask{repo: repo}
}

type UpdateTaskInput struct {
	ID       string     `json:"id"`
	Title    *string    `json:"title,omitempty"`
	Status   *string    `json:"status,omitempty"`
	Priority *string    `json:"priority,omitempty"`
	DueDate  *time.Time `json:"due_date,omitempty"`
}

func (uc UpdateTask) Execute(ctx context.Context, in UpdateTaskInput) error {
	task, err := uc.repo.GetByID(ctx, in.ID)
	if err != nil {
		return fmt.Errorf("get task: %w", err)
	}

	if in.Title != nil {
		task.Title = *in.Title
	}
	if in.Status != nil {
		task.Status = domain.TaskStatus(*in.Status)
	}
	if in.Priority != nil {
		task.Priority = domain.Priority(*in.Priority)
	}
	if in.DueDate != nil {
		task.DueDate = in.DueDate
	}

	if err := task.IsValid(); err != nil {
		return fmt.Errorf("validate task: %w", err)
	}

	if err := uc.repo.Save(ctx, task); err != nil {
		return fmt.Errorf("save task: %w", err)
	}

	return nil
}