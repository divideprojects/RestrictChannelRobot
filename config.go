package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)

// Environmental variables.
var (
	enableWebhook bool
	webhookPort   int
	botToken      string
	apiUrl        string
	webhookUrl    string
)

func init() {
	flag.StringVar(&botToken, "BOT_TOKEN", os.Getenv("BOT_TOKEN"), "Bot token for running the bot")
	flag.StringVar(&apiUrl, "API_URL", os.Getenv("API_URL"), "Api Server used to connect bot to")
	flag.BoolVar(&enableWebhook, "USE_WEBHOOKS", func() bool {
		return os.Getenv("USE_WEBHOOKS") == "yes" || os.Getenv("USE_WEBHOOKS") == "true"
	}(),
		"Enable webhooks",
	)
	flag.StringVar(&webhookUrl, "WEBHOOK_URL", os.Getenv("WEBHOOK_URL"), "URL for the Webhook")
	flag.IntVar(
		&webhookPort,
		"PORT",
		func(value string) int {
			if value == "" {
				return 0
			}
			val, err := strconv.Atoi(value)
			if err != nil {
				log.Fatal(err)
				return 0
			}
			return val
		}(os.Getenv("PORT")),
		"Port for the webhook",
	)

	if apiUrl == "" {
		apiUrl = "https://api.telegram.org"
	}
}
