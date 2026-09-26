// Package saga is a small orchestration-saga runner: steps run in order, and
// when one fails every completed step is compensated in reverse. Put the step
// that is hardest to undo (the payment capture) LAST.
package saga

import (
	"context"
	"errors"
	"fmt"
)

// Step is one participant action and its compensation.
type Step struct {
	Name       string
	Do         func(ctx context.Context) error
	Compensate func(ctx context.Context) error
}

// Saga is a named ordered list of steps.
type Saga struct {
	Name  string
	Steps []Step
}

// Result names what ran and what was undone.
type Result struct {
	Completed   []string
	Compensated []string
}

// Run executes the steps. On failure it compensates the completed steps in
// reverse and returns the step error joined with any compensation errors.
func (s Saga) Run(ctx context.Context) (Result, error) {
	var res Result
	for i, st := range s.Steps {
		if err := st.Do(ctx); err != nil {
			stepErr := fmt.Errorf("saga %s: step %s: %w", s.Name, st.Name, err)
			var compErrs []error
			for j := i - 1; j >= 0; j-- {
				done := s.Steps[j]
				if done.Compensate == nil {
					continue
				}
				if cerr := done.Compensate(ctx); cerr != nil {
					compErrs = append(compErrs, fmt.Errorf("compensate %s: %w", done.Name, cerr))
					continue
				}
				res.Compensated = append(res.Compensated, done.Name)
			}
			return res, errors.Join(append([]error{stepErr}, compErrs...)...)
		}
		res.Completed = append(res.Completed, st.Name)
	}
	return res, nil
}
