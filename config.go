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
	webhookDomain string
	databaseUrl   string
	databaseName  string
	webhookSecret string
)

func init() {
	flag.StringVar(&botToken, "BOT_TOKEN", os.Getenv("BOT_TOKEN"), "Bot token for running the bot")
	flag.BoolVar(&enableWebhook, "USE_WEBHOOKS", func() bool {
		return os.Getenv("USE_WEBHOOKS") == "yes" || os.Getenv("USE_WEBHOOKS") == "true"
	}(),
		"Enable webhooks",
	)
	flag.StringVar(&webhookSecret, "WEBHOOK_SECRET", os.Getenv("WEBHOOK_SECRET"), "Secret for webhook")
	flag.StringVar(&webhookDomain, "WEBHOOK_DOMAIN", os.Getenv("WEBHOOK_DOMAIN"), "URL for the Webhook")
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

	flag.StringVar(&databaseUrl, "DB_URI", os.Getenv("DB_URI"), "Database URI for MongoDB")
	flag.StringVar(&databaseName, "DB_NAME", os.Getenv("DB_NAME"), "Bot database name in MongoDB")
}
