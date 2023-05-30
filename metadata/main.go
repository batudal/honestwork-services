package main

import (
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/contrib/fibersentry"
	"github.com/gofiber/fiber/v2"
)

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		Debug:            true,
		AttachStacktrace: true,
	})
	if err != nil {
		panic(err)
	}
	app := fiber.New()
	app.Static("/metadata", "./genesis-metadata")
	app.Static("/metadata/starter", "./starter-metadata")
	ethereum_rpc, err := ethclient.Dial(os.Getenv("ETH_RPC"))
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	defer ethereum_rpc.Close()
	arbitrum_rpc, err := ethclient.Dial(os.Getenv("ARBITRUM_RPC"))
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	defer arbitrum_rpc.Close()
	go WriteStarterNFT()
	revenue_watcher := NewRevenueWatcher(ethereum_rpc, arbitrum_rpc)
	go revenue_watcher.WatchRevenues()
	app.Use(fibersentry.New(fibersentry.Config{
		Repanic:         true,
		WaitForDelivery: false,
	}))
	app.Listen(":3000")
}
