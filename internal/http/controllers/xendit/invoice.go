package xendit_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	xendit_service "raznar.id/invoice-broker/internal/services/xendit"
)



type XenditNewRequest struct {
	ID              string  `json:"id" binding:"required"`
	Description     string  `json:"description"`
	Amount          float64 `json:"amount" binding:"required"`
	InvoiceDuration int64   `json:"invoice_duration"`
}

func (x *XenditController) CreateInvoice(c *gin.Context) {
	var req XenditNewRequest

	// Gin uses ShouldBindJSON for body parsing
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Str("ip", c.ClientIP()).Msg("Failed to parse request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Call the service with the injected paymentConfig
	transaction, err := x.Services.Xendit.CreateInvoice(xendit_service.CreateInvoicePayload{
		PaymentConfig:   x.paymentConfig,
		ID:              req.ID,
		Description:     req.Description,
		Amount:          req.Amount,
		InvoiceDuration: req.InvoiceDuration,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"invoice_url": transaction.InvoiceUrl,
			"external_id": transaction.ExternalId,
		},
	})
}


