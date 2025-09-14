package app

import (
	"context"
	"fmt"

	"github.com/w0ikid/dekstop-todo-app/internal/domain"
)

type GetTask struct {
	repo domain.TaskRepository
}

func NewGetTask(repo domain.TaskRepository) GetTask {
	return GetTask{repo: repo}
}

type GetTaskInput struct {
	ID string `json:"id"`
}

type GetTaskOutput struct {
	Task *domain.Task `json:"task"`
}

func (uc GetTask) Execute(ctx context.Context, in GetTaskInput) (GetTaskOutput, error) {
	task, err := uc.repo.GetByID(ctx, in.ID)
	if err != nil {
		return GetTaskOutput{}, fmt.Errorf("get task: %w", err)
	}

	return GetTaskOutput{Task: task}, nil
}