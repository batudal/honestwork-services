package main

import (
	"context"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/getsentry/sentry-go"
)

func main() {
	_ = sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		Debug:            true,
		AttachStacktrace: true,
	})
	defer sentry.Flush(2 * time.Second)
	mongodb_client := NewMongoClient()
	defer mongodb_client.Disconnect(context.Background())
	rpc, err := ethclient.Dial(os.Getenv("ARBITRUM_RPC"))
	if err != nil {
		panic(err)
	}
	rw := NewRatingWatcher(mongodb_client, rpc)
	rw.WatchRatings()
}
