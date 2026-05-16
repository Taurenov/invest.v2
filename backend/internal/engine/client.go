package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const DisclaimerRU = "Прогноз сгенерирован моделью и не является инвестиционной рекомендацией. Решения о сделках принимайте самостоятельно или с лицензированным консультантом."

const DisclaimerEN = "This forecast is model-generated and is not investment advice. Make your own decisions or consult a licensed advisor."

type PredictResult struct {
	PredictedValue float64 `json:"predicted_value"`
	ChangePercent  float64 `json:"change_percent"`
	Confidence     float64 `json:"confidence"`
	ModelVersion   string  `json:"model_version"`
}

type Client struct {
	base   string
	client *http.Client
}

func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = "http://127.0.0.1:50052"
	}
	return &Client{
		base: baseURL,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) Predict(ctx context.Context, symbol string, prices []float64, horizonDays int) (*PredictResult, error) {
	body, _ := json.Marshal(map[string]any{
		"prices":       prices,
		"horizon_days": horizonDays,
		"symbol":       symbol,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+"/v1/predict", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("engine predict: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("engine status %d: %s", resp.StatusCode, b)
	}

	var out PredictResult
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ROI(ctx context.Context, initial, current float64) (float64, error) {
	body, _ := json.Marshal(map[string]float64{"initial": initial, "current": current})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+"/v1/roi", bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var out struct {
		RoiPercent float64 `json:"roi_percent"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return 0, err
	}
	return out.RoiPercent, nil
}

func (c *Client) CAGR(ctx context.Context, initial, finalValue, years float64) (float64, error) {
	var out struct {
		CagrPercent float64 `json:"cagr_percent"`
	}
	if err := c.postJSON(ctx, "/v1/cagr", map[string]float64{
		"initial": initial, "final": finalValue, "years": years,
	}, &out); err != nil {
		return 0, err
	}
	return out.CagrPercent, nil
}

func (c *Client) Savings(ctx context.Context, monthly, annualRatePct float64, months int, initialBalance float64) (float64, error) {
	var out struct {
		FutureValue float64 `json:"future_value"`
	}
	if err := c.postJSON(ctx, "/v1/savings", map[string]any{
		"monthly": monthly, "annual_rate_pct": annualRatePct,
		"months": months, "initial_balance": initialBalance,
	}, &out); err != nil {
		return 0, err
	}
	return out.FutureValue, nil
}

func (c *Client) postJSON(ctx context.Context, path string, body, out any) error {
	b, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+path, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("engine %s: %s", path, raw)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
