package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fin-helper/backend/internal/engine"
	"github.com/fin-helper/backend/internal/http/middleware"
	"github.com/fin-helper/backend/internal/market"
	"github.com/fin-helper/backend/internal/repo"
	"github.com/google/uuid"
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
	narrative := buildNarrative(locale, symbol, exchange, pred.ChangePercent, horizon)
	companyName := symbol
	sector := ""
	if profile, ok := market.GetCompanyProfile(symbol, exchange); ok {
		companyName = profile.Name
		sector = profile.Sector
	}
	var uidPtr *uuid.UUID
	if uid, ok := middleware.UserIDFromContext(r.Context()); ok {
		uidPtr = &uid
	}
	_ = h.DB.SaveForecast(
		r.Context(),
		inst.ID,
		uidPtr,
		horizon,
		pred.ChangePercent,
		pred.Confidence,
		narrative,
		pred.ModelVersion,
	)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"data": map[string]any{
			"symbol":               symbol,
			"exchange":             exchange,
			"company_name":         companyName,
			"sector":               sector,
			"horizon_days":         horizon,
			"predicted_value":      pred.PredictedValue,
			"predicted_change_pct": pred.ChangePercent,
			"confidence":           pred.Confidence,
			"model_version":        pred.ModelVersion,
			"narrative":            narrative,
			"disclaimer":           disclaimer,
		},
	})
}

func (h *ForecastHandler) History(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	exchange := r.URL.Query().Get("exchange")
	if exchange == "" {
		exchange = "MOEX"
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	inst, err := h.DB.GetInstrument(r.Context(), symbol, exchange)
	if err != nil {
		http.Error(w, `{"error":"instrument not found"}`, http.StatusNotFound)
		return
	}
	var uidPtr *uuid.UUID
	if uid, ok := middleware.UserIDFromContext(r.Context()); ok {
		uidPtr = &uid
	}
	items, err := h.DB.ListForecastHistory(r.Context(), inst.ID, uidPtr, limit)
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"data": items})
}

func (h *ForecastHandler) PriceHistory(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	exchange := r.URL.Query().Get("exchange")
	if exchange == "" {
		exchange = "MOEX"
	}
	points, _ := strconv.Atoi(r.URL.Query().Get("points"))
	inst, err := h.DB.GetInstrument(r.Context(), symbol, exchange)
	if err != nil {
		http.Error(w, `{"error":"instrument not found"}`, http.StatusNotFound)
		return
	}
	data, err := h.DB.ListPricePointsByInstrument(r.Context(), inst.ID, points)
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"data": map[string]any{
			"symbol":   symbol,
			"exchange": exchange,
			"points":   data,
		},
	})
}

func buildNarrative(locale, symbol, exchange string, change float64, horizon int) string {
	name := symbol
	sector := ""
	if profile, ok := market.GetCompanyProfile(symbol, exchange); ok {
		name = profile.Name
		sector = profile.Sector
	}
	dirRu := "рост"
	dirEn := "growth"
	if change < 0 {
		dirRu = "снижение"
		dirEn = "decline"
	}
	if locale == "en" {
		base := name + " (" + symbol + ")"
		if sector != "" {
			base += ", " + sector
		}
		return base + ": linear model suggests " + dirEn + " ~" + formatPct(abs(change)) + "% over " + strconv.Itoa(horizon) + " days based on recent MOEX prices."
	}
	base := name + " (" + symbol + ")"
	if sector != "" {
		base += ", сектор «" + sector + "»"
	}
	return base + ": по линейной модели ожидается " + dirRu + " ~" + formatPct(abs(change)) + "% за " + strconv.Itoa(horizon) + " дн. на основе истории MOEX."
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

func formatPct(v float64) string {
	return strconv.FormatFloat(v, 'f', 1, 64)
}
