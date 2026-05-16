package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/fin-helper/backend/internal/alerts"
	"github.com/fin-helper/backend/internal/domain"
	"github.com/fin-helper/backend/internal/http/middleware"
	"github.com/fin-helper/backend/internal/market"
	"github.com/fin-helper/backend/internal/repo"
	"github.com/google/uuid"
)

type PortfolioHandler struct {
	Store  repo.PortfolioStore
	Quotes market.Provider
}

func (h *PortfolioHandler) Get(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	view, err := h.Store.GetPortfolio(r.Context(), uid)
	if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	var totalValue, totalPnL float64
	for i := range view.Holdings {
		q, err := h.Quotes.Quote(r.Context(), view.Holdings[i].Symbol, view.Holdings[i].Exchange)
		if err == nil {
			view.Holdings[i].CurrentPrice = q.Price
			view.Holdings[i].MarketValue = q.Price * view.Holdings[i].Quantity
			view.Holdings[i].PnL = view.Holdings[i].MarketValue - view.Holdings[i].CostBasis
			if view.Holdings[i].CostBasis > 0 {
				view.Holdings[i].PnLPercent = view.Holdings[i].PnL / view.Holdings[i].CostBasis * 100
			}
			totalValue += view.Holdings[i].MarketValue
			totalPnL += view.Holdings[i].PnL
		}
	}
	view.TotalValue = totalValue
	view.TotalPnL = totalPnL
	writeJSON(w, map[string]any{"data": view})
}

func (h *PortfolioHandler) Add(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	var req struct {
		Symbol   string  `json:"symbol"`
		Exchange string  `json:"exchange"`
		Quantity float64 `json:"quantity"`
		AvgCost  float64 `json:"avg_cost"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Quantity <= 0 || req.AvgCost <= 0 {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if req.Exchange == "" {
		req.Exchange = "MOEX"
	}
	holding, err := h.Store.AddHolding(r.Context(), uid, req.Symbol, req.Exchange, req.Quantity, req.AvgCost)
	if err != nil {
		http.Error(w, `{"error":"could not add"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"data": holding})
}

func (h *PortfolioHandler) Remove(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	hid, err := uuid.Parse(r.PathValue("holdingID"))
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}
	if err := h.Store.RemoveHolding(r.Context(), uid, hid); err != nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type WatchlistHandler struct {
	Store repo.WatchlistStore
}

func (h *WatchlistHandler) List(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	items, err := h.Store.ListWatchlist(r.Context(), uid)
	if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"data": items})
}

func (h *WatchlistHandler) Add(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	var req struct {
		Symbol   string `json:"symbol"`
		Exchange string `json:"exchange"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Symbol == "" {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if req.Exchange == "" {
		req.Exchange = "MOEX"
	}
	item, err := h.Store.AddToWatchlist(r.Context(), uid, req.Symbol, req.Exchange)
	if err != nil {
		http.Error(w, `{"error":"could not add"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"data": item})
}

func (h *WatchlistHandler) Remove(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	iid, err := uuid.Parse(r.PathValue("instrumentID"))
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}
	_ = h.Store.RemoveFromWatchlist(r.Context(), uid, iid)
	w.WriteHeader(http.StatusNoContent)
}

type SettingsHandler struct {
	Store repo.SettingsStore
}

func (h *SettingsHandler) Get(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	s, err := h.Store.GetSettings(r.Context(), uid)
	if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"data": s})
}

func (h *SettingsHandler) Patch(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	var req domain.UserSettings
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"bad json"}`, http.StatusBadRequest)
		return
	}
	s, err := h.Store.UpdateSettings(r.Context(), uid, req)
	if err != nil {
		http.Error(w, `{"error":"could not update"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"data": s})
}

type AnalyticsHandler struct {
	Store repo.AnalyticsStore
	Tx    repo.TransactionStore
}

func (h *AnalyticsHandler) Report(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	from, to := parseDateRange(r)
	report, err := h.Store.Analytics(r.Context(), uid, from, to)
	if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"data": report})
}

func (h *AnalyticsHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	from, to := parseDateRange(r)
	report, err := h.Store.Analytics(r.Context(), uid, from, to)
	if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=analytics.csv")
	_, _ = w.Write([]byte("month,income,expense\n"))
	for _, b := range report.ByMonth {
		_, _ = w.Write([]byte(b.Label + "," + formatFloat(b.Income) + "," + formatFloat(b.Expense) + "\n"))
	}
}

type SummaryHandler struct {
	Store repo.SummaryStore
}

func (h *SummaryHandler) Get(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	exchange := r.URL.Query().Get("exchange")
	if exchange == "" {
		exchange = "MOEX"
	}
	s, err := h.Store.GetOrCreateSummary(r.Context(), symbol, exchange)
	if err != nil {
		http.Error(w, `{"error":"summary unavailable"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"data": s})
}

type AlertsHandler struct {
	Hub *alerts.Hub
}

func (h *AlertsHandler) List(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	writeJSON(w, map[string]any{"data": h.Hub.List(uid)})
}

func (h *AlertsHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	h.Hub.MarkAllRead(uid)
	w.WriteHeader(http.StatusNoContent)
}

// CheckGoals notifies when goal completed
func CheckGoals(hub *alerts.Hub, goals []domain.Goal) {
	for _, g := range goals {
		if g.CurrentAmount >= g.TargetAmount {
			hub.Push(g.UserID, "goal_done", "Цель достигнута", g.Title)
		}
	}
}
