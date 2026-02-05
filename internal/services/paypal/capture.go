package paypal_service

import (
	"fmt"
	"net/http"
	"raznar.id/invoice-broker/configs"
)

func (p *PaypalService) CaptureOrder(cfg *configs.GatewayConfig, orderID string) error {
	token, err := p.getAccessToken(cfg)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", paypalEndpoint(cfg)+"/v2/checkout/orders/"+orderID+"/capture", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("paypal capture failed with status %s", resp.Status)
	}

	return nil
}
