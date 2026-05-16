package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fin-helper/backend/internal/alerts"
	"github.com/fin-helper/backend/internal/domain"
	"github.com/fin-helper/backend/internal/http/middleware"
	"github.com/fin-helper/backend/internal/repo"
	"github.com/google/uuid"
)

type MeTransactions struct {
	Store repo.TransactionStore
}

func (h *MeTransactions) List(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	items, err := h.Store.ListByUser(r.Context(), uid)
	if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"data": items})
}

func (h *MeTransactions) Create(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	var req struct {
		CategoryID  *uuid.UUID             `json:"category_id"`
		Kind        domain.TransactionKind `json:"kind"`
		Amount      float64                `json:"amount"`
		Currency    string                 `json:"currency"`
		Description string                 `json:"description"`
		OccurredAt  time.Time              `json:"occurred_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Amount <= 0 {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if req.Currency == "" {
		req.Currency = "RUB"
	}
	if req.OccurredAt.IsZero() {
		req.OccurredAt = time.Now()
	}
	t, err := h.Store.CreateTransaction(r.Context(), uid, repo.CreateTransactionInput{
		CategoryID: req.CategoryID, Kind: req.Kind, Amount: req.Amount,
		Currency: req.Currency, Description: req.Description, OccurredAt: req.OccurredAt,
	})
	if err != nil {
		http.Error(w, `{"error":"could not create"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"data": t})
}

func (h *MeTransactions) Update(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	txID, err := uuid.Parse(r.PathValue("txID"))
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}
	var req struct {
		CategoryID  *uuid.UUID              `json:"category_id"`
		Kind        *domain.TransactionKind `json:"kind"`
		Amount      *float64                `json:"amount"`
		Currency    *string                 `json:"currency"`
		Description *string                 `json:"description"`
		OccurredAt  *time.Time              `json:"occurred_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"bad json"}`, http.StatusBadRequest)
		return
	}
	t, err := h.Store.UpdateTransaction(r.Context(), uid, txID, repo.UpdateTransactionInput{
		CategoryID: req.CategoryID, Kind: req.Kind, Amount: req.Amount,
		Currency: req.Currency, Description: req.Description, OccurredAt: req.OccurredAt,
	})
	if err != nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}
	writeJSON(w, map[string]any{"data": t})
}

func (h *MeTransactions) Delete(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	txID, err := uuid.Parse(r.PathValue("txID"))
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}
	if err := h.Store.DeleteTransaction(r.Context(), uid, txID); err != nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type MeCategories struct {
	Store repo.CategoryStore
}

func (h *MeCategories) List(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	items, err := h.Store.ListCategoriesByUser(r.Context(), uid)
	if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"data": items})
}

func (h *MeCategories) Create(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	var req repo.CreateCategoryInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	c, err := h.Store.CreateCategory(r.Context(), uid, req)
	if err != nil {
		http.Error(w, `{"error":"could not create"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"data": c})
}

type MeGoals struct {
	Store  repo.GoalStore
	Alerts *alerts.Hub
}

func (h *MeGoals) List(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	items, err := h.Store.ListGoalsByUser(r.Context(), uid)
	if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	if h.Alerts != nil {
		CheckGoals(h.Alerts, items)
	}
	writeJSON(w, map[string]any{"data": items})
}

func (h *MeGoals) Create(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	var req repo.CreateGoalInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" || req.TargetAmount <= 0 {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if req.Currency == "" {
		req.Currency = "RUB"
	}
	if req.GoalType == "" {
		req.GoalType = "savings"
	}
	g, err := h.Store.CreateGoal(r.Context(), uid, req)
	if err != nil {
		http.Error(w, `{"error":"could not create"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"data": g})
}

func (h *MeGoals) Contribute(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	goalID, err := uuid.Parse(r.PathValue("goalID"))
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}
	var req struct {
		Amount float64 `json:"amount"`
		Note   string  `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Amount <= 0 {
		http.Error(w, `{"error":"invalid amount"}`, http.StatusBadRequest)
		return
	}
	g, err := h.Store.Contribute(r.Context(), uid, goalID, req.Amount, req.Note)
	if err != nil {
		http.Error(w, `{"error":"could not contribute"}`, http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]any{"data": g})
}
