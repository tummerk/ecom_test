package server

import (
	"context"
	"ecom_test/internal/domain/entity"
	"ecom_test/internal/server/dto"
	"encoding/json"
	"net/http"
	"strconv"
)

type TaskService interface {
	GetByID(ctx context.Context, id int) (*entity.Task, error)
	GetAll(ctx context.Context) ([]entity.Task, error)
	Create(ctx context.Context, task *entity.Task) (int, error)
	Update(ctx context.Context, task *entity.Task) error
	Delete(ctx context.Context, id int) error
}

type TaskHandler struct {
	service TaskService
}

func NewTaskHandler(service TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "invalid_body")
		return
	}

	task := &entity.Task{
		Title:       req.Title,
		Description: req.Description,
	}

	id, err := h.service.Create(r.Context(), task)
	if err != nil {
		h.handleError(r.Context(), w, err)
		return
	}

	h.sendJSON(w, http.StatusCreated, dto.CreateTaskResponse{ID: id})
}

func (h *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.GetAll(r.Context())
	if err != nil {
		h.handleError(r.Context(), w, err)
		return
	}

	resp := dto.GetAllTasksResponse{
		Tasks: make([]dto.TaskListItemResponse, 0, len(tasks)),
	}
	for _, t := range tasks {
		resp.Tasks = append(resp.Tasks, dto.TaskListItemResponse{
			ID:          t.ID,
			Title:       t.Title,
			IsCompleted: t.IsCompleted,
		})
	}

	h.sendJSON(w, http.StatusOK, resp)
}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	task, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(r.Context(), w, err)
		return
	}

	h.sendJSON(w, http.StatusOK, dto.GetTaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		IsCompleted: task.IsCompleted,
	})
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	var req dto.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "invalid_body")
		return
	}

	task := &entity.Task{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		IsCompleted: req.IsCompleted,
	}

	if err := h.service.Update(r.Context(), task); err != nil {
		h.handleError(r.Context(), w, err)
		return
	}

	h.sendJSON(w, http.StatusOK, dto.UpdateTaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		IsCompleted: task.IsCompleted,
	})
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	if err := h.service.Delete(r.Context(), id); err != nil {
		h.handleError(r.Context(), w, err)
		return
	}

	h.sendJSON(w, http.StatusOK, dto.DeleteTaskResponse{Status: "success"})
}

func (h *TaskHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetAll(w, r)
		case http.MethodPost:
			h.Create(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/todos/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetByID(w, r)
		case http.MethodPut:
			h.Update(w, r)
		case http.MethodDelete:
			h.Delete(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
