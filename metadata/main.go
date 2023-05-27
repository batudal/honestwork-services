package main

import (
	"log"
	"os"

	"github.com/getsentry/sentry-go"
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
	log.Fatal(app.Listen(":3000"))
}
