package market

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fin-helper/backend/internal/domain"
)

// MOEX ISS — публичный API Московской биржи (без ключа).
type MOEX struct {
	client *http.Client
}

func NewMOEX() *MOEX {
	return &MOEX{client: &http.Client{Timeout: 8 * time.Second}}
}

func (m *MOEX) Quote(ctx context.Context, symbol, exchange string) (*domain.Quote, error) {
	if exchange != "" && exchange != "MOEX" {
		return nil, fmt.Errorf("moex provider supports MOEX only")
	}
	url := fmt.Sprintf(
		"https://iss.moex.com/iss/engines/stock/markets/shares/securities/%s.json?iss.meta=off&iss.only=marketdata&marketdata.columns=LAST,CHANGE",
		symbol,
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var payload moexResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if len(payload.Marketdata.Data) == 0 {
		return nil, fmt.Errorf("no marketdata for %s", symbol)
	}
	row := payload.Marketdata.Data[0]
	if len(row) < 2 || row[0] == nil {
		return nil, fmt.Errorf("empty quote for %s", symbol)
	}

	price, _ := row[0].(float64)
	change, _ := row[1].(float64)
	var changePct float64
	if price > 0 && change != 0 {
		changePct = (change / (price - change)) * 100
	}

	return &domain.Quote{
		Symbol:    symbol,
		Exchange:  "MOEX",
		Price:     price,
		ChangePct: changePct,
	}, nil
}

type moexResponse struct {
	Marketdata struct {
		Data [][]any `json:"data"`
	} `json:"marketdata"`
}
