package persistance

import (
	"context"
	"ecom_test/internal/domain"
	"ecom_test/internal/domain/entity"
	"sync"
)

type TaskRepository struct {
	mu        sync.RWMutex
	data      map[int]entity.Task
	currentID int
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		data:      make(map[int]entity.Task),
		currentID: 0,
	}
}

func (r *TaskRepository) GetByID(ctx context.Context, id int) (*entity.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if v, ok := r.data[id]; ok {
		taskCopy := v
		return &taskCopy, nil
	}
	return nil, domain.ErrTaskNotFound
}

func (r *TaskRepository) Create(ctx context.Context, task *entity.Task) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	task.ID = r.currentID
	r.data[task.ID] = *task
	r.currentID++

	return task.ID, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; !ok {
		return domain.ErrTaskNotFound
	}

	delete(r.data, id)
	return nil
}

func (r *TaskRepository) Update(ctx context.Context, task *entity.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[task.ID]; !ok {
		return domain.ErrTaskNotFound
	}

	r.data[task.ID] = *task
	return nil
}

// сделал проверку контекста только здесь потому что остальные методы работают мнгновенно или почти мнгновенно
func (r *TaskRepository) GetAll(ctx context.Context) ([]entity.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]entity.Task, 0, len(r.data))
	for _, v := range r.data {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			tasks = append(tasks, v)
		}
	}
	return tasks, nil
}
