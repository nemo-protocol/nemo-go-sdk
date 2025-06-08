package models

type UID struct {
	Id [32]byte
}

type Coin struct {
	Id    UID
	Value uint64
}

type CoinInfo struct {
	CoinType  string
	Amount    float64
}