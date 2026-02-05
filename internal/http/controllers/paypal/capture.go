package paypal_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)
func (p *PaypalController) Return(c *gin.Context) {
	orderID := c.Query("token")

	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing order token",
		})
		return
	}

	if err := p.Services.Paypal.CaptureOrder(p.paymentConfig, orderID); err != nil {
		log.Error().
			Err(err).
			Str("order_id", orderID).
			Msg("failed to capture paypal order")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "payment capture failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"order_id": orderID,
			"status":   "paid",
		},
	})
}
