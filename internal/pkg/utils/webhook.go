package utils

import "strings"

func IsDiscordWebhook(u string) bool {
	return strings.HasPrefix(u, "https://discord.com/api/webhooks/")
}
