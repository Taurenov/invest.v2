package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/fin-helper/backend/internal/market"
)

type QuotesHandler struct {
	Provider market.Provider
}

func (h *QuotesHandler) Get(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	exchange := r.URL.Query().Get("exchange")
	if exchange == "" {
		exchange = "MOEX"
	}

	q, err := h.Provider.Quote(r.Context(), symbol, exchange)
	if err != nil {
		http.Error(w, `{"error":"quote unavailable"}`, http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"data": q})
}
