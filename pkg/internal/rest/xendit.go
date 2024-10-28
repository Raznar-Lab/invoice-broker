package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/xendit/xendit-go/v5"
	xenInvoice "github.com/xendit/xendit-go/v5/invoice"
	"raznar.id/invoice-broker/configs"
	"raznar.id/invoice-broker/pkg/internal/database"
	"raznar.id/invoice-broker/pkg/internal/database/models"
	"raznar.id/invoice-broker/pkg/internal/invoice"
)

type XenditInvoiceRequest struct {
	ID              string  `json:"id"`
	Description     string  `json:"description"`
	Amount          float64 `json:"amount"`
	InvoiceDuration int64   `json:"invoice_duration"`
}

type XenditInvoiceResponse struct {
	TransactionID string `json:"transaction_id"`
	Link          string `json:"link"`
	ShortLink     string `json:"short_link"`
}

func initXendit(app *fiber.App, conf configs.PaymentConfig, webConfig configs.WebConfig, db *database.Database) {
	xenditClient := xendit.NewClient(conf.APIKey)

	app.Post(fmt.Sprintf("/gateway/%s/invoice", conf.Label), func(c *fiber.Ctx) error {
		xenditHeader := invoice.XenditHeader{
			CallbackToken: c.Get("x-callback-token"),
			WebhookID:     c.Get("webhook-id"),
		}

		if !invoice.IsCBValid(conf, xenditHeader) {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		var invoiceData xenInvoice.Invoice
		if err := c.BodyParser(&invoiceData); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		transaction := db.GetTransaction(invoiceData.GetId())
		if transaction == nil {
			return c.SendStatus(fiber.StatusNotFound)
		}

		transaction.Status = invoiceData.GetStatus().String()
		db.SilentSave()

		body := c.Body()
		go func() {
			apiConf := webConfig.GetAPIConfig(transaction.Organization)
			if apiConf == nil {
				fmt.Printf("An error occured when forwarding webhook data, api conf with organization %s is null\n", transaction.Organization)
				return
			}

			invoice.ForwardWebhookData(body, apiConf.CallbackURLS, []string{"x-callback-token", "X-Callback-Token"}, apiConf.Token)
		}()

		return c.SendStatus(fiber.StatusOK)
	})

	app.Post(fmt.Sprintf("/gateway/%s/invoice/new", conf.Label), func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		tokenParts := strings.Split(authHeader, "Bearer ")
		if len(tokenParts) < 2 {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		apiToken := tokenParts[1]
		apiData := getApiData(apiToken, c.IP(), webConfig.APIConfigs)
		if apiData == nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		var req XenditInvoiceRequest
		if err := c.BodyParser(&req); err != nil {
			return c.SendStatus(fiber.StatusBadGateway)
		}

		feeAmount := req.Amount * 0.01
		invoiceDuration := fmt.Sprintf("%d", req.InvoiceDuration)

		invoiceRequest := xenditClient.InvoiceApi.CreateInvoice(context.Background()).CreateInvoiceRequest(xenInvoice.CreateInvoiceRequest{
			ExternalId:  fmt.Sprintf("RAZNAR - B#%s", req.ID),
			Amount:      req.Amount + feeAmount,
			Description: &req.Description,
			Fees: []xenInvoice.InvoiceFee{
				{Type: "QR Fee", Value: float32(feeAmount)},
			},
			InvoiceDuration: &invoiceDuration,
			PaymentMethods:  []string{"QRIS"},
		})

		transaction, _, xenditError := invoiceRequest.Execute()
		if xenditError != nil {
			return c.SendStatus(fiber.StatusBadGateway)
		}

		err := db.AddTransaction(&models.TransactionModel{
			Organization:  apiData.Organization,
			TransactionID: *transaction.Id,
			Amount:        req.Amount,
			CreatedAt:     time.Now().Format(time.RFC3339),
			Gateway:       conf.Label,
			Status:        transaction.GetStatus().String(),
		})
		if err != nil {
			return c.SendStatus(fiber.StatusBadGateway)
		}

		shortLinkID := uuid.NewString()[0:5]
		err = db.AddShortener(&models.ShortenerModel{
			ID:   shortLinkID,
			Link: transaction.InvoiceUrl,
		})
		if err != nil {
			return c.SendStatus(fiber.StatusBadGateway)
		}

		responseData := XenditInvoiceResponse{
			TransactionID: transaction.GetId(),
			Link:          transaction.InvoiceUrl,
			ShortLink:     fmt.Sprintf("https://%s/%s", c.Get("host"), shortLinkID),
		}

		responseJSON, err := json.Marshal(responseData)
		if err != nil {
			return c.SendStatus(fiber.StatusBadGateway)
		}

		return c.Status(fiber.StatusCreated).Send(responseJSON)
	})

	app.Get(fmt.Sprintf("/gateway/%s/invoice/data", conf.Label), func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		tokenParts := strings.Split(authHeader, "Bearer ")
		if len(tokenParts) < 2 {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		apiToken := tokenParts[1]
		apiData := getApiData(apiToken, c.IP(), webConfig.APIConfigs)
		if apiData == nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		transactionID := c.Query("id")
		transaction := db.GetTransaction(transactionID)
		if transaction == nil {
			return c.SendStatus(fiber.StatusNotFound)
		}

		responseJSON, err := json.Marshal(transaction)
		if err != nil {
			return c.SendStatus(fiber.StatusBadGateway)
		}

		return c.Status(fiber.StatusOK).Send(responseJSON)
	})
}
