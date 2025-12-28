package service

import (
	"context"
	"ecom_test/internal/domain"
	"ecom_test/internal/domain/entity"
	"errors"
	"testing"
)

type MockTaskRepository struct {
	GetByIDFunc func(ctx context.Context, id int) (*entity.Task, error)
	CreateFunc  func(ctx context.Context, task *entity.Task) (int, error)
	GetAllFunc  func(ctx context.Context) ([]entity.Task, error)
	UpdateFunc  func(ctx context.Context, task *entity.Task) error
	DeleteFunc  func(ctx context.Context, id int) error
}

func (m *MockTaskRepository) GetByID(ctx context.Context, id int) (*entity.Task, error) {
	return m.GetByIDFunc(ctx, id)
}
func (m *MockTaskRepository) Create(ctx context.Context, task *entity.Task) (int, error) {
	return m.CreateFunc(ctx, task)
}
func (m *MockTaskRepository) GetAll(ctx context.Context) ([]entity.Task, error) {
	return m.GetAllFunc(ctx)
}
func (m *MockTaskRepository) Update(ctx context.Context, task *entity.Task) error {
	return m.UpdateFunc(ctx, task)
}
func (m *MockTaskRepository) Delete(ctx context.Context, id int) error {
	return m.DeleteFunc(ctx, id)
}

func TestTaskService_Create(t *testing.T) {
	tests := []struct {
		name    string
		task    *entity.Task
		mockFn  func(ctx context.Context, task *entity.Task) (int, error)
		wantID  int
		wantErr error
	}{
		{
			name: "Success create task",
			task: &entity.Task{Title: "Купить хлеб"},
			mockFn: func(ctx context.Context, task *entity.Task) (int, error) {
				return 1, nil
			},
			wantID:  1,
			wantErr: nil,
		},
		{
			name:    "Error nil task",
			task:    nil,
			mockFn:  nil,
			wantID:  0,
			wantErr: domain.ErrEmptyTask, // ПРОВЕРКА НА NIL
		},
		{
			name:    "Error empty title",
			task:    &entity.Task{Title: ""},
			mockFn:  nil,
			wantID:  0,
			wantErr: domain.ErrEmptyTitle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockTaskRepository{CreateFunc: tt.mockFn}
			svc := NewTaskService(repo)

			gotID, err := svc.Create(context.Background(), tt.task)

			if gotID != tt.wantID {
				t.Errorf("Create() gotID = %v, want %v", gotID, tt.wantID)
			}

			if tt.wantErr != nil {
				if err == nil || !errors.Is(err, tt.wantErr) {
					t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else if err != nil {
				t.Errorf("Create() unexpected error = %v", err)
			}
		})
	}
}

func TestTaskService_GetAll(t *testing.T) {
	dbErr := errors.New("db error")
	tests := []struct {
		name    string
		mockFn  func(ctx context.Context) ([]entity.Task, error)
		wantLen int
		wantErr error
	}{
		{
			name: "Success get all",
			mockFn: func(ctx context.Context) ([]entity.Task, error) {
				return []entity.Task{{ID: 1}, {ID: 2}}, nil
			},
			wantLen: 2,
			wantErr: nil,
		},
		{
			name: "Repo error",
			mockFn: func(ctx context.Context) ([]entity.Task, error) {
				return nil, dbErr
			},
			wantLen: 0,
			wantErr: dbErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockTaskRepository{GetAllFunc: tt.mockFn}
			svc := NewTaskService(repo)
			tasks, err := svc.GetAll(context.Background())

			if tt.wantErr != nil {
				if err == nil || !errors.Is(err, tt.wantErr) {
					t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if len(tasks) != tt.wantLen {
				t.Errorf("Expected len %d, got %d", tt.wantLen, len(tasks))
			}
		})
	}
}

func TestTaskService_Update(t *testing.T) {
	errNotFound := errors.New("not found")

	tests := []struct {
		name    string
		task    *entity.Task
		mockFn  func(ctx context.Context, task *entity.Task) error
		wantErr error
	}{
		{
			name: "Success update",
			task: &entity.Task{ID: 1, Title: "Updated"},
			mockFn: func(ctx context.Context, task *entity.Task) error {
				return nil
			},
			wantErr: nil,
		},
		{
			name:    "Error nil task", // ДОБАВИЛИ КЕЙС ДЛЯ NIL
			task:    nil,
			mockFn:  nil,
			wantErr: domain.ErrEmptyTask,
		},
		{
			name:    "Error invalid ID",
			task:    &entity.Task{ID: -1, Title: "Bad ID"},
			wantErr: domain.ErrInvalidID,
		},
		{
			name:    "Error empty title",
			task:    &entity.Task{ID: 1, Title: ""},
			wantErr: domain.ErrEmptyTitle,
		},
		{
			name: "Error task not found in repo",
			task: &entity.Task{ID: 1, Title: "Ghost"},
			mockFn: func(ctx context.Context, task *entity.Task) error {
				return errNotFound
			},
			wantErr: errNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockTaskRepository{UpdateFunc: tt.mockFn}
			svc := NewTaskService(repo)
			err := svc.Update(context.Background(), tt.task)

			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("Expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestTaskService_Delete(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		mockFn  func(ctx context.Context, id int) error
		wantErr error
	}{
		{
			name: "Success delete",
			id:   1,
			mockFn: func(ctx context.Context, id int) error {
				return nil
			},
			wantErr: nil,
		},
		{
			name:    "Error invalid ID",
			id:      -10,
			wantErr: domain.ErrInvalidID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockTaskRepository{DeleteFunc: tt.mockFn}
			svc := NewTaskService(repo)
			err := svc.Delete(context.Background(), tt.id)

			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("Expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestTaskService_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		mockFn  func(ctx context.Context, id int) (*entity.Task, error)
		wantErr error
	}{
		{
			name: "Success get task by id",
			id:   1,
			mockFn: func(ctx context.Context, id int) (*entity.Task, error) {
				return &entity.Task{ID: 1, Title: "Test"}, nil
			},
			wantErr: nil,
		},
		{
			name:    "Error invalid ID",
			id:      -5,
			mockFn:  nil,
			wantErr: domain.ErrInvalidID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockTaskRepository{GetByIDFunc: tt.mockFn}
			svc := NewTaskService(repo)

			_, err := svc.GetByID(context.Background(), tt.id)

			if tt.wantErr != nil && (err == nil || !errors.Is(err, tt.wantErr)) {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
