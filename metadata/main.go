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
	app.Static("/metadata", "./metadata")
	revenue_watcher := NewRevenueWatcher()
	go revenue_watcher.WatchRevenues()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	log.Fatal(app.Listen(":3000"))
}
