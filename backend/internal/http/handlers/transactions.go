package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/fin-helper/backend/internal/repo"
	"github.com/google/uuid"
)

type TransactionsHandler struct {
	Store repo.TransactionStore
}

// GET /api/v1/users/{userID}/transactions
func (h *TransactionsHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	userID, err := uuid.Parse(r.PathValue("userID"))
	if err != nil {
		http.Error(w, `{"error":"invalid user id"}`, http.StatusBadRequest)
		return
	}

	items, err := h.Store.ListByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"data": items,
		"meta": map[string]int{"count": len(items)},
	})
}
