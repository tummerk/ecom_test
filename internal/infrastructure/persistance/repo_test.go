package persistance

import (
	"context"
	"ecom_test/internal/domain"
	"ecom_test/internal/domain/entity"
	"errors"
	"testing"
)

func TestTaskRepository_CRUD(t *testing.T) {
	repo := NewTaskRepository()
	ctx := context.Background()

	task := &entity.Task{Title: "Test Task", Description: "Description"}
	id, err := repo.Create(ctx, task)

	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}
	if id != 0 {
		t.Errorf("Expected first ID to be 0, got %d", id)
	}

	savedTask, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get task by ID: %v", err)
	}
	if savedTask.Title != task.Title {
		t.Errorf("Expected title %s, got %s", task.Title, savedTask.Title)
	}

	savedTask.Title = "Updated Title"
	err = repo.Update(ctx, savedTask)
	if err != nil {
		t.Fatalf("Failed to update task: %v", err)
	}

	updatedTask, _ := repo.GetByID(ctx, id)
	if updatedTask.Title != "Updated Title" {
		t.Errorf("Update failed: expected 'Updated Title', got '%s'", updatedTask.Title)
	}

	tasks, err := repo.GetAll(ctx)
	if err != nil || len(tasks) != 1 {
		t.Errorf("GetAll failed: expected 1 task, got %d", len(tasks))
	}

	err = repo.Delete(ctx, id)
	if err != nil {
		t.Fatalf("Failed to delete task: %v", err)
	}

	_, err = repo.GetByID(ctx, id)
	if !errors.Is(err, domain.ErrTaskNotFound) {
		t.Errorf("Expected ErrTaskNotFound after deletion, got %v", err)
	}
}

func TestTaskRepository_Errors(t *testing.T) {
	repo := NewTaskRepository()
	ctx := context.Background()

	t.Run("Update non-existent task", func(t *testing.T) {
		err := repo.Update(ctx, &entity.Task{ID: 999, Title: "Ghost"})
		if !errors.Is(err, domain.ErrTaskNotFound) {
			t.Errorf("Expected ErrTaskNotFound, got %v", err)
		}
	})

	t.Run("Delete non-existent task", func(t *testing.T) {
		err := repo.Delete(ctx, 999)
		if !errors.Is(err, domain.ErrTaskNotFound) {
			t.Errorf("Expected ErrTaskNotFound, got %v", err)
		}
	})
}
