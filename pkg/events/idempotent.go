package events

import (
	"context"
	"database/sql"
	"sync"
)

// Seen records which event ids a consumer has already handled.
type Seen interface {
	// MarkSeen records id and reports whether it was new.
	MarkSeen(ctx context.Context, consumer, id string) (bool, error)
}

// Idempotent drops redeliveries of an event this consumer already handled.
func Idempotent(consumer string, seen Seen, h Handler) Handler {
	return func(ctx context.Context, e Envelope) error {
		fresh, err := seen.MarkSeen(ctx, consumer, e.ID)
		if err != nil {
			return err
		}
		if !fresh {
			return nil
		}
		return h(ctx, e)
	}
}

// MemorySeen is an in-process Seen for tests.
type MemorySeen struct {
	mu  sync.Mutex
	ids map[string]struct{}
}

// NewMemorySeen returns an empty MemorySeen.
func NewMemorySeen() *MemorySeen { return &MemorySeen{ids: map[string]struct{}{}} }

// MarkSeen implements Seen.
func (m *MemorySeen) MarkSeen(_ context.Context, consumer, id string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	k := consumer + "/" + id
	if _, ok := m.ids[k]; ok {
		return false, nil
	}
	m.ids[k] = struct{}{}
	return true, nil
}

// SQLSeen stores handled ids in the consumer's own database
// (table processed_messages, created by the service's migrations).
type SQLSeen struct{ DB *sql.DB }

// MarkSeen implements Seen with an insert that ignores duplicates.
func (s SQLSeen) MarkSeen(ctx context.Context, consumer, id string) (bool, error) {
	res, err := s.DB.ExecContext(ctx,
		`INSERT INTO processed_messages (consumer, message_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		consumer, id)
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	return n == 1, err
}
