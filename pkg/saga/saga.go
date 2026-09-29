// Package saga is a small orchestration-saga runner: steps run in order, and
// when one fails every completed step is compensated in reverse. Put the step
// that is hardest to undo (the payment capture) LAST.
package saga

import (
	"context"
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

// StepError is the failure of one step. Its message is the step's own error,
// so a caller answers with exactly what the monolith's single transaction
// answered ("out of stock", not "saga CreateOrder: step inventory: out of
// stock"); Saga and Step say where it happened, for logs, and Compensation
// holds any compensation that failed in turn.
type StepError struct {
	Saga, Step   string
	Err          error
	Compensation []error
}

func (e *StepError) Error() string { return e.Err.Error() }

// Unwrap exposes the step's error and every failed compensation to errors.Is/As.
func (e *StepError) Unwrap() []error { return append([]error{e.Err}, e.Compensation...) }

// Run executes the steps. On failure it compensates the completed steps in
// reverse and returns a *StepError carrying the step's error and any
// compensation errors.
func (s Saga) Run(ctx context.Context) (Result, error) {
	var res Result
	for i, st := range s.Steps {
		if err := st.Do(ctx); err != nil {
			stepErr := &StepError{Saga: s.Name, Step: st.Name, Err: err}
			for j := i - 1; j >= 0; j-- {
				done := s.Steps[j]
				if done.Compensate == nil {
					continue
				}
				if cerr := done.Compensate(ctx); cerr != nil {
					stepErr.Compensation = append(stepErr.Compensation, fmt.Errorf("compensate %s: %w", done.Name, cerr))
					continue
				}
				res.Compensated = append(res.Compensated, done.Name)
			}
			return res, stepErr
		}
		res.Completed = append(res.Completed, st.Name)
	}
	return res, nil
}
