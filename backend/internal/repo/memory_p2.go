package repo

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/fin-helper/backend/internal/domain"
	"github.com/google/uuid"
)

type MemoryPortfolioStore struct {
	mu       sync.RWMutex
	holdings map[uuid.UUID][]domain.Holding
}

func NewMemoryPortfolioStore() *MemoryPortfolioStore {
	uid := uuid.MustParse("00000000-0000-4000-8000-000000000001")
	return &MemoryPortfolioStore{
		holdings: map[uuid.UUID][]domain.Holding{
			uid: {{
				ID: uuid.New(), InstrumentID: uuid.New(), Symbol: "SBER", Exchange: "MOEX", Name: "Сбербанк",
				Quantity: 10, AvgCost: 250, Currency: "RUB", CostBasis: 2500,
			}},
		},
	}
}

func (m *MemoryPortfolioStore) GetPortfolio(_ context.Context, userID uuid.UUID) (*domain.PortfolioView, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	h := append([]domain.Holding{}, m.holdings[userID]...)
	var cost float64
	for _, x := range h {
		cost += x.CostBasis
	}
	return &domain.PortfolioView{
		ID: uuid.New(), Name: "Основной", Holdings: h, TotalCost: cost,
	}, nil
}

func (m *MemoryPortfolioStore) AddHolding(_ context.Context, userID uuid.UUID, symbol, exchange string, qty, avgCost float64) (*domain.Holding, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	h := domain.Holding{
		ID: uuid.New(), InstrumentID: uuid.New(), Symbol: symbol, Exchange: exchange, Name: symbol,
		Quantity: qty, AvgCost: avgCost, Currency: "RUB", CostBasis: qty * avgCost,
	}
	m.holdings[userID] = append(m.holdings[userID], h)
	return &h, nil
}

func (m *MemoryPortfolioStore) RemoveHolding(_ context.Context, userID, holdingID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	list := m.holdings[userID]
	for i, h := range list {
		if h.ID == holdingID {
			m.holdings[userID] = append(list[:i], list[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

type MemoryWatchlistStore struct {
	mu    sync.RWMutex
	items map[uuid.UUID][]domain.WatchlistItem
}

func NewMemoryWatchlistStore() *MemoryWatchlistStore {
	uid := uuid.MustParse("00000000-0000-4000-8000-000000000001")
	return &MemoryWatchlistStore{items: map[uuid.UUID][]domain.WatchlistItem{
		uid: {
			{InstrumentID: uuid.New(), Symbol: "SBER", Exchange: "MOEX", Name: "Сбербанк"},
			{InstrumentID: uuid.New(), Symbol: "GAZP", Exchange: "MOEX", Name: "Газпром"},
			{InstrumentID: uuid.New(), Symbol: "LKOH", Exchange: "MOEX", Name: "Лукойл"},
		},
	}}
}

func (m *MemoryWatchlistStore) ListWatchlist(_ context.Context, userID uuid.UUID) ([]domain.WatchlistItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]domain.WatchlistItem{}, m.items[userID]...), nil
}

func (m *MemoryWatchlistStore) AddToWatchlist(_ context.Context, userID uuid.UUID, symbol, exchange string) (*domain.WatchlistItem, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	w := domain.WatchlistItem{InstrumentID: uuid.New(), Symbol: symbol, Exchange: exchange, Name: symbol}
	m.items[userID] = append(m.items[userID], w)
	return &w, nil
}

func (m *MemoryWatchlistStore) RemoveFromWatchlist(_ context.Context, userID, instrumentID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	list := m.items[userID]
	for i, w := range list {
		if w.InstrumentID == instrumentID {
			m.items[userID] = append(list[:i], list[i+1:]...)
			return nil
		}
	}
	return nil
}

type MemorySettingsStore struct {
	mu   sync.RWMutex
	data map[uuid.UUID]domain.UserSettings
}

func NewMemorySettingsStore() *MemorySettingsStore {
	uid := uuid.MustParse("00000000-0000-4000-8000-000000000001")
	return &MemorySettingsStore{data: map[uuid.UUID]domain.UserSettings{
		uid: {UserID: uid, Locale: "ru", BaseCurrency: "RUB", Theme: "dark", Timezone: "Europe/Moscow"},
	}}
}

func (m *MemorySettingsStore) GetSettings(_ context.Context, userID uuid.UUID) (*domain.UserSettings, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if s, ok := m.data[userID]; ok {
		return &s, nil
	}
	s := domain.UserSettings{UserID: userID, Locale: "ru", BaseCurrency: "RUB", Theme: "dark", Timezone: "Europe/Moscow"}
	return &s, nil
}

func (m *MemorySettingsStore) UpdateSettings(_ context.Context, userID uuid.UUID, in domain.UserSettings) (*domain.UserSettings, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	in.UserID = userID
	m.data[userID] = in
	return &in, nil
}

type MemoryAnalyticsStore struct {
	Tx TransactionStore
}

func (m *MemoryAnalyticsStore) Analytics(ctx context.Context, userID uuid.UUID, from, to time.Time) (*domain.AnalyticsReport, error) {
	txs, err := m.Tx.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	report := &domain.AnalyticsReport{From: from.Format("2006-01-02"), To: to.Format("2006-01-02")}
	monthMap := map[string]*domain.AnalyticsBucket{}
	catMap := map[string]*domain.CategoryStat{}
	for _, t := range txs {
		if t.OccurredAt.Before(from) || !t.OccurredAt.Before(to) {
			continue
		}
		key := t.OccurredAt.Format("2006-01")
		if monthMap[key] == nil {
			monthMap[key] = &domain.AnalyticsBucket{Label: key}
		}
		if t.Kind == domain.KindIncome {
			monthMap[key].Income += t.Amount
		} else {
			monthMap[key].Expense += t.Amount
		}
		ck := t.Description + "|" + string(t.Kind)
		if catMap[ck] == nil {
			catMap[ck] = &domain.CategoryStat{Name: t.Description, Kind: string(t.Kind)}
		}
		catMap[ck].Total += t.Amount
	}
	for _, b := range monthMap {
		report.ByMonth = append(report.ByMonth, *b)
	}
	for _, c := range catMap {
		report.ByCategory = append(report.ByCategory, *c)
	}
	return report, nil
}

type MemorySummaryStore struct{}

func (MemorySummaryStore) EnsureInstrument(_ context.Context, _, _, _ string) (uuid.UUID, error) {
	return uuid.New(), nil
}

func (MemorySummaryStore) GetOrCreateSummary(_ context.Context, symbol, exchange string) (*domain.CompanySummary, error) {
	return &domain.CompanySummary{
		Symbol: symbol, Exchange: exchange,
		SummaryText: symbol + " — демо-сводка. Подключите PostgreSQL для кэширования обзоров.",
		KeyMetrics:  map[string]any{"demo": true},
	}, nil
}
