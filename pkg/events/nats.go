package events

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// NATS is a JetStream-backed Bus. Every subject lives in one stream named
// after the first subject token (`order.placed` → stream ORDER); publishes
// carry Nats-Msg-Id = envelope id, so JetStream's duplicate window drops a
// retried publish.
type NATS struct {
	nc *nats.Conn
	js jetstream.JetStream
}

// ConnectNATS dials url and opens JetStream.
func ConnectNATS(url, name string) (*NATS, error) {
	nc, err := nats.Connect(url, nats.Name(name))
	if err != nil {
		return nil, fmt.Errorf("connect nats %s: %w", url, err)
	}
	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("jetstream: %w", err)
	}
	return &NATS{nc: nc, js: js}, nil
}

func streamFor(subject string) string {
	head, _, _ := strings.Cut(subject, ".")
	return strings.ToUpper(head)
}

func (n *NATS) ensure(ctx context.Context, subject string) (jetstream.Stream, error) {
	name := streamFor(subject)
	return n.js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     name,
		Subjects: []string{strings.ToLower(name) + ".>"},
	})
}

// Publish implements Bus.
func (n *NATS) Publish(ctx context.Context, subject string, e Envelope) error {
	if _, err := n.ensure(ctx, subject); err != nil {
		return err
	}
	body, err := json.Marshal(e)
	if err != nil {
		return err
	}
	_, err = n.js.Publish(ctx, subject, body, jetstream.WithMsgID(e.ID))
	return err
}

// DurableName makes name a legal JetStream consumer name. JetStream rejects
// `.`, `*`, `>`, whitespace and path separators in a durable, so a name built
// from a service and a subject (`inventory-service.order-placed`) would stop
// the service at startup; each such character becomes `-`.
func DurableName(name string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case '.', '*', '>', ' ', '\t', '\n', '/', '\\':
			return '-'
		}
		return r
	}, name)
}

// Subscribe implements Bus with a durable pull consumer.
func (n *NATS) Subscribe(ctx context.Context, subject, durable string, h Handler) error {
	st, err := n.ensure(ctx, subject)
	if err != nil {
		return err
	}
	cons, err := st.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       DurableName(durable),
		FilterSubject: subject,
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return err
	}
	_, err = cons.Consume(func(m jetstream.Msg) {
		var e Envelope
		if err := json.Unmarshal(m.Data(), &e); err != nil {
			_ = m.Term()
			return
		}
		if err := h(ctx, e); err != nil {
			_ = m.Nak()
			return
		}
		_ = m.Ack()
	})
	return err
}

// Close drains the connection.
func (n *NATS) Close() error { return n.nc.Drain() }
