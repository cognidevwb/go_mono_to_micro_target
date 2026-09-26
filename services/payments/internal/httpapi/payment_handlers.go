// Package httpapi hosts the HTTP handlers for Payment on the monolith's own
// router, kept (gin), under /v1/payment: bind + validate the request, call
// the ported service, respond as JSON. Errors are reported as RFC 9457
// problem+json. context.Context flows from the request into every call.
package httpapi

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/acme/shop/services/payments/internal/payments"
)

// chargeRequest is the wire shape POST /v1/payment/charge accepts.
type chargeRequest struct {
	OrderID uint    `json:"order_id" binding:"required"`
	Amount  float64 `json:"amount" binding:"required,gt=0"`
}

// problemDetail is an RFC 9457 problem+json response body.
type problemDetail struct {
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
	c.AbortWithStatusJSON(status, problemDetail{
		Type:   "about:blank",
		Title:  title,
		Status: status,
		Detail: detail,
	})
}

// RegisterRoutes mounts the payment endpoints under /v1/payment.
func RegisterRoutes(r gin.IRouter, db *gorm.DB, svc *payments.Service) {
	g := r.Group("/v1/payment")
	g.POST("/charge", func(c *gin.Context) { charge(c, db, svc) })
}

func charge(c *gin.Context, db *gorm.DB, svc *payments.Service) {
	var req chargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeProblem(c, http.StatusBadRequest, "invalid charge request", err)
		return
	}

	var payment payments.Payment
	err := db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := svc.Charge(tx, req.OrderID, req.Amount); err != nil {
			return err
		}
		return tx.Where("order_id = ?", req.OrderID).Order("id DESC").First(&payment).Error
	})
	if err != nil {
		slog.Error("charge failed", "order_id", req.OrderID, "err", err)
		writeProblem(c, http.StatusBadGateway, "payment declined", err)
		return
	}

	c.JSON(http.StatusCreated, payment)
}
