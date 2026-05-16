package domain

import "github.com/google/uuid"

type Instrument struct {
	ID       uuid.UUID `json:"id"`
	Symbol   string    `json:"symbol"`
	Exchange string    `json:"exchange"`
	Name     string    `json:"name"`
}

type Quote struct {
	Symbol   string  `json:"symbol"`
	Exchange string  `json:"exchange"`
	Price    float64 `json:"price"`
	ChangePct float64 `json:"change_pct"`
}

type Forecast struct {
	Symbol              string  `json:"symbol"`
	HorizonDays         int     `json:"horizon_days"`
	PredictedChangePct  float64 `json:"predicted_change_pct"`
	PredictedValue      float64 `json:"predicted_value"`
	Confidence          float64 `json:"confidence"`
	Narrative           string  `json:"narrative"`
	Disclaimer          string  `json:"disclaimer"`
	ModelVersion        string  `json:"model_version"`
}
