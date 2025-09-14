package app

import (
	"context"
	"fmt"

	"github.com/w0ikid/dekstop-todo-app/internal/domain"
)

type GetDashboard struct {
	repo domain.TaskRepository
}

func NewGetDashboard(repo domain.TaskRepository) GetDashboard {
	return GetDashboard{repo: repo}
}

type GetDashboardOutput struct {
	ActiveCount    int            `json:"active_count"`
	CompletedCount int            `json:"completed_count"`
	OverdueCount   int            `json:"overdue_count"`
	DueToday       []*domain.Task `json:"due_today"`
	DueThisWeek    []*domain.Task `json:"due_this_week"`
	RecentTasks    []*domain.Task `json:"recent_tasks"`
}

func (uc GetDashboard) Execute(ctx context.Context) (GetDashboardOutput, error) {
	activeTasks, err := uc.repo.GetByStatus(ctx, domain.StatusActive)
	if err != nil {
		return GetDashboardOutput{}, fmt.Errorf("get active tasks: %w", err)
	}

	completedTasks, err := uc.repo.GetByStatus(ctx, domain.StatusCompleted)
	if err != nil {
		return GetDashboardOutput{}, fmt.Errorf("get completed tasks: %w", err)
	}

	listUC := NewListTasks(uc.repo)

	dueToday, err := listUC.getTasksDueToday(ctx)
	if err != nil {
		return GetDashboardOutput{}, fmt.Errorf("get due today: %w", err)
	}

	dueWeek, err := listUC.getTasksDueThisWeek(ctx)
	if err != nil {
		return GetDashboardOutput{}, fmt.Errorf("get due this week: %w", err)
	}

	overdue, err := listUC.getOverdueTasks(ctx)
	if err != nil {
		return GetDashboardOutput{}, fmt.Errorf("get overdue: %w", err)
	}

	// Последние 5 активных задач
	recent := activeTasks
	if len(recent) > 5 {
		recent = recent[:5]
	}

	return GetDashboardOutput{
		ActiveCount:    len(activeTasks),
		CompletedCount: len(completedTasks),
		OverdueCount:   len(overdue),
		DueToday:       dueToday,
		DueThisWeek:    dueWeek,
		RecentTasks:    recent,
	}, nil
}