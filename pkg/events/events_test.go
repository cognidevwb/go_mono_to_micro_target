package events

import (
	"context"
	"testing"
)

func TestIdempotentDropsARedelivery(t *testing.T) {
	bus := NewMemory()
	calls := 0
	h := Idempotent("test", NewMemorySeen(), func(context.Context, Envelope) error { calls++; return nil })
	if err := bus.Subscribe(context.Background(), "order.placed", "test", h); err != nil {
		t.Fatal(err)
	}
	e, err := New("order.placed", "orders", map[string]int{"id": 1})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if err := bus.Publish(context.Background(), "order.placed", e); err != nil {
			t.Fatal(err)
		}
	}
	if calls != 1 {
		t.Fatalf("handler ran %d times, want 1", calls)
	}
}

func TestEnvelopeRoundTrip(t *testing.T) {
	e, err := New("x", "src", map[string]string{"k": "v"})
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]string
	if err := e.Decode(&got); err != nil || got["k"] != "v" {
		t.Fatalf("decode: %v %v", got, err)
	}
}
