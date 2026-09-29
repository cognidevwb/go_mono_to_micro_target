package saga

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestAFailedStepCompensatesTheCompletedOnesInReverse(t *testing.T) {
	ok := func(context.Context) error { return nil }
	s := Saga{Name: "t", Steps: []Step{
		{Name: "a", Do: ok, Compensate: ok},
		{Name: "b", Do: ok, Compensate: ok},
		{Name: "c", Do: func(context.Context) error { return errors.New("boom") }},
	}}
	res, err := s.Run(context.Background())
	if err == nil {
		t.Fatal("want an error")
	}
	if !reflect.DeepEqual(res.Compensated, []string{"b", "a"}) {
		t.Fatalf("compensated %v", res.Compensated)
	}
}

func TestAFailedStepAnswersWithItsOwnMessage(t *testing.T) {
	outOfStock := errors.New("out of stock")
	s := Saga{Name: "CreateOrder", Steps: []Step{{Name: "inventory", Do: func(context.Context) error { return outOfStock }}}}
	_, err := s.Run(context.Background())
	if err == nil || err.Error() != "out of stock" || !errors.Is(err, outOfStock) {
		t.Fatalf("err = %v, want the step's own message", err)
	}
	var se *StepError
	if !errors.As(err, &se) || se.Step != "inventory" {
		t.Fatalf("err = %#v, want a StepError naming the step", err)
	}
}
