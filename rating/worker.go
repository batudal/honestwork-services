package main

import (
	"context"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/getsentry/sentry-go"
	"github.com/honestworkdao/honestwork-services/rating/abi/hwescrow"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RatingWatcher struct {
	mc  *mongo.Client
	rpc *ethclient.Client
}

type IRatingWatcher interface {
	WatchRatings()
	FetchAllRatings()
	FetchAllListers() []string
	FetchAllMembers() []string
	UpdateRating(address string)
	FetchAggregatedRating(address string) float64
}

func NewRatingWatcher(mc *mongo.Client, rpc *ethclient.Client) *RatingWatcher {
	return &RatingWatcher{
		mc:  mc,
		rpc: rpc,
	}
}

func (r *RatingWatcher) WatchRatings() {
	for {
		r.FetchAllRatings()
		time.Sleep(time.Duration(1) * time.Hour)
	}
}

func (r *RatingWatcher) FetchAllRatings() {
	listers := r.FetchAllListers()
	members := r.FetchAllMembers()
	all := make([]string, len(listers)+len(members))
	copy(all, listers)
	copy(all[len(listers):], members[:])
	for _, user := range all {
		r.UpdateRating(user)
	}
}

func (r *RatingWatcher) FetchAllListers() []string {
	coll := r.mc.Database("honestwork-cluster").Collection("jobs")
	opts := options.Find().SetProjection(bson.D{{"useraddress", 1}})
	var listers []string
	cursor, err := coll.Find(context.Background(), bson.D{}, opts)
	if err != nil {
		return []string{}
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var job Job
		if err = cursor.Decode(&job); err != nil {
			sentry.CaptureException(err)
			continue
		}
		listers = append(listers, job.UserAddress)
	}
	return listers
}

func (r *RatingWatcher) FetchAllMembers() []string {
	coll := r.mc.Database("honestwork-cluster").Collection("users")
	opts := options.Find().SetProjection(bson.D{{"address", 1}})
	var members []string
	cursor, err := coll.Find(context.Background(), bson.D{}, opts)
	if err != nil {
		return []string{}
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var user User
		if err = cursor.Decode(&user); err != nil {
			sentry.CaptureException(err)
			continue
		}
		members = append(members, user.Address)
	}
	return members
}

func (r *RatingWatcher) UpdateRating(address string) {
	rating := r.FetchAggregatedRating(address)
	collection := r.mc.Database("honestwork-cluster").Collection("users")
	filter := bson.D{{"address", address}}
	update := bson.D{{"$set", bson.D{{"rating", rating}}}}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		sentry.CaptureException(err)
	}
}

func (r *RatingWatcher) FetchAggregatedRating(user_address string) float64 {
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
