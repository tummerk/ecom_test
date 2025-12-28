package service

import (
	"context"
	"ecom_test/internal/domain"
	"ecom_test/internal/domain/entity"
)

type TaskRepository interface {
	GetByID(ctx context.Context, id int) (*entity.Task, error)
	GetAll(ctx context.Context) ([]entity.Task, error)
	Create(ctx context.Context, task *entity.Task) (int, error)
	Update(ctx context.Context, task *entity.Task) error
	Delete(ctx context.Context, id int) error
}

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) GetByID(ctx context.Context, id int) (*entity.Task, error) {
	if id < 0 {
		return nil, domain.Wrap(domain.ErrInvalidID, "GetByID", id)
	}

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.Wrap(err, "GetByID", id)
	}
	return task, nil
}

func (s *TaskService) GetAll(ctx context.Context) ([]entity.Task, error) {
	tasks, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, domain.Wrap(err, "GetAll", 0)
	}
	return tasks, nil
}

func (s *TaskService) Create(ctx context.Context, task *entity.Task) (int, error) {
	if task == nil {
		return 0, domain.Wrap(domain.ErrEmptyTask, "Create", 0)
	}

	if task.Title == "" {
		return 0, domain.Wrap(domain.ErrEmptyTitle, "Create", 0)
	}

	id, err := s.repo.Create(ctx, task)
	if err != nil {
		return 0, domain.Wrap(err, "Create", 0)
	}
	return id, nil
}

func (s *TaskService) Update(ctx context.Context, task *entity.Task) error {
	if task == nil {
		return domain.Wrap(domain.ErrEmptyTask, "Update", 0)
	}
	if task.ID < 0 {
		return domain.Wrap(domain.ErrInvalidID, "Update", task.ID)
	}
	if task.Title == "" {
		return domain.Wrap(domain.ErrEmptyTitle, "Update", task.ID)
	}

	err := s.repo.Update(ctx, task)
	if err != nil {
		return domain.Wrap(err, "Update", task.ID)
	}
	return nil
}

func (s *TaskService) Delete(ctx context.Context, id int) error {
	if id < 0 {
		return domain.Wrap(domain.ErrInvalidID, "Delete", id)
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return domain.Wrap(err, "Delete", id)
	}
	return nil
}
