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
