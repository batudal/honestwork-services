package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

type Metadata struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Image       string        `json:"image"`
	ExternalUrl string        `json:"external_url"`
	Attributes  []interface{} `json:"attributes"`
}

type TierAttribute struct {
	TraitType string `json:"trait_type"`
	Value     int    `json:"value"`
	MaxValue  int    `json:"max_value"`
}

type GrossRevenueAttribute struct {
	TraitType string `json:"trait_type"`
	Value     int    `json:"value"`
}

type RevenueTierAttribute struct {
	TraitType string `json:"trait_type"`
	Value     string `json:"value"`
}

type Revenue struct {
	TokenId     int    `json:"token_id"`
	Amount      int    `json:"amount"`
	Tier        int    `json:"tier"`
	RevenueTier string `json:"revenue_tier"`
}

type RevenueWatcher struct {
}

func NewRevenueWatcher() *RevenueWatcher {
	return &RevenueWatcher{}
}

func (r *RevenueWatcher) WatchRevenues() {
	for {
		fetchAllRevenues()
		log.Println("Fetched revenues and refreshed metadata")
		time.Sleep(time.Duration(5) * time.Minute)
	}
}

func fetchAllRevenues() {
	total_supply := FetchTotalSupply()
	for i := 1; i <= total_supply; i++ {
		revenue := fetchRevenue(i)
		writeJSON(revenue)
	}
}

func fetchRevenue(token_id int) Revenue {
	amount := FetchNFTRevenue(token_id)
	tier := FetchTokenTier(token_id)
	revenue_tier := getRevenueTier(amount)
	revenue := Revenue{
		TokenId:     token_id,
		Amount:      amount,
		Tier:        tier,
		RevenueTier: revenue_tier}
	return revenue
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

func writeJSON(revenue Revenue) {
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
	file, _ := json.MarshalIndent(data, "", " ")
	_ = ioutil.WriteFile(fmt.Sprintf("./genesis-metadata/%v", revenue.TokenId), file, 0644)
}

func FetchTotalSupply() int {
	client, err := ethclient.Dial(os.Getenv("ETH_RPC"))
	if err != nil {
		sentry.CaptureException(err)
	}
	nft_address_hex := common.HexToAddress(os.Getenv("GENESIS_ADDRESS"))
	instance, err := honestworknft.NewHonestworknft(nft_address_hex, client)
	if err != nil {
		sentry.CaptureException(err)
	}
	total_supply, err := instance.TotalSupply(nil)
	if err != nil {
		sentry.CaptureException(err)
	}
	client.Close()
	return int(total_supply.Int64())
}

func FetchNFTRevenue(token_id int) int {
	client, err := ethclient.Dial(os.Getenv("ARBITRUM_RPC"))
	if err != nil {
		sentry.CaptureException(err)
	}
	defer client.Close()
	registry_address_hex := common.HexToAddress(os.Getenv("REGISTRY_ADDRESS"))
	instance, err := hwregistry.NewHwregistry(registry_address_hex, client)
	if err != nil {
		sentry.CaptureException(err)
	}
	revenue, err := instance.GetNFTGrossRevenue(nil, big.NewInt(int64(token_id)))
	if err != nil {
		sentry.CaptureException(err)
	}
	revenue_normalized := new(big.Int)
	revenue_normalized.Div(revenue, big.NewInt(1000000000000000000))
	return int(revenue_normalized.Int64())
}

func FetchTokenTier(token_id int) int {
	client, err := ethclient.Dial(os.Getenv("ETH_RPC"))
	if err != nil {
		sentry.CaptureException(err)
	}
	nft_address_hex := common.HexToAddress(os.Getenv("GENESIS_ADDRESS"))
	instance, err := honestworknft.NewHonestworknft(nft_address_hex, client)
	if err != nil {
		sentry.CaptureException(err)
	}
	token_id_big := big.NewInt(int64(token_id))
	state, err := instance.GetTokenTier(nil, token_id_big)
	if err != nil {
		sentry.CaptureException(err)
	}
	client.Close()
	return int(state.Int64())
}
