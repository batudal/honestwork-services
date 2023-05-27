package main

import (
	"context"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/getsentry/sentry-go"
	"github.com/honestworkdao/honestwork-services/rating/abi/hwescrow"
	"go.mongodb.org/mongo-driver/mongo"
)

type RatingWatcher struct {
	mc  *mongo.Client
	rpc *ethclient.Client
}

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

func NewRatingWatcher(mc *mongo.Client, rpc *ethclient.Client) *RatingWatcher {
	return &RatingWatcher{
		mc:  mc,
		rpc: rpc,
	}
}

func (r *RatingWatcher) WatchRatings() {
	for {
		r.fetchAllRatings()
		time.Sleep(time.Duration(30) * time.Minute)
	}
}

func (r *RatingWatcher) fetchAllRatings() {
	listers := fetchAllListers()
	members := fetchAllMembers()
	for _, lister := range listers {
		updateRating(lister)
	}
	for _, member := range members {
		updateRating(member)
	}
}

func search(length int, f func(index int) bool) int {
	for index := 0; index < length; index++ {
		if f(index) {
			return index
		}
	}
	return -1
}

func fetchAllListers() []string {
	return []string{}
}

func fetchAllMembers() []string {
	return []string{}
}

func updateRating(address string) {
	rating := FetchAggregatedRating(address)
	_ = rating
}

func FetchAggregatedRating(user_address string) float64 {
	var client *ethclient.Client
	client, err := ethclient.Dial(os.Getenv("ARBITRUM_RPC"))
	if err != nil {
		sentry.CaptureException(err)
	}
	defer client.Close()
	escrow_address_hex := common.HexToAddress(os.Getenv("ESCROW_ADDRESS"))
	instance, err := hwescrow.NewHwescrow(escrow_address_hex, client)
	if err != nil {
		sentry.CaptureException(err)
	}
	user_address_hex := common.HexToAddress(user_address)
	rating, err := instance.GetAggregatedRating(nil, user_address_hex)
	if err != nil {
		sentry.CaptureException(err)
	}
	precision, err := instance.GetPrecision(nil)
	if err != nil {
		sentry.CaptureException(err)
	}
	rating_normalized := float64(rating.Int64()) / float64(precision.Int64())
	return rating_normalized
}
