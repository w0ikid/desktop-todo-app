package app

import (
	"context"
	"fmt"
	"time"
	"errors"
	"github.com/w0ikid/dekstop-todo-app/internal/domain"
)

type ListTasks struct {
	repo domain.TaskRepository
}

func NewListTasks(repo domain.TaskRepository) ListTasks {
	return ListTasks{repo: repo}
}

type ListTasksInput struct {
	Status   *string `json:"status,omitempty"`
	Priority *string `json:"priority,omitempty"`
	Filter   *string `json:"filter,omitempty"` // "today", "week", "overdue"
}

type ListTasksOutput struct {
	Tasks []*domain.Task `json:"tasks"`
	Total int            `json:"total"`
}

func (uc ListTasks) Execute(ctx context.Context, in ListTasksInput) (ListTasksOutput, error) {
	var tasks []*domain.Task
	var err error

	// Применяем фильтры
	if in.Filter != nil {
		switch *in.Filter {
		case "today":
			tasks, err = uc.getTasksDueToday(ctx)
		case "week":
			tasks, err = uc.getTasksDueThisWeek(ctx)
		case "overdue":
			tasks, err = uc.getOverdueTasks(ctx)
		default:
			return ListTasksOutput{}, errors.New("invalid filter")
		}
	} else if in.Status != nil {
		status := domain.TaskStatus(*in.Status)
		tasks, err = uc.repo.GetByStatus(ctx, status)
	} else {
		tasks, err = uc.repo.GetAll(ctx)
	}

	if err != nil {
		return ListTasksOutput{}, fmt.Errorf("get tasks: %w", err)
	}

	// Фильтр по приоритету
	if in.Priority != nil {
		priority := domain.Priority(*in.Priority)
		filtered := make([]*domain.Task, 0)
		for _, task := range tasks {
			if task.Priority == priority {
				filtered = append(filtered, task)
			}
		}
		tasks = filtered
	}

	return ListTasksOutput{
		Tasks: tasks,
		Total: len(tasks),
	}, nil
}

func (uc ListTasks) getTasksDueToday(ctx context.Context) ([]*domain.Task, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	return uc.repo.GetDueBetween(ctx, startOfDay, endOfDay)
}

func (uc ListTasks) getTasksDueThisWeek(ctx context.Context) ([]*domain.Task, error) {
	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, now.Location())
	endOfWeek := startOfWeek.Add(7 * 24 * time.Hour)

	return uc.repo.GetDueBetween(ctx, startOfWeek, endOfWeek)
}

func (uc ListTasks) getOverdueTasks(ctx context.Context) ([]*domain.Task, error) {
	now := time.Now()
	startTime := time.Date(2000, 1, 1, 0, 0, 0, 0, now.Location())
	
	tasks, err := uc.repo.GetDueBetween(ctx, startTime, now)
	if err != nil {
		return nil, err
	}

	// Только активные просроченные задачи
	overdue := make([]*domain.Task, 0)
	for _, task := range tasks {
		if task.Status == domain.StatusActive && task.IsOverdue() {
			overdue = append(overdue, task)
		}
	}

	return overdue, nil
}