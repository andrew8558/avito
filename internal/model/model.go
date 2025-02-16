package model

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SendCoinRequest struct {
	ToUser string `json:"to_user"`
	Amount int32  `json:"amount"`
}

type InfoResponse struct {
	Coins       int32       `json:"coins"`
	Inventory   []Item      `json:"inventory"`
	CoinHistory CoinHistory `json:"coin_history"`
}

type Item struct {
	Type     string `json:"type" db:"type"`
	Quantity int32  `json:"quantity" db:"quantity"`
}

type CoinHistory struct {
	Received []ReceiveCoinEvent `json:"received"`
	Sent     []SendCoinEvent    `json:"sent"`
}

type SendCoinEvent struct {
	ToUser string `json:"to_user" db:"to_user"`
	Amount int32  `json:"amount" db:"amount"`
}

type ReceiveCoinEvent struct {
	FromUser string `json:"from_user" db:"from_user"`
	Amount   int32  `json:"amount" db:"amount"`
}
