package rest

import (
	"context"
	"encoding/json"
	"time"

	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/xendit/xendit-go/v5"
	xenInvoice "github.com/xendit/xendit-go/v5/invoice"
	"raznar.id/invoice-broker/configs"
	"raznar.id/invoice-broker/pkg/internal/database"
	"raznar.id/invoice-broker/pkg/internal/database/models"
	"raznar.id/invoice-broker/pkg/internal/invoice"
)

type XenditNewRequest struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	// the format is in seconds
	InvoiceDuration int64 `json:"invoice_duration"`
}

func initXendit(app *fiber.App, conf configs.PaymentConfig, webConfig configs.WebConfig, db *database.Database) {
	xenClient := xendit.NewClient(conf.APIKey)
	app.Post(fmt.Sprintf("/gateway/%s/invoice", conf.Label), func(c *fiber.Ctx) (err error) {
		xenHeader := invoice.XenditHeader{
			CallbackToken: c.Get("x-callback-token"),
			WebhookID:     c.Get("webhook-id"),
		}

		code := invoice.IsCBValid(conf, xenHeader)
		if code != fiber.StatusOK {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		var invoiceData *xenInvoice.Invoice
		body := c.Body()
		err = c.BodyParser(&invoiceData)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		fmt.Println(invoiceData.GetId())
		model := db.Get(invoiceData.GetId())
		if model == nil {
			fmt.Println(model)
			return c.SendStatus(fiber.StatusNotFound)
		}

		go func() {
			invoice.ForwardWebhookData(body, webConfig.GetAPIConfig(model.Organization).CallbackURLS, "x-callback-token", xenHeader.CallbackToken)
			model.Status = invoiceData.GetStatus().String()
			db.Save()
		}()


		return c.SendStatus(code)
	})

	app.Post(fmt.Sprintf("/gateway/%s/invoice/new", conf.Label), func(c *fiber.Ctx) (err error) {
		apiTokenHeader := strings.Split(c.Get("Authorization"), "Bearer ")
		if len(apiTokenHeader) < 2 {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		apiToken := apiTokenHeader[1]
		apiData := getApiData(apiToken, c.IP(), webConfig.APIConfigs)
		if apiData == nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		var res XenditNewRequest
		err = c.BodyParser(&res)
		if err != nil {
			fmt.Printf("An error occured when parsing a request data from %s: %s", c.IP(), err.Error())
			return c.SendStatus(fiber.StatusBadGateway)
		}

		amount := res.Amount
		feeAmount := (res.Amount * 0.01)

		invDuration := fmt.Sprintf("%d", res.InvoiceDuration)
		req := xenClient.InvoiceApi.CreateInvoice(context.Background()).CreateInvoiceRequest(xenInvoice.CreateInvoiceRequest{
			ExternalId:  fmt.Sprintf("RAZNAR - B#%s", res.ID),
			Amount:      amount + feeAmount,
			Description: &res.Description,
			Fees: []xenInvoice.InvoiceFee{
				{
					Type:  "QR Fee",
					Value: float32(feeAmount),
				},
			},
			InvoiceDuration: &invDuration,
			PaymentMethods: []string{
				"QRIS",
			},
		})

		transaction, _, xenerr := req.Execute()
		if xenerr != nil {
			fmt.Printf("An error occured when creating an invoice link from %s: %s", c.IP(), xenerr.Error())
			return c.SendStatus(fiber.StatusBadGateway)
		}

		err = db.Add(&models.TransactionModel{
			Organization:  apiData.Organization,
			TransactionID: *transaction.Id,
			Amount:        amount,
			CreatedAt:     time.Now().Format(time.RFC3339),
			Gateway:       conf.Label,
			Status:        transaction.GetStatus().String(),
		})

		if err != nil {
			fmt.Printf("An error occured when saving invoice data from %s: %s", c.IP(), xenerr.Error())
			return c.SendStatus(fiber.StatusBadGateway)
		}

		return c.Status(fiber.StatusCreated).SendString(transaction.InvoiceUrl)
	})

	app.Get(fmt.Sprintf("/gateway/%s/invoice/data", conf.Label), func(c *fiber.Ctx) (err error) {
		apiTokenHeader := strings.Split(c.Get("Authorization"), "Bearer ")
		if len(apiTokenHeader) < 2 {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		apiToken := apiTokenHeader[1]
		apiData := getApiData(apiToken, c.IP(), webConfig.APIConfigs)
		if apiData == nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		transactionId := c.Query("id")
		model := db.Get(transactionId)
		if model == nil {
			return c.SendStatus(fiber.StatusNotFound)
		}

		res, err := json.Marshal(&model)
		if err != nil {
			return c.SendStatus(fiber.StatusBadGateway)
		}

		return c.Status(fiber.StatusOK).Send(res)
	})
}
