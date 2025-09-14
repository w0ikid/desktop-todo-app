package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/w0ikid/dekstop-todo-app/internal/db/sqlc"
	"github.com/w0ikid/dekstop-todo-app/internal/domain"
)

type taskRepository struct {
	queries *db.Queries
	pool      *pgxpool.Pool
}

func NewTaskRepository(queries *db.Queries, pool *pgxpool.Pool) domain.TaskRepository {
	return &taskRepository{queries: queries, pool: pool}
}

func (r *taskRepository) Save(ctx context.Context, task *domain.Task) error {
	params := db.SaveTaskParams{
		ID:       task.ID,
		Title:    task.Title,
		Status:   string(task.Status),
		Priority: string(task.Priority),
		CreatedAt: pgtype.Timestamp{
			Time:  task.CreatedAt,
			Valid: true,
		},
	}

	if task.DueDate != nil {
		params.DueDate = pgtype.Timestamp{
			Time:  *task.DueDate,
			Valid: true,
		}
	}

	return r.queries.SaveTask(ctx, params)
}

func (r *taskRepository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	dbTask, err := r.queries.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTaskNotFound
		}
		return nil, err
	}

	return r.convertDBTaskToDomain(dbTask), nil
}

func (r *taskRepository) GetAll(ctx context.Context) ([]*domain.Task, error) {
	dbTasks, err := r.queries.GetAllTasks(ctx)
	if err != nil {
		return nil, err
	}

	tasks := make([]*domain.Task, 0, len(dbTasks))
	for _, dbTask := range dbTasks {
		tasks = append(tasks, r.convertDBTaskToDomain(dbTask))
	}

	return tasks, nil
}

func (r *taskRepository) GetByStatus(ctx context.Context, status domain.TaskStatus) ([]*domain.Task, error) {
	dbTasks, err := r.queries.GetTasksByStatus(ctx, string(status))
	if err != nil {
		return nil, err
	}

	tasks := make([]*domain.Task, 0, len(dbTasks))
	for _, dbTask := range dbTasks {
		tasks = append(tasks, r.convertDBTaskToDomain(dbTask))
	}

	return tasks, nil
}

func (r *taskRepository) GetDueBetween(ctx context.Context, startDate, endDate time.Time) ([]*domain.Task, error) {
	params := db.GetTasksDueBetweenParams{
		DueDate: pgtype.Timestamp{
			Time:  startDate,
			Valid: true,
		},
		DueDate_2: pgtype.Timestamp{
			Time:  endDate,
			Valid: true,
		},
	}

	dbTasks, err := r.queries.GetTasksDueBetween(ctx, params)
	if err != nil {
		return nil, err
	}

	tasks := make([]*domain.Task, 0, len(dbTasks))
	for _, dbTask := range dbTasks {
		tasks = append(tasks, r.convertDBTaskToDomain(dbTask))
	}

	return tasks, nil
}

func (r *taskRepository) Delete(ctx context.Context, id string) error {
	return r.queries.DeleteTask(ctx, id)
}

func (r *taskRepository) WithTx(ctx context.Context, fn func(repo domain.TaskRepository) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	txRepo := &taskRepository{
		queries: r.queries.WithTx(tx),
		pool:      r.pool,
	}

	if err := fn(txRepo); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *taskRepository) convertDBTaskToDomain(dbTask db.Task) *domain.Task {
	task := &domain.Task{
		ID:        dbTask.ID,
		Title:     dbTask.Title,
		Status:    domain.TaskStatus(dbTask.Status),
		Priority:  domain.Priority(dbTask.Priority),
		CreatedAt: dbTask.CreatedAt.Time,
	}

	if dbTask.DueDate.Valid {
		task.DueDate = &dbTask.DueDate.Time
	}

	return task
}