package events

import (
	"context"
	"sync"
)

// Handler consumes one delivered event. Returning an error asks the broker to
// redeliver; delivery is at-least-once, so handlers must be idempotent — wrap
// them in Idempotent.
type Handler func(ctx context.Context, e Envelope) error

// Bus publishes and subscribes. Subjects are dotted names (`order.placed`).
type Bus interface {
	Publish(ctx context.Context, subject string, e Envelope) error
	Subscribe(ctx context.Context, subject, durable string, h Handler) error
	Close() error
}

// Memory is an in-process Bus for tests and for `event_backbone: none`.
type Memory struct {
	mu   sync.RWMutex
	subs map[string][]Handler
}

// NewMemory returns an empty in-process bus.
func NewMemory() *Memory { return &Memory{subs: map[string][]Handler{}} }

// Publish delivers e to every subscriber of subject, synchronously.
func (m *Memory) Publish(ctx context.Context, subject string, e Envelope) error {
	m.mu.RLock()
	hs := append([]Handler(nil), m.subs[subject]...)
	m.mu.RUnlock()
	for _, h := range hs {
		if err := h(ctx, e); err != nil {
			return err
		}
	}
	return nil
}

// Subscribe registers h for subject.
func (m *Memory) Subscribe(_ context.Context, subject, _ string, h Handler) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.subs[subject] = append(m.subs[subject], h)
	return nil
}

// Close is a no-op for the in-process bus.
func (m *Memory) Close() error { return nil }
