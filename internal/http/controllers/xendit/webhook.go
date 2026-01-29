package xendit_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	xenInvoice "github.com/xendit/xendit-go/v5/invoice"
	xendit_service "raznar.id/invoice-broker/internal/services/xendit"
)

func (x *XenditController) ValidateWebhook(c *gin.Context) {
	var invoice xenInvoice.Invoice
	if err := c.ShouldBindJSON(&invoice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	cbToken := c.GetHeader("x-callback-token")

	payload := xendit_service.ValidationPayload{
		PaymentConfig: x.paymentConfig,
		Invoice:       &invoice,
		CallbackToken: cbToken,
	}

	if !x.Services.Xendit.ValidateWebhook(payload) {
		c.Status(http.StatusUnauthorized)
		return
	}

	c.Status(http.StatusOK)
}
