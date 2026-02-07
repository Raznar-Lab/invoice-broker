package xendit_service

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/xendit/xendit-go/v5"
	xenInvoice "github.com/xendit/xendit-go/v5/invoice"
	"raznar.id/invoice-broker/configs"
)

type CreateInvoicePayload struct {
	PaymentConfig   *configs.PaymentConfig // Added PaymentConfig directly
	ID              string
	Description     string
	Amount          float64
	InvoiceDuration int64
}

func (x *XenditService) CreateInvoice(payload CreateInvoicePayload) (*xenInvoice.Invoice, error) {
	// 1. Initialize Logger with context
	logger := log.With().
		Str("external_id", payload.ID).
		Str("service", "xendit_create").
		Logger()

	// 2. Initialize Client using the passed config
	xenClient := xendit.NewClient(payload.PaymentConfig.APIKey)
	
	// Business Logic: 1% Fee
	feeAmount := payload.Amount * 0.01
	totalAmount := payload.Amount + feeAmount
	invDuration := fmt.Sprintf("%d", payload.InvoiceDuration)

	createReq := xenInvoice.CreateInvoiceRequest{
		ExternalId:  fmt.Sprintf("RAZNAR - B#%s", payload.ID),
		Amount:      totalAmount,
		Description: &payload.Description,
		Fees: []xenInvoice.InvoiceFee{
			{
				Type:  "QR Fee",
				Value: float32(feeAmount),
			},
		},
		InvoiceDuration: &invDuration,
		PaymentMethods: []string{"QRIS"},
	}

	// 3. Execute Request
	logger.Debug().Float64("total_amount", totalAmount).Msg("executing xendit invoice creation")
	
	resp, _, err := xenClient.InvoiceApi.
		CreateInvoice(context.Background()).
		CreateInvoiceRequest(createReq).
		Execute()

	if err != nil {
		logger.Error().Err(err).Msg("failed to create xendit invoice")
		return nil, err
	}

	logger.Info().
		Str("invoice_url", resp.InvoiceUrl).
		Msg("xendit invoice created successfully")

	return resp, nil
}