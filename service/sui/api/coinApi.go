package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"strconv"
)


type CoinData struct {
	CoinType            string `json:"coinType"`
	CoinObjectId        string `json:"coinObjectId"`
	Version             string `json:"version"`
	Digest              string `json:"digest"`
	Balance             string `json:"balance"`
	PreviousTransaction string `json:"previousTransaction"`
}

type CoinPage struct {
	Data        []CoinData `json:"data"`
	NextCursor  string     `json:"nextCursor,omitempty"`
	HasNextPage bool       `json:"hasNextPage"`
}
func SplitCoin(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, amount uint64, address string, expectGas uint64) ([]*sui_types.Argument, *sui_types.ObjectRef, error) {
	coinType := "0x2::sui::SUI"
	sd, err := sui_types.NewAddressFromHex(address)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get SUI address: %w", err)
	}

	coinJson, err := client.GetCoins(context.Background(), *sd, &coinType, nil, 100)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get SUI coins: %w", err)
	}

	var coins CoinPage
	b, _ := json.Marshal(coinJson)
	err = json.Unmarshal(b, &coins)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse SUI coins JSON: %w", err)
	}

	if len(coins.Data) == 0 {
		return nil, nil, errors.New("account has no SUI coins")
	}

	// 计算需要的总金额
	totalNeeded := amount + expectGas

	// 检查第一个 coin 是否足够
	firstCoinBalance, err := strconv.ParseUint(coins.Data[0].Balance, 10, 64)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse coin balance: %w", err)
	}

	if firstCoinBalance < totalNeeded {
		return nil, nil, fmt.Errorf("insufficient balance: have %d, need %d", firstCoinBalance, totalNeeded)
	}

	// 保存第一个 coin 作为 gas payment
	coinId, err := sui_types.NewObjectIdFromHex(coins.Data[0].CoinObjectId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse coin ID: %w", err)
	}

	version, err := strconv.ParseUint(coins.Data[0].Version, 10, 64)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse version: %w", err)
	}

	digest, err := sui_types.NewDigest(coins.Data[0].Digest)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse digest: %w", err)
	}

	gasObjectRef := &sui_types.ObjectRef{
		ObjectId: *coinId,
		Version:  sui_types.SequenceNumber(version),
		Digest:   *digest,
	}

	// 如果有第二个 coin，使用它来分割
	if len(coins.Data) > 1 {
		secondCoinArg, err := GetObjectArg(client, coins.Data[1].CoinObjectId, true)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get second coin object arg: %w", err)
		}

		primaryCoin, err := ptb.Input(sui_types.CallArg{Object: secondCoinArg})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to input second coin: %w", err)
		}

		// 准备交易金额参数
		amt, err := ptb.Pure(amount)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to pure amount: %w", err)
		}

		// 从第二个 coin 中分割出交易金额
		splitCoinResult := ptb.Command(
			sui_types.Command{
				SplitCoins: &struct {
					Argument  sui_types.Argument
					Arguments []sui_types.Argument
				}{
					Argument:  primaryCoin,
					Arguments: []sui_types.Argument{amt},
				},
			},
		)

		if splitCoinResult.Result == nil {
			return nil, nil, errors.New("split coin command should always give a Result")
		}

		mainCoin := &sui_types.Argument{
			NestedResult: &struct {
				Result1 uint16
				Result2 uint16
			}{Result1: *splitCoinResult.Result, Result2: 0},
		}

		return []*sui_types.Argument{mainCoin}, gasObjectRef, nil
	}

	// 如果只有一个 coin，我们需要先创建一个新的 coin
	return nil, nil, errors.New("only one coin available, please get more coins first")
}