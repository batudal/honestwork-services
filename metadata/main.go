package main

import (
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/contrib/fibersentry"
	"github.com/gofiber/fiber/v2"
)

func main() {
	_ = sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		Debug:            true,
		AttachStacktrace: true,
	})
	app := fiber.New()
	app.Static("/metadata", "./genesis-metadata")
	revenue_watcher := NewRevenueWatcher()
	go revenue_watcher.WatchRevenues()
	app.Use(fibersentry.New(fibersentry.Config{
		Repanic:         true,
		WaitForDelivery: false,
	}))
	app.Listen(":3000")
}
