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

type User struct {
	Address      string        `json:"address" validate:"required,eth_addr" bson:"address"`
	Salt         string        `json:"salt" bson:"salt"`
	Username     string        `json:"username" validate:"required,min=5,max=50" bson:"username"`
	ShowEns      *bool         `json:"showens" validate:"required,boolean" bson:"showens"`
	EnsName      string        `json:"ensname" bson:"ensname"` // custom
	Title        string        `json:"title" validate:"required,min=5,max=50" bson:"title"`
	ImageUrl     string        `json:"imageurl" validate:"omitempty,url" bson:"imageurl"`
	ShowNFT      *bool         `json:"shownft" validate:"boolean" bson:"shownft"`
	NFTUrl       string        `json:"nfturl" validate:"omitempty,url" bson:"nfturl"`
	NFTAddress   string        `json:"nftaddress" validate:"omitempty,eth_addr" bson:"nftaddress"`
	NFTId        string        `json:"nftid" bson:"nftid"` // custom
	Email        string        `json:"email" validate:"omitempty,email" bson:"email"`
	Timezone     *int64        `json:"timezone" validate:"required,min=-12,max=14" bson:"timezone"`
	Bio          string        `json:"bio" bson:"bio"` // custom
	Links        []string      `json:"links" validate:"required,min=1,max=3,dive,omitempty,url" bson:"links"`
	Rating       int64         `json:"rating" bson:"rating"`
	Watchlist    []*Watchlist  `json:"watchlist" bson:"watchlist"`
	Favorites    []*Favorite   `json:"favorites" bson:"favorites"`
	DmsOpen      *bool         `json:"dmsopen" validate:"required,boolean" bson:"dmsopen"`
	Applications []Application `json:"application" bson:"application"`
}

type FavoriteInput struct {
	Address string `json:"address" bson:"address"`
	Slot    int    `json:"slot" bson:"slot"`
}

type Favorite struct {
	Input    *FavoriteInput `json:"input" bson:"input"`
	Username string         `json:"username" bson:"username"`
	Title    string         `json:"title" bson:"title"`
	ImageUrl string         `json:"imageurl" bson:"imageurl"`
}

type Watchlist struct {
	Input    *WatchlistInput `json:"input" bson:"input"`
	Username string          `json:"username" bson:"username"`
	Title    string          `json:"title" bson:"title"`
	ImageUrl string          `json:"imageurl" bson:"imageurl"`
}

type WatchlistInput struct {
	Address string `json:"address" bson:"address"`
	Slot    int    `json:"slot" bson:"slot"`
}
