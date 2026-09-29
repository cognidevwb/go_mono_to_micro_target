package app

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// compensationLedger records which saga steps were already undone, so a
// repeated compensation is a no-op.
const compensationLedger = `CREATE TABLE IF NOT EXISTS inventory_compensations (
    step_key   TEXT        PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
)`

// compensateRequest identifies the reservation a CreateOrder saga made.
type compensateRequest struct {
	SagaID    string `json:"saga_id"`
	ProductID uint   `json:"product_id"`
	Qty       int    `json:"qty"`
}

func init() {
	schemas = append(schemas, compensationLedger)
	mounts = append(mounts, func(r *gin.Engine, d *Deps) {
		r.POST("/v1/inventory/create_order/compensate", compensateHandler(d))
	})
}

// compensateHandler releases a reservation in inventory's own transaction. The
// ledger insert and the release commit together, so the release happens once.
func compensateHandler(d *Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req compensateRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.SagaID == "" || req.ProductID == 0 || req.Qty <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		if d == nil || d.DB == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "database unavailable"})
			return
		}
		key := req.SagaID + ":" + strconv.FormatUint(uint64(req.ProductID), 10)
		err := d.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
			res := tx.Exec(`INSERT INTO inventory_compensations (step_key) VALUES (?) ON CONFLICT DO NOTHING`, key)
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return nil // already compensated
			}
			return tx.Exec(`UPDATE stock_items SET reserved = reserved - ? WHERE product_id = ? AND reserved >= ?`,
				req.Qty, req.ProductID, req.Qty).Error
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "compensation failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"compensated": true})
	}
}
