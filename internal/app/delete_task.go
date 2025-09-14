package app

import (
	"context"
	"fmt"

	"github.com/w0ikid/dekstop-todo-app/internal/domain"
)

type DeleteTask struct {
	repo domain.TaskRepository
}

func NewDeleteTask(repo domain.TaskRepository) DeleteTask {
	return DeleteTask{repo: repo}
}

type DeleteTaskInput struct {
	ID string `json:"id"`
}

func (uc DeleteTask) Execute(ctx context.Context, in DeleteTaskInput) error {
	// Проверяем существование
	_, err := uc.repo.GetByID(ctx, in.ID)
	if err != nil {
		return fmt.Errorf("get task: %w", err)
	}

	if err := uc.repo.Delete(ctx, in.ID); err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	return nil
}