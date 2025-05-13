package models

type UID struct {
	Id [32]byte
}

type Coin struct {
	Id    UID
	Value uint64
}
