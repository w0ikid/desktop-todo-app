package app

import (
	"context"
	"fmt"

	"github.com/w0ikid/dekstop-todo-app/internal/domain"
)

type CompleteTask struct {
	repo domain.TaskRepository
}

func NewCompleteTask(repo domain.TaskRepository) CompleteTask {
	return CompleteTask{repo: repo}
}

type CompleteTaskInput struct {
	ID string `json:"id"`
}

func (uc CompleteTask) Execute(ctx context.Context, in CompleteTaskInput) error {
	return uc.repo.WithTx(ctx, func(repo domain.TaskRepository) error {
		task, err := repo.GetByID(ctx, in.ID)
		if err != nil {
			return fmt.Errorf("get task: %w", err)
		}

		if err := task.Complete(); err != nil {
			return fmt.Errorf("complete task: %w", err)
		}

		if err := repo.Save(ctx, task); err != nil {
			return fmt.Errorf("save task: %w", err)
		}

		return nil
	})
}