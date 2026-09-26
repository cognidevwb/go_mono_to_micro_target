package catalog

import (
	"errors"
	"sync"
	"time"

	"gorm.io/gorm"
)

// priceCacheTTL bounds how long a replica may serve a stale price before
// re-reading catalog's own database; the cache is invalidated on write too.
const priceCacheTTL = 30 * time.Second

type priceCacheEntry struct {
	price     float64
	expiresAt time.Time
}

// priceCache is per-replica, TTL-bound and invalidated on write. catalog is
// the owner of this data, so a short-lived cache here (not a copy held by
// another service) is an accepted tradeoff rather than shared mutable state.
var (
	priceMu    sync.RWMutex
	priceCache = map[uint]priceCacheEntry{}
)

// Service owns the catalog context.
type Service struct{ db *gorm.DB }

// NewService builds the catalog service.
func NewService(db *gorm.DB) *Service { return &Service{db: db} }

// List returns every product.
func (s *Service) List() ([]Product, error) {
	var products []Product
	err := s.db.Find(&products).Error
	return products, err
}

// Create stores a new product.
func (s *Service) Create(p *Product) error {
	if p.Price <= 0 {
		return errors.New("price must be positive")
	}
	if err := s.db.Create(p).Error; err != nil {
		return err
	}
	invalidatePrice(p.ID)
	return nil
}

// PriceOf returns a product's price, cached in-process with a TTL.
func (s *Service) PriceOf(productID uint) (float64, error) {
	priceMu.RLock()
	entry, ok := priceCache[productID]
	priceMu.RUnlock()
	if ok && time.Now().Before(entry.expiresAt) {
		return entry.price, nil
	}
	var p Product
	if err := s.db.First(&p, productID).Error; err != nil {
		return 0, err
	}
	priceMu.Lock()
	priceCache[productID] = priceCacheEntry{price: p.Price, expiresAt: time.Now().Add(priceCacheTTL)}
	priceMu.Unlock()
	return p.Price, nil
}

// invalidatePrice drops productID from the cache so the next read is fresh.
func invalidatePrice(productID uint) {
	priceMu.Lock()
	delete(priceCache, productID)
	priceMu.Unlock()
}
