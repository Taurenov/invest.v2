package ws

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/fin-helper/backend/internal/market"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Hub struct {
	mu       sync.Mutex
	clients  map[*websocket.Conn]struct{}
	provider market.Provider
	symbols  []struct{ Symbol, Exchange string }
}

func NewHub(provider market.Provider) *Hub {
	return &Hub{
		clients:  make(map[*websocket.Conn]struct{}),
		provider: provider,
		symbols: []struct{ Symbol, Exchange string }{
			{"SBER", "MOEX"}, {"GAZP", "MOEX"}, {"LKOH", "MOEX"},
		},
	}
}

func (h *Hub) SetSymbols(symbols []struct{ Symbol, Exchange string }) {
	if len(symbols) == 0 {
		return
	}
	h.mu.Lock()
	h.symbols = symbols
	h.mu.Unlock()
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	h.mu.Lock()
	h.clients[conn] = struct{}{}
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients, conn)
		h.mu.Unlock()
		_ = conn.Close()
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}

func (h *Hub) RunBroadcast(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			h.broadcastQuotes(ctx)
		}
	}
}

func (h *Hub) broadcastQuotes(ctx context.Context) {
	type msg struct {
		Type   string `json:"type"`
		Symbol string `json:"symbol"`
		Price  float64 `json:"price"`
		Change float64 `json:"change_pct"`
	}

	var batch []msg
	h.mu.Lock()
	syms := append([]struct{ Symbol, Exchange string }{}, h.symbols...)
	h.mu.Unlock()

	for _, s := range syms {
		q, err := h.provider.Quote(ctx, s.Symbol, s.Exchange)
		if err != nil {
			log.Printf("quote %s: %v", s.Symbol, err)
			continue
		}
		batch = append(batch, msg{Type: "quote", Symbol: q.Symbol, Price: q.Price, Change: q.ChangePct})
	}
	if len(batch) == 0 {
		return
	}

	payload, _ := json.Marshal(map[string]any{"type": "quotes", "data": batch})

	h.mu.Lock()
	defer h.mu.Unlock()
	for c := range h.clients {
		if err := c.WriteMessage(websocket.TextMessage, payload); err != nil {
			_ = c.Close()
			delete(h.clients, c)
		}
	}
}
