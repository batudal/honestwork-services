package main

import (
	"context"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/getsentry/sentry-go"
)

// todo: add subscriber

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		Debug:            true,
		AttachStacktrace: true,
	})
	if err != nil {
		panic(err)
	}
	defer sentry.Flush(2 * time.Second)
	mongodb_client, err := NewMongoClient()
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	defer mongodb_client.Disconnect(context.Background())
	rpc, err := ethclient.Dial(os.Getenv("ARBITRUM_RPC"))
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	defer rpc.Close()
	deal_watcher := NewDealWatcher(mongodb_client, rpc)
	deal_watcher.WatchDeals()
}
