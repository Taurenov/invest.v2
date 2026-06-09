package handlers

import (
	"net/http"
	"strconv"

	"github.com/fin-helper/backend/internal/market"
)

func ListMOEXInstruments(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 500
	}
	items, err := market.DefaultCatalog.List(r.Context(), q, limit)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "catalog unavailable")
		return
	}
	writeJSON(w, map[string]any{"data": items, "count": len(items)})
}

func MarketCatalog(w http.ResponseWriter, r *http.Request) {
	ListMOEXInstruments(w, r)
}

func RefreshMarketCatalog(w http.ResponseWriter, r *http.Request) {
	if err := market.DefaultCatalog.Refresh(r.Context()); err != nil {
		writeError(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	items, _ := market.DefaultCatalog.All(r.Context())
	writeJSON(w, map[string]any{"data": items, "count": len(items)})
}
