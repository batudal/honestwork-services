package main

type Metadata struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Image       string        `json:"image"`
	ExternalUrl string        `json:"external_url"`
	Attributes  []interface{} `json:"attributes"`
}

type StarterMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	ExternalUrl string `json:"external_url"`
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
