// Package httpapi serves payments's read endpoints on the monolith's gin router.
package httpapi

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/acme/shop/services/payments/internal/payments"
)

// problem answers an RFC 9457 problem+json body.
func problem(c *gin.Context, status int, title string) {
	c.Header("Content-Type", "application/problem+json")
	c.AbortWithStatusJSON(status, gin.H{
		"type":   "about:blank",
		"title":  title,
		"status": status,
	})
}

// RegisterRoutes mounts the payment routes under /v1/payment.
func RegisterRoutes(r gin.IRouter, db *gorm.DB, svc *payments.Service) {
	g := r.Group("/v1/payment")
	g.GET("/order/:orderId", func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("orderId"), 10, 64)
		if err != nil || id == 0 {
			problem(c, http.StatusBadRequest, "invalid order id")
			return
		}
		var list []payments.Payment
		if err := db.WithContext(c.Request.Context()).Where("order_id = ?", id).Order("id").Find(&list).Error; err != nil {
			problem(c, http.StatusInternalServerError, "could not load payments")
			return
		}
		if len(list) == 0 {
			problem(c, http.StatusNotFound, "payment not found")
			return
		}
		c.JSON(http.StatusOK, list)
	})
}
