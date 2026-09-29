package app

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/acme/shop/services/inventory/internal/inventory"
)

// reserveRequest is the argument list of inventory.Service.Reserve.
type reserveRequest struct {
	ProductID uint `json:"product_id"`
	Qty       int  `json:"qty"`
}

func init() {
	mounts = append(mounts, func(r *gin.Engine, d *Deps) {
		r.POST("/v1/inventory/reserve", reserveHandler(d))
	})
}

// reserveHandler serves the ported Reserve to orders. The stock row is locked
// for the transaction first, so the ported read-check-then-write cannot
// oversell under concurrent reservations.
func reserveHandler(d *Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req reserveRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.ProductID == 0 || req.Qty <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		if d == nil || d.DB == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "database unavailable"})
			return
		}
		svc := inventory.NewService(d.DB)
		var ok bool
		err := d.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
			var id uint
			if err := tx.Raw("SELECT id FROM stock_items WHERE product_id = ? FOR UPDATE", req.ProductID).
				Scan(&id).Error; err != nil {
				return err
			}
			if id == 0 {
				return gorm.ErrRecordNotFound
			}
			var err error
			ok, err = svc.Reserve(tx, req.ProductID, req.Qty)
			return err
		})
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "stock item not found"})
		case err != nil:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "reservation failed"})
		default:
			c.JSON(http.StatusOK, gin.H{"reserved": ok})
		}
	}
}
