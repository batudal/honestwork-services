package main

type Job struct {
	Email          string        `json:"email" validate:"required,email" bson:"email"`
	Slot           int           `json:"slot" bson:"slot"`
	UserAddress    string        `json:"useraddress" validate:"required,eth_addr" bson:"useraddress"`
	Username       string        `json:"username" validate:"required,min=5,max=50" bson:"username"`
	TokenPaid      string        `json:"tokenpaid" validate:"eth_addr" bson:"tokenpaid"`
	Title          string        `json:"title" validate:"required,min=5,max=50" bson:"title"`
	Description    string        `json:"description" validate:"required" bson:"description"` // custom
	Tags           []string      `json:"tags" validate:"required,min=1,max=3,dive,omitempty,min=2,max=20" bson:"tags"`
	Links          []string      `json:"links" validate:"required,min=1,max=3,dive,omitempty,url" bson:"links"`
	Budget         int           `json:"budget" validate:"required,min=200,max=100000" bson:"budget"`
	Timezone       *int64        `json:"timezone" validate:"required,min=-12,max=14" bson:"timezone"`
	TokensAccepted []Network     `json:"tokensaccepted" validate:"required,min=1" bson:"tokensaccepted"`
	StickyDuration int64         `json:"stickyduration" validate:"omitempty,lte=30" bson:"stickyduration"`
	CreatedAt      int64         `json:"createdat" bson:"createdat"`
	TxHash         string        `json:"txhash" validate:"required" bson:"txhash"`
	ImageUrl       string        `json:"imageurl" validate:"url" bson:"imageurl"`
	Applications   []Application `json:"application" bson:"application"`
	DealNetworkId  int           `json:"dealnetworkid" bson:"dealnetworkid"`
	DealId         int           `json:"dealid" bson:"dealid"`
}
type Application struct {
	UserAddress string `json:"useraddress" validate:"required" bson:"useraddress"`
	JobId       string `json:"jobid" validate:"required" bson:"jobid"`
	CoverLetter string `json:"coverletter" validate:"required" bson:"coverletter"`
	Date        int64  `json:"date" validate:"required" bson:"date"`
}

type Network struct {
	Id     int64   `json:"id" bson:"id"`
	Tokens []Token `json:"tokens" bson:"tokens"`
}
type Token struct {
	Address string `json:"address" bson:"address"`
}
