package handlers

import (
	"net/http"
	"strconv"
	"time"
)

func parseDateRange(r *http.Request) (time.Time, time.Time) {
	now := time.Now()
	to := now
	from := now.AddDate(0, -11, 0)
	if v := r.URL.Query().Get("from"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			from = t
		}
	}
	if v := r.URL.Query().Get("to"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			to = t.Add(24 * time.Hour)
		}
	}
	return from, to
}

func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', 2, 64)
}
