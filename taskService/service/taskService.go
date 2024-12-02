package service

import (
	"context"
	"errors"
	"taskService/domain"
	"taskService/repository"
)

type TaskService struct {
	repo *repository.TaskRepo
}

func NewTaskService(r *repository.TaskRepo) *TaskService {
	return &TaskService{repo: r}
}

func (s *TaskService) Create(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	if task.Name == "" || task.Description == "" {
		return nil, errors.New("task name and description must be provided")
	}

	// Call the Create method from the repository to insert the task
	createdTask, err := s.repo.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	return createdTask, nil
}
