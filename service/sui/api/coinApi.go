package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"sort"
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

func RemainCoinAndGas(client *client.Client, address string, expectGas uint64) ([]CoinData, *sui_types.ObjectRef, error) {
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

	var gasCoin *CoinData
	var remainingCoins []CoinData
	var totalBalance uint64
	minDiff := ^uint64(0) // 最大的 uint64 值

	for _, coin := range coins.Data {
		balance, err := strconv.ParseUint(coin.Balance, 10, 64)
		if err != nil {
			continue
		}

		if balance >= expectGas {
			diff := balance - expectGas
			if diff < minDiff {
				if gasCoin != nil {
					remainingCoins = append(remainingCoins, *gasCoin)
					gasBalance, _ := strconv.ParseUint(coin.Balance, 10, 64)
					totalBalance += gasBalance
				}
				gasCoin = &coin
				minDiff = diff
				continue
			}
		}

		remainingCoins = append(remainingCoins, coin)
		totalBalance += balance
	}

	if gasCoin == nil {
		return nil, nil, errors.New("no suitable coin for gas found")
	}

	coinId, err := sui_types.NewObjectIdFromHex(gasCoin.CoinObjectId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse coin ID: %w", err)
	}

	version, err := strconv.ParseUint(gasCoin.Version, 10, 64)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse version: %w", err)
	}

	digest, err := sui_types.NewDigest(gasCoin.Digest)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse digest: %w", err)
	}

	gasObjectRef := &sui_types.ObjectRef{
		ObjectId: *coinId,
		Version:  version,
		Digest:   *digest,
	}

	return remainingCoins, gasObjectRef, nil
}

func MergeCoin(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, remainingCoins []CoinData, minMergeAmount uint64) ([]*sui_types.Argument, []CoinData, error) {
	if len(remainingCoins) == 0 {
		return nil, nil, errors.New("no coins to merge")
	}

	sort.Slice(remainingCoins, func(i, j int) bool {
		balanceI, _ := strconv.ParseUint(remainingCoins[i].Balance, 10, 64)
		balanceJ, _ := strconv.ParseUint(remainingCoins[j].Balance, 10, 64)
		return balanceI > balanceJ
	})

	var unusedCoins []CoinData
	for i, coin := range remainingCoins {
		balance, err := strconv.ParseUint(coin.Balance, 10, 64)
		if err != nil {
			continue
		}
		if balance >= minMergeAmount {
			coinArg, err := GetObjectArg(client, coin.CoinObjectId, true, "", "", "")
			if err != nil {
				continue
			}
			primaryCoin, err := ptb.Input(sui_types.CallArg{Object: coinArg})
			if err != nil {
				continue
			}

			unusedCoins = append(unusedCoins, remainingCoins[:i]...)
			unusedCoins = append(unusedCoins, remainingCoins[i+1:]...)

			return []*sui_types.Argument{&primaryCoin}, unusedCoins, nil
		}
	}

	var bestCombination []int
	var bestTotal uint64
	var currentTotal uint64

	for i := 0; i < len(remainingCoins); i++ {
		balance, err := strconv.ParseUint(remainingCoins[i].Balance, 10, 64)
		if err != nil {
			continue
		}
		currentTotal = balance
		currentCombination := []int{i}

		for j := i + 1; j < len(remainingCoins); j++ {
			nextBalance, err := strconv.ParseUint(remainingCoins[j].Balance, 10, 64)
			if err != nil {
				continue
			}

			if currentTotal+nextBalance >= minMergeAmount {
				currentCombination = append(currentCombination, j)
				currentTotal += nextBalance
				break
			}
			currentTotal += nextBalance
			currentCombination = append(currentCombination, j)
		}

		if currentTotal >= minMergeAmount && (len(bestCombination) == 0 || currentTotal < bestTotal) {
			bestCombination = currentCombination
			bestTotal = currentTotal
		}
	}

	if len(bestCombination) == 0 {
		return nil, remainingCoins, fmt.Errorf("cannot find suitable combination for amount %d", minMergeAmount)
	}

	// merge coin
	firstCoinArg, err := GetObjectArg(client, remainingCoins[bestCombination[0]].CoinObjectId, true, "", "", "")
	if err != nil {
		return nil, remainingCoins, fmt.Errorf("failed to get first coin object arg: %w", err)
	}

	primaryCoin, err := ptb.Input(sui_types.CallArg{Object: firstCoinArg})
	if err != nil {
		return nil, remainingCoins, fmt.Errorf("failed to input first coin: %w", err)
	}

	var coinsToMerge []sui_types.Argument
	for i := 1; i < len(bestCombination); i++ {
		coinToMerge, err := GetObjectArg(client, remainingCoins[bestCombination[i]].CoinObjectId, true, "", "", "")
		if err != nil {
			continue
		}
		coinArgument, err := ptb.Input(sui_types.CallArg{Object: coinToMerge})
		if err != nil {
			continue
		}
		coinsToMerge = append(coinsToMerge, coinArgument)
	}

	if len(coinsToMerge) > 0 {
		mergeResult := ptb.Command(
			sui_types.Command{
				MergeCoins: &struct {
					Argument  sui_types.Argument
					Arguments []sui_types.Argument
				}{
					Argument:  primaryCoin,
					Arguments: coinsToMerge,
				},
			},
		)

		if mergeResult.Result == nil {
			return nil, remainingCoins, errors.New("merge coins command should give a Result")
		}

		resultArg := &sui_types.Argument{
			NestedResult: &struct {
				Result1 uint16
				Result2 uint16
			}{Result1: *mergeResult.Result, Result2: 0},
		}

		usedIndexes := make(map[int]bool)
		for _, idx := range bestCombination {
			usedIndexes[idx] = true
		}

		for i := 0; i < len(remainingCoins); i++ {
			if !usedIndexes[i] {
				unusedCoins = append(unusedCoins, remainingCoins[i])
			}
		}

		return []*sui_types.Argument{resultArg}, unusedCoins, nil
	}

	return nil, remainingCoins, errors.New("failed to merge coins")
}