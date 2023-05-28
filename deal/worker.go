package main

import (
	"context"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/getsentry/sentry-go"
	"github.com/honestworkdao/honestwork-services/deal/abi/hwescrow"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// todo: test

type DealWatcher struct {
	mc  *mongo.Client
	rpc *ethclient.Client
}

func NewDealWatcher(mc *mongo.Client, rpc *ethclient.Client) *DealWatcher {
	return &DealWatcher{
		mc:  mc,
		rpc: rpc,
	}
}

func (r *DealWatcher) WatchDeals() {
	for {
		updateDeals(r.mc, r.rpc)
		time.Sleep(time.Duration(1) * time.Hour)
	}
}

func updateDeals(mc *mongo.Client, rpc *ethclient.Client) {
	escrow_address_hex := common.HexToAddress(os.Getenv("ESCROW_ADDRESS"))
	instance, err := hwescrow.NewHwescrow(escrow_address_hex, rpc)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	deals, err := instance.GetDeals(nil)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	if len(deals) == 0 {
		return
	}
	coll := mc.Database("honestwork-cluster").Collection("deals")
	for i, deal := range deals {
		job := Job{}
		err = coll.FindOne(context.TODO(), bson.D{{"useraddress", deal.Recruiter.String()}, {"slot", int(deal.JobId.Int64())}}).Decode(&job)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				continue
			}
			sentry.CaptureException(err)
		}
		if job.DealId != -1 {
			continue
		}
		_, err := coll.UpdateOne(context.TODO(), bson.D{
			{"useraddress", deal.Recruiter.String()},
			{"slot", int(deal.JobId.Int64())},
		}, bson.D{
			{"$set", bson.D{
				{"dealid", i},
				{"dealnetworkid", 42161},
			}},
		})
		if err != nil {
			sentry.CaptureException(err)
		}
	}
}
