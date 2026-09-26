// Package events is the shared event contract every service speaks: one
// envelope, one bus interface, and an idempotent-consumer wrapper. The broker
// behind the interface is chosen once, in this package; services never import
// a broker client directly.
package events

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// Envelope is the wire shape of every event. Data carries the payload as raw
// JSON so a consumer can decode it into its own copy of the contract.
type Envelope struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	Source      string          `json:"source"`
	Time        time.Time       `json:"time"`
	TraceParent string          `json:"traceparent,omitempty"`
	Data        json.RawMessage `json:"data"`
}

// New builds an envelope with a fresh id for payload.
func New(eventType, source string, payload any) (Envelope, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return Envelope{}, fmt.Errorf("marshal %s payload: %w", eventType, err)
	}
	return Envelope{ID: NewID(), Type: eventType, Source: source, Time: time.Now().UTC(), Data: data}, nil
}

// Decode unmarshals the payload into v.
func (e Envelope) Decode(v any) error { return json.Unmarshal(e.Data, v) }

// NewID returns a random 128-bit hex id.
func NewID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b[:])
}
