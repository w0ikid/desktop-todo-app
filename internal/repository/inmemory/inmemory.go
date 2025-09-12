package inmemory

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"github.com/w0ikid/dekstop-todo-app/internal/domain"
)

type InMemoryTaskRepository struct {
	tasks map[string]*domain.Task
	mu    sync.RWMutex
}

func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		tasks: make(map[string]*domain.Task),
	}
}

func (r *InMemoryTaskRepository) GetTaskByID(id string) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	task, ok := r.tasks[id]
	if !ok {
		return nil, errors.New("task not found")
	}
	return task, nil
}

func (r *InMemoryTaskRepository) GetAllTasks() ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	if len(r.tasks) == 0 {
		return []*domain.Task{}, nil
	}

	allTasks := make([]*domain.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		allTasks = append(allTasks, task)
	}
	return allTasks, nil
}

func (r *InMemoryTaskRepository) SaveTask(task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if task.ID == "" {
		return errors.New("task ID cannot be empty")
	}
	r.tasks[task.ID] = task
	return nil
}

func (r *InMemoryTaskRepository) DeleteTask(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[id]; !ok {
		return errors.New("task not found")
	}
	delete(r.tasks, id)
	return nil
}

func (r *InMemoryTaskRepository) GetTasksByStatus(status domain.TaskStatus) ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	filteredTasks := make([]*domain.Task, 0)
	for _, task := range r.tasks {
		if task.Status == status {
			filteredTasks = append(filteredTasks, task)
		}
	}
	return filteredTasks, nil
}

func (r *InMemoryTaskRepository) GetSortedTasks(sortBy string, order string) ([]*domain.Task, error) {
	tasks, err := r.GetAllTasks()
	if err != nil {
		return nil, err
	}
	
	isDesc := strings.ToLower(order) == "desc"
	
	sort.Slice(tasks, func(i, j int) bool {
		var less bool
		switch strings.ToLower(sortBy) {
		case "createdat":
			less = tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
		case "duedate":
			if tasks[i].DueDate == nil {
				return true
			}
			if tasks[j].DueDate == nil {
				return false
			}
			less = tasks[i].DueDate.Before(*tasks[j].DueDate)
		case "priority":
			priorityWeight := func(p domain.TaskPriority) int {
				switch p {
				case domain.PriorityHigh:
					return 3
				case domain.PriorityMedium:
					return 2
				case domain.PriorityLow:
					return 1
				default:
					return 0
				}
			}
			less = priorityWeight(tasks[i].Priority) < priorityWeight(tasks[j].Priority)
		default:
			less = tasks[i].ID < tasks[j].ID
		}
		
		if isDesc {
			return !less
		}
		return less
	})

	return tasks, nil
}