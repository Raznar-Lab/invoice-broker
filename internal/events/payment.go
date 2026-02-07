package events

import (
	"fmt"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/rs/zerolog/log"
	"raznar.id/invoice-broker/configs"
	"raznar.id/invoice-broker/internal/notifications"
	"raznar.id/invoice-broker/internal/notifications/webhook"
	"raznar.id/invoice-broker/internal/pkg/utils"
)

type PaymentWebhookEvent struct {
	ExternalID  string  `json:"external_id"`
	ID          string  `json:"id"`
	Status      string  `json:"status"`
	Fee         float64 `json:"fee"`
	Amount      float64 `json:"amount"`
	Gateway     string  `json:"gateway"`
	Currency    string  `json:"currency"`
	Description string  `json:"description"`

	PaymentConfig *configs.PaymentConfig `json:"-"`
	Raw           any                    `json:"-"`
}

func (p *PaymentWebhookEvent) emit() error {
	if err := p.sendCallbacks(); err != nil {
		return err
	}
	return p.sendWebhooks()
}

func (p *PaymentWebhookEvent) sendWebhooks() error {
	log.Debug().Strs("urls", p.PaymentConfig.Notification.Webhooks).Msg("Sending webhook notifications")
	var discordURLs, normalURLs []string

	for _, u := range p.PaymentConfig.Notification.Webhooks {
		if utils.IsDiscordWebhook(u) {
			discordURLs = append(discordURLs, u)
		} else {
			normalURLs = append(normalURLs, u)
		}
	}
	if len(discordURLs) > 0 {
		color := 0x3498db
		statusEmoji := "ℹ️"

		switch p.Status {
		case "success", "completed", "paid":
			color = 0x2ecc71
			statusEmoji = "✅"
		case "failed", "expired", "cancelled":
			color = 0xe74c3c
			statusEmoji = "❌"
		}

		embed := map[string]any{
			"title":       fmt.Sprintf("💳 %s", p.ExternalID),
			"description": fmt.Sprintf("**%s**", p.Description),
			"color":       color,
			"fields": []map[string]any{
				{
					"name":   "🆔 Transaction ID",
					"value":  fmt.Sprintf("`%s`", p.ID),
					"inline": true,
				},
				{
					"name":   "🚦 Status",
					"value":  fmt.Sprintf("%s %s", statusEmoji, strings.Title(p.Status)),
					"inline": true,
				},
				{
					"name":   "\u200b",
					"value":  "\u200b",
					"inline": true,
				},
				{
					"name":   "💰 Amount",
					"value":  fmt.Sprintf("**%s %s**", strings.ToUpper(p.Currency), humanize.Commaf(p.Amount)),
					"inline": true,
				},
				{
					"name":   "💰 Fee",
					"value":  fmt.Sprintf("**%s %s**", strings.ToUpper(p.Currency), humanize.Commaf(p.Fee)),
					"inline": true,
				},
				{
					"name":   "\u200b",
					"value":  "\u200b",
					"inline": true,
				},
				{
					"name":   "🏦 Gateway",
					"value":  strings.ToUpper(p.Gateway),
					"inline": false,
				},
			},
			"footer": map[string]any{
				"text": "Invoice Broker Service",
			},
			"timestamp": time.Now().Format(time.RFC3339),
		}
		payload := map[string]any{
			"embeds": []any{embed},
		}

		if err := notifications.Send(notifications.Payload{
			Content: payload,
			URLS:    discordURLs,
		}); err != nil {
			return err
		}
	}

	if len(normalURLs) == 0 {
		return nil
	}

	return notifications.Send(notifications.Payload{
		Content: p,
		URLS:    normalURLs,
	})

}

func (p *PaymentWebhookEvent) sendCallbacks() error {
	urls := p.PaymentConfig.CallbackURLs

	log.Debug().Strs("urls", urls).Msg("Sending callbacks")
	return webhook.Send(
		webhook.Payload{
			Content: p.Raw,
			URLS:    urls,
			Headers: map[string]string{
				p.PaymentConfig.WebhookHeader: p.PaymentConfig.WebhookToken,
			},
		},
	)

}
