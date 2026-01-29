package paypal_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	paypal_service "raznar.id/invoice-broker/internal/services/paypal"
)

type PayPalNewRequest struct {
	ID          string  `json:"id" binding:"required"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount" binding:"required"`
	ReturnURL   string  `json:"return_url" binding:"required"`
	CancelURL   string  `json:"cancel_url" binding:"required"`
}

func (p *PaypalController) CreateInvoice(c *gin.Context) {
	var req PayPalNewRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().
			Err(err).
			Str("ip", c.ClientIP()).
			Msg("failed to parse paypal create invoice request")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	approvalURL, err := p.Services.Paypal.CreateInvoice(paypal_service.CreateInvoicePayload{
		PaymentConfig: p.paymentConfig,
		ID:            req.ID,
		Description:   req.Description,
		Amount:        req.Amount,
		ReturnURL:     req.ReturnURL,
		CancelURL:     req.CancelURL,
	})

	if err != nil {
		log.Error().
			Err(err).
			Str("external_id", req.ID).
			Msg("failed to create paypal invoice")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create invoice",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"approval_url": approvalURL,
			"external_id":  req.ID,
		},
	})
}
