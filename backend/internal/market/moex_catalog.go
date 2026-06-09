package market

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// CatalogItem — строка биржевого каталога MOEX (основной режим TQBR).
type CatalogItem struct {
	Symbol      string  `json:"symbol"`
	Exchange    string  `json:"exchange"`
	Name        string  `json:"name"`
	ShortName   string  `json:"short_name"`
	Sector      string  `json:"sector"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price"`
	ChangePct   float64 `json:"change_pct"`
}

// Catalog кэширует полный список акций MOEX и периодически обновляет котировки.
type Catalog struct {
	mu     sync.RWMutex
	items  []CatalogItem
	loaded time.Time
	ttl    time.Duration
	client *http.Client
}

var DefaultCatalog = &Catalog{
	ttl:    2 * time.Minute,
	client: &http.Client{Timeout: 20 * time.Second},
}

func (c *Catalog) Refresh(ctx context.Context) error {
	url := "https://iss.moex.com/iss/engines/stock/markets/shares/boards/TQBR/securities.json" +
		"?iss.meta=off&iss.only=securities,marketdata" +
		"&securities.columns=SECID,SHORTNAME,SECNAME" +
		"&marketdata.columns=SECID,LAST,CHANGE"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var raw struct {
		Securities struct {
			Data [][]any `json:"data"`
		} `json:"securities"`
		Marketdata struct {
			Data [][]any `json:"data"`
		} `json:"marketdata"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return err
	}

	quotes := make(map[string][2]float64)
	for _, row := range raw.Marketdata.Data {
		if len(row) < 3 {
			continue
		}
		sym, _ := row[0].(string)
		last, _ := toFloat(row[1])
		chg, _ := toFloat(row[2])
		var pct float64
		if last > 0 && chg != 0 {
			pct = (chg / (last - chg)) * 100
		}
		quotes[sym] = [2]float64{last, pct}
	}

	items := make([]CatalogItem, 0, len(raw.Securities.Data))
	for _, row := range raw.Securities.Data {
		if len(row) < 3 {
			continue
		}
		sym, _ := row[0].(string)
		short, _ := row[1].(string)
		name, _ := row[2].(string)
		if sym == "" {
			continue
		}
		item := CatalogItem{
			Symbol:    sym,
			Exchange:  "MOEX",
			Name:      name,
			ShortName: short,
			Sector:    "Прочее",
		}
		if profile, ok := GetCompanyProfile(sym, "MOEX"); ok {
			item.Sector = profile.Sector
			item.Description = profile.Description
			if item.Name == "" {
				item.Name = profile.Name
			}
		}
		if q, ok := quotes[sym]; ok {
			item.Price = q[0]
			item.ChangePct = q[1]
		}
		items = append(items, item)
	}

	if len(items) == 0 {
		return fmt.Errorf("moex catalog empty")
	}

	c.mu.Lock()
	c.items = items
	c.loaded = time.Now()
	c.mu.Unlock()
	return nil
}

func (c *Catalog) List(ctx context.Context, query string, limit int) ([]CatalogItem, error) {
	c.mu.RLock()
	stale := len(c.items) == 0 || time.Since(c.loaded) > c.ttl
	c.mu.RUnlock()

	if stale {
		if err := c.Refresh(ctx); err != nil {
			c.mu.RLock()
			empty := len(c.items) == 0
			c.mu.RUnlock()
			if empty {
				return fallbackCatalog(query, limit), nil
			}
		}
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	q := strings.ToUpper(strings.TrimSpace(query))
	out := make([]CatalogItem, 0, limit)
	for _, it := range c.items {
		if q != "" && !strings.Contains(it.Symbol, q) &&
			!strings.Contains(strings.ToUpper(it.Name), q) &&
			!strings.Contains(strings.ToUpper(it.ShortName), q) {
			continue
		}
		out = append(out, it)
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out, nil
}

func (c *Catalog) All(ctx context.Context) ([]CatalogItem, error) {
	return c.List(ctx, "", 0)
}

func fallbackCatalog(query string, limit int) []CatalogItem {
	items := ListMOEXProfiles()
	out := make([]CatalogItem, 0)
	q := strings.ToUpper(strings.TrimSpace(query))
	for _, p := range items {
		if q != "" && !strings.Contains(p.Symbol, q) && !strings.Contains(strings.ToUpper(p.Name), q) {
			continue
		}
		out = append(out, CatalogItem{
			Symbol: p.Symbol, Exchange: p.Exchange, Name: p.Name,
			Sector: p.Sector, Description: p.Description,
		})
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out
}

func toFloat(v any) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case json.Number:
		f, err := x.Float64()
		return f, err == nil
	default:
		return 0, false
	}
}
