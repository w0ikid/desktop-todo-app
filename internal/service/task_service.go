package service

import (
    "errors"
    "time"
    "github.com/google/uuid"
    "github.com/w0ikid/dekstop-todo-app/internal/domain"
)

type TaskService struct {
    repo domain.TaskRepository
}

func NewTaskService(repo domain.TaskRepository) *TaskService {
    return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(title string) (*domain.Task, error) {
    // validation
	if title == "" {
        return nil, errors.New("task title cannot be empty")
    }
    task := &domain.Task{
        ID:        uuid.New().String(),
        Title:     title,
        Status:    domain.TaskStatusActive,
        CreatedAt: time.Now(),
    }
    err := s.repo.SaveTask(task)
    return task, err
}

func (s *TaskService) MarkTaskAsCompleted(id string) error {
    task, err := s.repo.GetTaskByID(id)
    if err != nil {
        return err
    }
    task.Status = domain.TaskStatusCompleted
    return s.repo.SaveTask(task)
}

func (s *TaskService) MarkTaskAsActive(id string) error {
    task, err := s.repo.GetTaskByID(id)
    if err != nil {
        return err
    }
    task.Status = domain.TaskStatusActive
    return s.repo.SaveTask(task)
}

func (s *TaskService) DeleteTask(id string) error {
    return s.repo.DeleteTask(id)
}

func (s *TaskService) GetAllTasks() ([]*domain.Task, error) {
    return s.repo.GetAllTasks()
}