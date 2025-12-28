package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"ecom_test/internal/domain"
)

func (h *TaskHandler) sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func (h *TaskHandler) sendError(w http.ResponseWriter, status int, message string) {
	h.sendJSON(w, status, map[string]string{"error": message})
}

func (h *TaskHandler) handleError(ctx context.Context, w http.ResponseWriter, err error) {
	logger(ctx).Error(err.Error())
	if errors.Is(err, domain.ErrTaskNotFound) {
		h.sendError(w, http.StatusNotFound, "task_not_found")
		return
	}

	if errors.Is(err, domain.ErrEmptyTitle) || errors.Is(err, domain.ErrInvalidID) {
		h.sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	h.sendError(w, http.StatusInternalServerError, "internal_server_error")
}
