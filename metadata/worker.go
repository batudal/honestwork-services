package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/getsentry/sentry-go"
	"github.com/honestworkdao/honestwork-services/metadata/abi/honestworknft"
	"github.com/honestworkdao/honestwork-services/metadata/abi/hwregistry"
)

// todo: multichain

type RevenueWatcher struct {
	arbitrum_rpc *ethclient.Client
	ethereum_rpc *ethclient.Client
}

type IRevenueWatcher interface {
	WatchRevenues()
}

func NewRevenueWatcher(ethereum_rpc *ethclient.Client, arbitrum_rpc *ethclient.Client) *RevenueWatcher {
	return &RevenueWatcher{
		arbitrum_rpc: arbitrum_rpc,
		ethereum_rpc: ethereum_rpc,
	}
}

func (r *RevenueWatcher) WatchRevenues() {
	for {
		r.fetchAllRevenues()
		time.Sleep(time.Duration(5) * time.Minute)
	}
}

func (r *RevenueWatcher) fetchAllRevenues() {
	total_supply, err := r.fetchTotalSupply()
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	for i := 1; i <= total_supply; i++ {
		revenue, err := r.fetchRevenue(i)
		if err != nil {
			sentry.CaptureException(err)
			continue
		}
		err = writeJSON(revenue)
		if err != nil {
			sentry.CaptureException(err)
			continue
		}
	}
}

func (r *RevenueWatcher) fetchRevenue(token_id int) (Revenue, error) {
	amount, err := r.fetchNFTRevenue(token_id)
	if err != nil {
		sentry.CaptureException(err)
		return Revenue{}, err
	}
	tier, err := r.fetchTokenTier(token_id)
	if err != nil {
		sentry.CaptureException(err)
		return Revenue{}, err
	}
	revenue_tier := getRevenueTier(amount)
	revenue := Revenue{
		TokenId:     token_id,
		Amount:      amount,
		Tier:        tier,
		RevenueTier: revenue_tier}
	return revenue, nil
}

func getRevenueTier(amount int) string {
	revenueTiers := []string{
		"< $1000",
		"$1000 - $10,000",
		"$10,000 - $100,000",
		"HonestChad",
	}
	if amount < 1000 {
		return revenueTiers[0]
	} else if amount < 10000 {
		return revenueTiers[1]
	} else if amount < 100000 {
		return revenueTiers[2]
	} else {
		return revenueTiers[3]
	}
}

func writeJSON(revenue Revenue) error {
	data := Metadata{
		Name:        "HonestWork #" + strconv.Itoa(revenue.TokenId),
		Description: "Introducing HonestWork Genesis NFT - the ultimate freelancer membership to our platform. AI-generated visuals and 3 tiers to choose from make your NFT a unique key to unlock access to HonestWork features and benefits. What's more, HonestWork Genesis also records your revenue on the blockchain, enabling you to earn future airdrops based on your performance. Join HonestWork today and take your freelancing career to the next level!",
		Image:       "https://honestwork-userfiles.fra1.cdn.digitaloceanspaces.com/genesis-nft/" + strconv.Itoa(revenue.TokenId) + ".png",
		ExternalUrl: "https://honestwork.app",
		Attributes: []interface{}{
			TierAttribute{
				TraitType: "Tier",
				Value:     revenue.Tier,
				MaxValue:  3,
			},
			GrossRevenueAttribute{
				TraitType: "Gross Revenue",
				Value:     revenue.Amount,
			},
			RevenueTierAttribute{
				TraitType: "Revenue Tier",
				Value:     revenue.RevenueTier,
			}}}
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		sentry.CaptureException(err)
		return err
	}
	err = ioutil.WriteFile(fmt.Sprintf("./genesis-metadata/%v", revenue.TokenId), file, 0644)
	if err != nil {
		sentry.CaptureException(err)
		return err
	}
	return nil
}

func (r *RevenueWatcher) fetchTotalSupply() (int, error) {
	nft_address_hex := common.HexToAddress(os.Getenv("GENESIS_ADDRESS"))
	instance, err := honestworknft.NewHonestworknft(nft_address_hex, r.ethereum_rpc)
	if err != nil {
		sentry.CaptureException(err)
		return 0, err
	}
	total_supply, err := instance.TotalSupply(nil)
	if err != nil {
		sentry.CaptureException(err)
		return 0, err
	}
	return int(total_supply.Int64()), nil
}

func (r *RevenueWatcher) fetchNFTRevenue(token_id int) (int, error) {
	registry_address_hex := common.HexToAddress(os.Getenv("REGISTRY_ADDRESS"))
	instance, err := hwregistry.NewHwregistry(registry_address_hex, r.arbitrum_rpc)
	if err != nil {
		sentry.CaptureException(err)
		return 0, err
	}
	revenue, err := instance.GetNFTGrossRevenue(nil, big.NewInt(int64(token_id)))
	if err != nil {
		sentry.CaptureException(err)
		return 0, err
	}
	revenue_normalized := new(big.Int)
	revenue_normalized.Div(revenue, big.NewInt(1000000000000000000))
	return int(revenue_normalized.Int64()), nil
}

func (r *RevenueWatcher) fetchTokenTier(token_id int) (int, error) {
	nft_address_hex := common.HexToAddress(os.Getenv("GENESIS_ADDRESS"))
	instance, err := honestworknft.NewHonestworknft(nft_address_hex, r.ethereum_rpc)
	if err != nil {
		sentry.CaptureException(err)
		return 0, err
	}
	token_id_big := big.NewInt(int64(token_id))
	state, err := instance.GetTokenTier(nil, token_id_big)
	if err != nil {
		sentry.CaptureException(err)
		return 0, err
	}
	return int(state.Int64()), nil
}
