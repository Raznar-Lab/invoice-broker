package router_xendit_gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/xendit/xendit-go/v5"
	xenInvoice "github.com/xendit/xendit-go/v5/invoice"
	"raznar.id/invoice-broker/internal/models"
)

type xenditCreateInvoiceRequest struct {
	Description     string   `json:"description"`
	Amount          float64  `json:"amount"`
	InvoiceDuration int64    `json:"invoice_duration,omitempty"`
	CallbackURLS    []string `json:"callback_urls,omitempty"`
}

type xenditInvoiceCreateResponse struct {
	Id            string `json:"id"`
	TransactionID string `json:"transaction_id"`
	Link          string `json:"link"`
	ShortLink     string `json:"short_link"`
}

type xenditInvoiceGetResponse struct {
	Id            string  `json:"id"`
	TransactionID string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
	Gateway       string  `json:"gateway"`
	Status        string  `json:"status"`
	CreatedAt     string  `json:"created_at"`
}

func (r XenditGatewayRouter) InvoiceCreateHandler(c *fiber.Ctx) (err error) {
	authHeader := c.Get("Authorization")
	tokenParts := strings.Split(authHeader, "Bearer ")
	if len(tokenParts) < 2 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	paymentConfig := r.Config.Gateway.Xendit

	apiToken := tokenParts[1]
	apiData := r.Config.GetAPIConfig(apiToken, c.IP())
	if apiData == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var req xenditCreateInvoiceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	feeAmount := req.Amount * 0.01

	invId := uuid.NewString()[0:8]

	invReqBody := xenInvoice.CreateInvoiceRequest{
		ExternalId:  fmt.Sprintf("RAZNAR - B#%s", invId),
		Amount:      req.Amount + feeAmount,
		Description: &req.Description,
		Fees: []xenInvoice.InvoiceFee{
			{Type: "QR Fee", Value: float32(feeAmount)},
		},
		PaymentMethods: []string{"QRIS"},
	}

	if req.InvoiceDuration > 0 {
		newInvDur := fmt.Sprintf("%d", req.InvoiceDuration)
		invReqBody.InvoiceDuration = &newInvDur
	}

	for _, cbURL := range apiData.CallbackURLS {
		_, err := url.Parse(cbURL)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
	}

	invoiceRequest := xendit.NewClient(paymentConfig.APIKey).InvoiceApi.CreateInvoice(context.Background()).CreateInvoiceRequest(invReqBody)

	transaction, _, xenditError := invoiceRequest.Execute()
	if xenditError != nil {
		fmt.Printf("%s", xenditError.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	err = r.DB.AddTransaction(&models.TransactionModel{
		Id:            invId,
		Organization:  apiData.Organization,
		TransactionID: *transaction.Id,
		Amount:        req.Amount,
		CreatedAt:     time.Now().Format(time.RFC3339),
		Gateway:       paymentConfig.Label,
		Status:        transaction.GetStatus().String(),
		CallbackURLS:  req.CallbackURLS,
	})
	if err != nil {
		fmt.Printf("%s", err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	shortLinkID := uuid.NewString()[0:5]
	err = r.DB.AddShortener(&models.ShortenerModel{
		ID:   shortLinkID,
		Link: transaction.InvoiceUrl,
	})
	if err != nil {
		return c.SendStatus(fiber.StatusBadGateway)
	}

	responseData := xenditInvoiceCreateResponse{
		Id:            invId,
		TransactionID: transaction.GetId(),
		Link:          transaction.InvoiceUrl,
		ShortLink:     fmt.Sprintf("%s://%s/%s", c.Protocol(), c.Get("host"), shortLinkID),
	}

	responseJSON, err := json.Marshal(responseData)
	if err != nil {
		return c.SendStatus(fiber.StatusBadGateway)
	}

	return c.Status(fiber.StatusCreated).Send(responseJSON)
}

func (r XenditGatewayRouter) InvoiceGetHandler(c *fiber.Ctx) (err error) {
	authHeader := c.Get("Authorization")
	tokenParts := strings.Split(authHeader, "Bearer ")
	if len(tokenParts) < 2 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	apiToken := tokenParts[1]
	apiData := r.Config.GetAPIConfig(apiToken, c.IP())
	if apiData == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	invoiceId := c.Params("id")
	transaction := r.DB.GetTransaction(invoiceId)
	if transaction == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	if transaction.Organization != apiData.Organization {
		return c.SendStatus(fiber.StatusNotFound)
	}

	responseJSON, err := json.Marshal(xenditInvoiceGetResponse{
		Id:            transaction.Id,
		TransactionID: transaction.TransactionID,
		Amount:        transaction.Amount,
		Gateway:       transaction.Gateway,
		Status:        transaction.Status,
		CreatedAt:     transaction.CreatedAt,
	})
	if err != nil {
		return c.SendStatus(fiber.StatusBadGateway)
	}

	return c.Status(fiber.StatusOK).Send(responseJSON)
}
