// Package saga holds the CreateOrder orchestration saga. It replaces the
// monolith's single database transaction: each participant is now a remote
// call with a compensation, and the hardest step to undo runs last.
package saga

import (
	"context"

	pkgsaga "github.com/acme/shop/pkg/saga"
)

// CreateOrderOrder is the participant order: inventory → payments.
var CreateOrderOrder = []string{"inventory", "payments"}

// CreateOrder builds the saga from one step per participant, in CreateOrderOrder.
// Callers supply each participant's Do/Compensate bound to its typed client;
// a participant with no step bound fails immediately via notBound.
func CreateOrder(steps map[string]pkgsaga.Step) pkgsaga.Saga {
	s := pkgsaga.Saga{Name: "CreateOrder"}
	for _, p := range CreateOrderOrder {
		st, ok := steps[p]
		if !ok {
			st = pkgsaga.Step{Name: p, Do: notBound(p)}
		}
		s.Steps = append(s.Steps, st)
	}
	return s
}

func notBound(p string) func(context.Context) error {
	return func(context.Context) error { return &UnboundError{Participant: p} }
}

// UnboundError reports a participant with no step bound yet.
type UnboundError struct{ Participant string }

func (e *UnboundError) Error() string { return "saga step not bound: " + e.Participant }
