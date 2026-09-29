package catalog

import (
	"testing"
	"time"
)

func TestPriceOfServesUnexpiredCacheEntry(t *testing.T) {
	priceMu.Lock()
	priceCache[42] = priceCacheEntry{price: 9.5, expiresAt: time.Now().Add(time.Minute)}
	priceMu.Unlock()
	t.Cleanup(func() {
		priceMu.Lock()
		delete(priceCache, 42)
		priceMu.Unlock()
	})

	// A nil db is safe: a fresh entry must never reach the database.
	got, err := NewService(nil).PriceOf(42)
	if err != nil {
		t.Fatalf("PriceOf: %v", err)
	}
	if got != 9.5 {
		t.Fatalf("price = %v, want 9.5", got)
	}
}

func TestCreateRejectsNonPositivePrice(t *testing.T) {
	for _, price := range []float64{0, -1} {
		if err := NewService(nil).Create(&Product{SKU: "x", Name: "x", Price: price}); err == nil {
			t.Fatalf("price %v: expected error", price)
		}
	}
}
