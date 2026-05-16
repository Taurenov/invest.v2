package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/fin-helper/backend/internal/engine"
)

type CalculatorHandler struct {
	Engine *engine.Client
}

func (h *CalculatorHandler) ROI(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Initial float64 `json:"initial"`
		Current float64 `json:"current"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"bad json"}`, http.StatusBadRequest)
		return
	}
	v, err := h.Engine.ROI(r.Context(), req.Initial, req.Current)
	if err != nil {
		http.Error(w, `{"error":"engine unavailable"}`, http.StatusServiceUnavailable)
		return
	}
	writeJSON(w, map[string]any{"data": map[string]float64{"roi_percent": v}})
}

func (h *CalculatorHandler) CAGR(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Initial float64 `json:"initial"`
		Final   float64 `json:"final"`
		Years   float64 `json:"years"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"bad json"}`, http.StatusBadRequest)
		return
	}
	v, err := h.Engine.CAGR(r.Context(), req.Initial, req.Final, req.Years)
	if err != nil {
		http.Error(w, `{"error":"engine unavailable"}`, http.StatusServiceUnavailable)
		return
	}
	writeJSON(w, map[string]any{"data": map[string]float64{"cagr_percent": v}})
}

func (h *CalculatorHandler) Savings(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Monthly        float64 `json:"monthly"`
		AnnualRatePct  float64 `json:"annual_rate_pct"`
		Months         int     `json:"months"`
		InitialBalance float64 `json:"initial_balance"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Months <= 0 {
		http.Error(w, `{"error":"bad json"}`, http.StatusBadRequest)
		return
	}
	v, err := h.Engine.Savings(r.Context(), req.Monthly, req.AnnualRatePct, req.Months, req.InitialBalance)
	if err != nil {
		http.Error(w, `{"error":"engine unavailable"}`, http.StatusServiceUnavailable)
		return
	}
	writeJSON(w, map[string]any{"data": map[string]float64{"future_value": v}})
}
