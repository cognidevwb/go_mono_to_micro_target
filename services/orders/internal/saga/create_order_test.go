package saga

import (
	"context"
	"testing"
)

func TestCreateOrderRunsThePivotLast(t *testing.T) {
	if n := len(CreateOrderOrder); n == 0 || CreateOrderOrder[n-1] != "payments" {
		t.Fatalf("order = %v", CreateOrderOrder)
	}
}

func TestCreateOrderFailsOnAnUnboundStep(t *testing.T) {
	if _, err := CreateOrder(nil).Run(context.Background()); err == nil {
		t.Fatal("an unbound saga must not succeed")
	}
}
