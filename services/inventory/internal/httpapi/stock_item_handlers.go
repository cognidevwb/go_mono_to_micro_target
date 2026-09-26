// routing · HTTP handlers for StockItem on the monolith's own router, kept (gin) under /v1/stockitem: bind + validate the request → call the ported service → JSON response, errors as RFC 9457 problem+json. context.Context flows from the request into every call. Routes carry the /v1 segment. (inventory-service)
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

// problem is a minimal RFC 9457 (application/problem+json) error body.
type problem struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail,omitempty"`
}

func writeProblem(c *gin.Context, status int, title string, err error) {
	detail := ""
	if err != nil {
		detail = err.Error()
	}
	c.Header("Content-Type", "application/problem+json")
	c.AbortWithStatusJSON(status, problem{Type: "about:blank", Title: title, Status: status, Detail: detail})
}

// RegisterRoutes mounts the stock item endpoints under /v1/stockitem.
func RegisterRoutes(r *gin.Engine, db *gorm.DB, svc *inventory.Service) {
	g := r.Group("/v1/stockitem")
	g.POST("", func(c *gin.Context) { createStockItem(c, db) })
	g.GET("/:productId", func(c *gin.Context) { getStockItem(c, db) })
	g.POST("/reserve", func(c *gin.Context) { reserveStockItem(c, db, svc) })
}

func createStockItem(c *gin.Context, db *gorm.DB) {
	var legacy acl.LegacyStockItem
	if err := c.ShouldBindJSON(&legacy); err != nil {
		writeProblem(c, http.StatusBadRequest, "invalid request body", err)
		return
	}
	item, err := acl.TranslateStockItem(legacy)
	if err != nil {
		writeProblem(c, http.StatusBadRequest, "invalid stock item", err)
		return
	}
	if err := db.Create(&item).Error; err != nil {
		writeProblem(c, http.StatusUnprocessableEntity, "could not create stock item", err)
		return
	}
	c.JSON(http.StatusCreated, acl.ViewOf(item))
}

func getStockItem(c *gin.Context, db *gorm.DB) {
	productID, err := strconv.ParseUint(c.Param("productId"), 10, 64)
	if err != nil {
		writeProblem(c, http.StatusBadRequest, "invalid product id", err)
		return
	}
	var item inventory.StockItem
	if err := db.Where("product_id = ?", uint(productID)).First(&item).Error; err != nil {
		writeProblem(c, http.StatusNotFound, "stock item not found", err)
		return
	}
	c.JSON(http.StatusOK, acl.ViewOf(item))
}

type reserveRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

func reserveStockItem(c *gin.Context, db *gorm.DB, svc *inventory.Service) {
	var req reserveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeProblem(c, http.StatusBadRequest, "invalid request body", err)
		return
	}
	var reserved bool
	err := db.Transaction(func(tx *gorm.DB) error {
		ok, err := svc.Reserve(tx, req.ProductID, req.Quantity)
		reserved = ok
		return err
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeProblem(c, http.StatusNotFound, "stock item not found", err)
			return
		}
		writeProblem(c, http.StatusInternalServerError, "could not reserve stock", err)
		return
	}
	if !reserved {
		writeProblem(c, http.StatusConflict, "insufficient stock", errors.New("not enough units on hand"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"reserved": true})
}
