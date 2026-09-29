package app

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/acme/shop/services/payments/internal/acl"
)

// compensationLedger records which saga steps were already undone, so a
// repeated compensation is a no-op.
const compensationLedger = `CREATE TABLE IF NOT EXISTS payment_compensations (
    step_key   TEXT        PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
)`

// compensateRequest identifies the charge a CreateOrder saga made.
type compensateRequest struct {
	SagaID  string `json:"saga_id"`
	OrderID uint   `json:"order_id"`
}

func init() {
	schemas = append(schemas, compensationLedger)
	mounts = append(mounts, func(r *gin.Engine, d *Deps) {
		r.POST("/v1/payments/create_order/compensate", compensatePaymentHandler(d))
	})
}

// compensatePaymentHandler voids the order's captured payments in payments's own
// transaction. The ledger insert and the void commit together, so it runs once.
func compensatePaymentHandler(d *Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req compensateRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.SagaID == "" || req.OrderID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		if d == nil || d.DB == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "database unavailable"})
			return
		}
		key := req.SagaID + ":" + strconv.FormatUint(uint64(req.OrderID), 10)
		err := d.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
			res := tx.Exec(`INSERT INTO payment_compensations (step_key) VALUES (?) ON CONFLICT DO NOTHING`, key)
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return nil // already compensated
			}
			return tx.Exec(`UPDATE payments SET status = ? WHERE order_id = ? AND status = ?`,
				string(acl.PaymentVoided), req.OrderID, string(acl.PaymentCaptured)).Error
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "compensation failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"compensated": true})
	}
}
