package domain

import (
	"errors"
	"fmt"
)

var (
	ErrTaskNotFound    = errors.New("task not found")
	ErrEmptyTask       = errors.New("empty task")
	ErrEmptyTitle      = errors.New("task title cannot be empty")
	ErrTaskAlreadyDone = errors.New("task is already completed")
	ErrInvalidID       = errors.New("invalid task identifier")
)

type TaskError struct {
	Op  string
	ID  int
	Err error
}

func (e *TaskError) Error() string {
	return fmt.Sprintf("ERROR: operation %s on task %d: %v", e.Op, e.ID, e.Err)
}

func (e *TaskError) Unwrap() error {
	return e.Err
}

func Wrap(err error, op string, id int) error {
	if err == nil {
		return nil
	}
	return &TaskError{
		Op:  op,
		ID:  id,
		Err: err,
	}
}
