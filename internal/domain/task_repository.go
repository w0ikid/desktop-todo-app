package domain

import (
	"context"
	"time"
)

type TaskRepository interface {
	Save(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, id string) (*Task, error)
	GetAll(ctx context.Context) ([]*Task, error)
	GetByStatus(ctx context.Context, status TaskStatus) ([]*Task, error)
	GetDueBetween(ctx context.Context, startDate, endDate time.Time) ([]*Task, error)
	Delete(ctx context.Context, id string) error
	WithTx(ctx context.Context, fn func(repo TaskRepository) error) error
}