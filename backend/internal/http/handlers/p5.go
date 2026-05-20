package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fin-helper/backend/internal/alerts"
	"github.com/fin-helper/backend/internal/domain"
	"github.com/fin-helper/backend/internal/http/middleware"
	"github.com/fin-helper/backend/internal/repo"
	"github.com/google/uuid"
)

type TagsHandler struct {
	Store repo.TagStore
}

func (h *TagsHandler) List(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	tags, err := h.Store.ListTags(r.Context(), uid)
	if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"data": tags})
}

func (h *TagsHandler) Create(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	var req struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	t, err := h.Store.CreateTag(r.Context(), uid, req.Name, req.Color)
	if err != nil {
		http.Error(w, `{"error":"could not create"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"data": t})
}

func (h *TagsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("tagID"))
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}
	_ = h.Store.DeleteTag(r.Context(), uid, id)
	w.WriteHeader(http.StatusNoContent)
}

func (h *TagsHandler) SetTxTags(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	txID, err := uuid.Parse(r.PathValue("txID"))
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}
	var req struct {
		TagIDs []uuid.UUID `json:"tag_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"bad json"}`, http.StatusBadRequest)
		return
	}
	if err := h.Store.SetTransactionTags(r.Context(), uid, txID, req.TagIDs); err != nil {
		http.Error(w, `{"error":"could not set tags"}`, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type BudgetsHandler struct {
	Store repo.BudgetStore
	Alerts *alerts.Hub
}

func (h *BudgetsHandler) List(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	items, err := h.Store.BudgetStatusThisMonth(r.Context(), uid)
	if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	if h.Alerts != nil {
		for _, s := range items {
			if s.Percent >= 100 {
				h.Alerts.Push(uid, "budget_exceeded", "Бюджет превышен", "Категория превысила лимит ("+strconv.FormatFloat(s.Percent, 'f', 0, 64)+"%)")
			}
		}
	}
	writeJSON(w, map[string]any{"data": items})
}

func (h *BudgetsHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	var req struct {
		CategoryID uuid.UUID `json:"category_id"`
		Amount     float64   `json:"amount"`
		Currency   string    `json:"currency"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Amount <= 0 {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	b, err := h.Store.UpsertBudget(r.Context(), uid, req.CategoryID, req.Amount, req.Currency)
	if err != nil {
		http.Error(w, `{"error":"could not save"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"data": b})
}

func (h *BudgetsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("budgetID"))
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}
	_ = h.Store.DeleteBudget(r.Context(), uid, id)
	w.WriteHeader(http.StatusNoContent)
}

type RecurringHandler struct {
	Store repo.RecurringStore
}

func (h *RecurringHandler) List(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	items, err := h.Store.ListRecurring(r.Context(), uid)
	if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"data": items})
}

func (h *RecurringHandler) Create(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	var req domain.RecurringTransaction
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Amount <= 0 || req.Kind == "" || req.Schedule == "" {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	req.UserID = uid
	if req.Currency == "" {
		req.Currency = "RUB"
	}
	created, err := h.Store.CreateRecurring(r.Context(), uid, req)
	if err != nil {
		http.Error(w, `{"error":"could not create"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"data": created})
}

func (h *RecurringHandler) Toggle(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("recurringID"))
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}
	active, _ := strconv.ParseBool(r.URL.Query().Get("active"))
	_ = h.Store.ToggleRecurring(r.Context(), uid, id, active)
	w.WriteHeader(http.StatusNoContent)
}

func (h *RecurringHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("recurringID"))
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}
	_ = h.Store.DeleteRecurring(r.Context(), uid, id)
	w.WriteHeader(http.StatusNoContent)
}

