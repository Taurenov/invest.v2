package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/fin-helper/backend/internal/repo"
	"github.com/google/uuid"
)

type GoalsHandler struct {
	Store repo.GoalStore
}

func (h *GoalsHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.PathValue("userID"))
	if err != nil {
		http.Error(w, `{"error":"invalid user id"}`, http.StatusBadRequest)
		return
	}
	items, err := h.Store.ListGoalsByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"data": items})
}
