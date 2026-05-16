package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fin-helper/backend/internal/engine"
	"github.com/fin-helper/backend/internal/repo"
)

type ForecastHandler struct {
	DB     *repo.Postgres
	Engine *engine.Client
}

func (h *ForecastHandler) Predict(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	exchange := r.URL.Query().Get("exchange")
	if exchange == "" {
		exchange = "MOEX"
	}
	horizon, _ := strconv.Atoi(r.URL.Query().Get("horizon_days"))
	if horizon <= 0 {
		horizon = 7
	}

	locale := r.URL.Query().Get("locale")
	disclaimer := engine.DisclaimerRU
	if locale == "en" {
		disclaimer = engine.DisclaimerEN
	}

	inst, err := h.DB.GetInstrument(r.Context(), symbol, exchange)
	if err != nil {
		http.Error(w, `{"error":"instrument not found"}`, http.StatusNotFound)
		return
	}

	prices, err := h.DB.ListPrices(r.Context(), inst.ID, 120)
	if err != nil || len(prices) < 2 {
		http.Error(w, `{"error":"insufficient price history"}`, http.StatusBadRequest)
		return
	}

	pred, err := h.Engine.Predict(r.Context(), symbol, prices, horizon)
	if err != nil {
		http.Error(w, `{"error":"forecast engine unavailable"}`, http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"data": map[string]any{
			"symbol":               symbol,
			"exchange":             exchange,
			"horizon_days":         horizon,
			"predicted_value":      pred.PredictedValue,
			"predicted_change_pct": pred.ChangePercent,
			"confidence":           pred.Confidence,
			"model_version":        pred.ModelVersion,
			"narrative": buildNarrative(locale, symbol, pred.ChangePercent, horizon),
			"disclaimer":           disclaimer,
		},
	})
}

func buildNarrative(locale, symbol string, change float64, horizon int) string {
	if locale == "en" {
		if change >= 0 {
			return symbol + ": linear model suggests ~" + formatPct(change) + "% over " + strconv.Itoa(horizon) + " days."
		}
		return symbol + ": linear model suggests decline ~" + formatPct(-change) + "% over " + strconv.Itoa(horizon) + " days."
	}
	if change >= 0 {
		return symbol + ": линейная модель указывает на ~" + formatPct(change) + "% за " + strconv.Itoa(horizon) + " дн."
	}
	return symbol + ": линейная модель указывает на снижение ~" + formatPct(-change) + "% за " + strconv.Itoa(horizon) + " дн."
}

func formatPct(v float64) string {
	return strconv.FormatFloat(v, 'f', 1, 64)
}
