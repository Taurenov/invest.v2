package domain

import "github.com/google/uuid"

type UserSettings struct {
	UserID       uuid.UUID `json:"user_id"`
	Locale       string    `json:"locale"`
	BaseCurrency string    `json:"base_currency"`
	Theme        string    `json:"theme"`
	Timezone     string    `json:"timezone"`
}

type Alert struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Read    bool   `json:"read"`
}

type AnalyticsBucket struct {
	Label   string  `json:"label"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

type CategoryStat struct {
	Name    string  `json:"name"`
	Kind    string  `json:"kind"`
	Total   float64 `json:"total"`
}

type AnalyticsReport struct {
	From       string           `json:"from"`
	To         string           `json:"to"`
	ByMonth    []AnalyticsBucket `json:"by_month"`
	ByCategory []CategoryStat   `json:"by_category"`
}

type CompanySummary struct {
	Symbol      string         `json:"symbol"`
	Exchange    string         `json:"exchange"`
	SummaryText string         `json:"summary_text"`
	KeyMetrics  map[string]any `json:"key_metrics"`
}
