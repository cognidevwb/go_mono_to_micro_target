// Package httpapi serves inventory's StockItem over HTTP as RFC 9457
// problem+json on errors.
package httpapi

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/acme/shop/services/inventory/internal/acl"
	"github.com/acme/shop/services/inventory/internal/inventory"
)

// RegisterRoutes mounts the StockItem read route under /v1/stockitem.
func RegisterRoutes(r gin.IRouter, db *gorm.DB, svc *inventory.Service) {
	g := r.Group("/v1/stockitem")
	g.GET("/:product_id", func(c *gin.Context) { getStockItem(c, db) })
}

func problem(c *gin.Context, status int, detail string) {
	c.Header("Content-Type", "application/problem+json")
	c.AbortWithStatusJSON(status, gin.H{
		"type":   "about:blank",
		"title":  http.StatusText(status),
		"status": status,
		"detail": detail,
	})
}

func getStockItem(c *gin.Context, db *gorm.DB) {
	pid, err := strconv.ParseUint(c.Param("product_id"), 10, 64)
	if err != nil || pid == 0 {
		problem(c, http.StatusBadRequest, "product_id must be a positive integer")
		return
	}
	if db == nil {
		problem(c, http.StatusServiceUnavailable, "database unavailable")
		return
	}
	var item inventory.StockItem
	err = db.WithContext(c.Request.Context()).Where("product_id = ?", pid).First(&item).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		problem(c, http.StatusNotFound, "stock item not found")
	case err != nil:
		problem(c, http.StatusInternalServerError, "could not read stock item")
	default:
		c.JSON(http.StatusOK, acl.ToLegacy(item))
	}
}
