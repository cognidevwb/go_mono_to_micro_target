package app

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/acme/shop/services/catalog/internal/catalog"
)

// priceOfRequest is the argument list of catalog.Service.PriceOf.
type priceOfRequest struct {
	ProductID uint `json:"productId" binding:"required"`
}

// priceOfResponse is its result.
type priceOfResponse struct {
	Price float64 `json:"price"`
}

// Serves PriceOf to other services (orders); mounted beside wire's routes.
func init() {
	mounts = append(mounts, func(r *gin.Engine, d *Deps) {
		svc := catalog.NewService(d.DB)
		r.POST("/v1/catalog/price_of", func(c *gin.Context) {
			var req priceOfRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			price, err := svc.PriceOf(req.ProductID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, priceOfResponse{Price: price})
		})
	})
}
