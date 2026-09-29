package app

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/acme/shop/services/payments/internal/acl"
	"github.com/acme/shop/services/payments/internal/httpapi"
	"github.com/acme/shop/services/payments/internal/payments"
)

func init() {
	mounts = append(mounts, func(r *gin.Engine, d *Deps) {
		r.POST("/v1/payments/charge", chargeHandler(d))
		if d != nil && d.DB != nil {
			httpapi.RegisterRoutes(r, d.DB, payments.NewService(d.DB, nil))
		}
	})
}

// chargeHandler serves the ported Charge to orders: the provider is charged and
// the payment row written in one local transaction.
func chargeHandler(d *Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req acl.ChargeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		if err := req.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if d == nil || d.DB == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "database unavailable"})
			return
		}
		svc := payments.NewService(d.DB, payments.NewGateway(os.Getenv("PAYMENTS_URL")))
		err := d.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
			return svc.Charge(tx, req.OrderID, req.Amount)
		})
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"charged": true})
	}
}
