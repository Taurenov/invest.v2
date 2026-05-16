package alerts

import (
	"sync"

	"github.com/fin-helper/backend/internal/domain"
	"github.com/google/uuid"
)

type Hub struct {
	mu     sync.RWMutex
	byUser map[uuid.UUID][]domain.Alert
}

func NewHub() *Hub {
	return &Hub{byUser: make(map[uuid.UUID][]domain.Alert)}
}

func (h *Hub) Push(userID uuid.UUID, alertType, title, message string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	list := h.byUser[userID]
	alert := domain.Alert{
		ID:      uuid.New().String(),
		Type:    alertType,
		Title:   title,
		Message: message,
		Read:    false,
	}
	h.byUser[userID] = append([]domain.Alert{alert}, list...)
	if len(h.byUser[userID]) > 50 {
		h.byUser[userID] = h.byUser[userID][:50]
	}
}

func (h *Hub) List(userID uuid.UUID) []domain.Alert {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make([]domain.Alert, len(h.byUser[userID]))
	copy(out, h.byUser[userID])
	return out
}

func (h *Hub) MarkAllRead(userID uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for i := range h.byUser[userID] {
		h.byUser[userID][i].Read = true
	}
}
