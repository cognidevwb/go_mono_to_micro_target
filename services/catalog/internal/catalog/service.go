package catalog

import (
	"errors"
	"sync"
	"time"

	"gorm.io/gorm"
)

// priceCacheTTL bounds how long a replica may serve a cached price.
const priceCacheTTL = 30 * time.Second

// priceCacheEntry is a cached price with its expiry.
type priceCacheEntry struct {
	price     float64
	expiresAt time.Time
}

// priceCache is per-replica: catalog is its only owner, and each entry expires
// after priceCacheTTL so replicas converge.
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
	return s.db.Create(p).Error
}

// PriceOf returns a product's price, cached in-process.
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
