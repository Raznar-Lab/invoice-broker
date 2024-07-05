package rest

import (
	"context"
	"time"

	"fmt"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/xendit/xendit-go/v5"
	xenInvoice "github.com/xendit/xendit-go/v5/invoice"
	"raznar.id/invoice-broker/config"
	"raznar.id/invoice-broker/internal/invoice"
)

type XenditNewRequest struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	// the format is in seconds
	InvoiceDuration int64 `json:"invoice_duration"`
}

func initXendit(app *fiber.App, conf config.PaymentConfig, rdb *redis.Client) {
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

		go func() {
			if invoiceData != nil {
				cmd := rdb.Set(context.Background(), fmt.Sprintf("%s:%s", conf.Label, invoiceData.GetId()), invoiceData.GetStatus().String(), time.Until(time.Now().Add(time.Hour*24)))
				cmdErr := cmd.Err()
				if cmdErr != nil {
					fmt.Printf("An error occured while trying to update a transaction with id %s on redis: %s", invoiceData.GetId(), cmdErr.Error())
				}
			} else {
				fmt.Println("invoice data null")
			}

			invoice.ForwardWebhookData(body, conf.CallbackURLs, "x-callback-token", xenHeader.CallbackToken)
		}()
		return c.SendStatus(code)
	})

	app.Post(fmt.Sprintf("/gateway/%s/invoice/new", conf.Label), func(c *fiber.Ctx) (err error) {
		adminTokenHeader := strings.Split(c.Get("Authorization"), "Bearer ")
		if len(adminTokenHeader) < 2 {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		adminToken := adminTokenHeader[1]
		if !slices.Contains(conf.AdminToken, adminToken) || !slices.Contains(conf.AdminIP, c.IP()) {
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

		return c.Status(fiber.StatusCreated).SendString(transaction.InvoiceUrl)
	})
}
