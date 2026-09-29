package app

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/acme/shop/services/customers/internal/customers"
)

// isActiveRequest is the argument list of customers.Service.IsActive.
type isActiveRequest struct {
	CustomerID uint `json:"customerId" binding:"required"`
}

// isActiveResponse is its result.
type isActiveResponse struct {
	Active bool `json:"active"`
}

// Serves IsActive to other services (orders); mounted beside wire's routes.
func init() {
	mounts = append(mounts, func(r *gin.Engine, d *Deps) {
		svc := customers.NewService(d.DB)
		r.POST("/v1/customers/is_active", func(c *gin.Context) {
			var req isActiveRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			active, err := svc.IsActive(req.CustomerID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, isActiveResponse{Active: active})
		})
	})
}
