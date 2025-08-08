package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/move_types"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/nemo-protocol/nemo-go-sdk/service/sui/common/constant"
	"github.com/nemo-protocol/nemo-go-sdk/service/sui/common/models"
	"github.com/nemo-protocol/nemo-go-sdk/utils"
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

type PriceInfo struct {
	Logo    string `json:"logo"`
	Price   string `json:"price"`
	Decimal string `json:"decimal"`
	Name    string `json:"name"`
}

type PricePage struct {
	Data        map[string]PriceInfo `json:"data"`
	NextCursor  string               `json:"nextCursor,omitempty"`
	HasNextPage bool                 `json:"hasNextPage"`
}

func RemainCoinAndGas(client *client.Client, address string, expectGas uint64, coinType string) ([]CoinData, *sui_types.ObjectRef, error) {
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
		return nil, nil, errors.New(fmt.Sprintf("address:%v account has no %v coins", address, coinType))
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

		if balance >= expectGas && constant.IsGasCoinType(coinType) {
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

	if gasCoin == nil && constant.IsGasCoinType(coinType) {
		return nil, nil, errors.New("no suitable coin for gas found")
	}

	var gasObjectRef *sui_types.ObjectRef
	if gasCoin != nil {
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

		gasObjectRef = &sui_types.ObjectRef{
			ObjectId: *coinId,
			Version:  version,
			Digest:   *digest,
		}
	}

	return remainingCoins, gasObjectRef, nil
}

func MergeAllCoin(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, coinList []CoinData) (*sui_types.Argument, error) {
	if len(coinList) == 0{
		return nil, errors.New("coinList is null")
	}
	coinArg, err := GetObjectArg(client, coinList[0].CoinObjectId, true, "", "", "")
	if err != nil {
		return nil, err
	}
	primaryCoin, err := ptb.Input(sui_types.CallArg{Object: coinArg})
	if len(coinList) == 1{
		return &primaryCoin, nil
	}
	fmt.Printf("\n==primaryCoin:%v==\n",primaryCoin)
	var coinsToMerge []sui_types.Argument
	for i := 1; i < len(coinList); i++ {
		coinToMerge, err := GetObjectArg(client, coinList[i].CoinObjectId, true, "", "", "")
		if err != nil {
			continue
		}
		coinArgument, err := ptb.Input(sui_types.CallArg{Object: coinToMerge})
		if err != nil {
			continue
		}
		coinsToMerge = append(coinsToMerge, coinArgument)
	}
	ptb.Command(
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
	return &primaryCoin, nil
}

func MergeCoin(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, remainingCoins []CoinData, minMergeAmount uint64) ([]*sui_types.Argument, []CoinData, error) {
	if len(remainingCoins) == 0 {
		return nil, nil, errors.New("no coins to merge,please get one more coin exclude gas coin")
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
		ptb.Command(
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

		usedIndexes := make(map[int]bool)
		for _, idx := range bestCombination {
			usedIndexes[idx] = true
		}

		for i := 0; i < len(remainingCoins); i++ {
			if !usedIndexes[i] {
				unusedCoins = append(unusedCoins, remainingCoins[i])
			}
		}

		return []*sui_types.Argument{&primaryCoin}, unusedCoins, nil
	}

	return nil, remainingCoins, errors.New("failed to merge coins")
}

func SplitOrMergeCoin(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, remainingCoins []CoinData, netSyIn uint64) (splitCoin sui_types.Argument, unusedCoins []CoinData, err error) {
	sort.Slice(remainingCoins, func(i, j int) bool {
		balanceI, _ := strconv.ParseUint(remainingCoins[i].Balance, 10, 64)
		balanceJ, _ := strconv.ParseUint(remainingCoins[j].Balance, 10, 64)
		return balanceI > balanceJ
	})

	if len(remainingCoins) > 0 {
		balance, err := strconv.ParseUint(remainingCoins[0].Balance, 10, 64)
		if err != nil {
			return sui_types.Argument{}, nil, err
		}
		if balance >= netSyIn {
			coinArg, err := GetObjectArg(client, remainingCoins[0].CoinObjectId, true, "", "", "")
			if err != nil {
				return sui_types.Argument{}, nil, err
			}
			primaryCoin, err := ptb.Input(sui_types.CallArg{Object: coinArg})
			if err != nil {
				return sui_types.Argument{}, nil, err
			}
			splitResult, err := SplitCoinFromMerged(ptb, primaryCoin, netSyIn)
			if err != nil {
				return sui_types.Argument{}, nil, err
			}
			unusedCoins = append(unusedCoins, remainingCoins[1:]...)
			return splitResult, unusedCoins, nil
		}
	}

	mergedCoins, unusedCoins, err := MergeCoin(ptb, client, remainingCoins, netSyIn)
	if err != nil {
		return sui_types.Argument{}, remainingCoins, fmt.Errorf("failed to merge coins: %w", err)
	}

	if len(mergedCoins) == 0 {
		return sui_types.Argument{}, remainingCoins, errors.New("no coins merged")
	}

	splitResult, err := SplitCoinFromMerged(ptb, *mergedCoins[0], netSyIn)
	if err != nil {
		return sui_types.Argument{}, remainingCoins, fmt.Errorf("failed to split merged coin: %w", err)
	}

	return splitResult, unusedCoins, nil
}

func SwapToUnderlyingCoin(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, coinArgument *sui_types.Argument) (*sui_types.Argument, error) {
	if constant.IsSui(nemoConfig.UnderlyingCoinType) {
		return BurnSCoin(ptb, client, nemoConfig.CoinType, nemoConfig.UnderlyingCoinType, coinArgument)
	} else if constant.IsBuck(nemoConfig.UnderlyingCoinType) {
		return BurnToBuck(ptb, client, nemoConfig, coinArgument)
	}
	return nil, errors.New("invalid underlying coin！")
}

func GetCoinPriceInfo() map[string]PriceInfo {
	priceInfoUrl := "https://api.nemoprotocol.com/api/v1/market/info"
	priceInfoByte, err := utils.SendGetRpc(priceInfoUrl)
	if err != nil {
		return map[string]PriceInfo{}
	}
	pricePage := &PricePage{}
	_ = json.Unmarshal(priceInfoByte, pricePage)

	return pricePage.Data
}

func CoinIntoBalance(ptb *sui_types.ProgrammableTransactionBuilder, coinArgument *sui_types.Argument, coinType string) (*sui_types.Argument, error) {
	sui02Package, err := sui_types.NewObjectIdFromHex("0x0000000000000000000000000000000000000000000000000000000000000002")
	if err != nil {
		return nil, err
	}

	moduleName := "coin"
	functionName := "into_balance"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	coinTypeStructTag, err := GetStructTag(coinType)
	if err != nil {
		return nil, err
	}
	type1Tag := move_types.TypeTag{
		Struct: coinTypeStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, type1Tag)

	var arguments []sui_types.Argument

	arguments = append(arguments, *coinArgument)

	command := ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *sui02Package,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)
	return &command, nil
}

func CoinFromBalance(ptb *sui_types.ProgrammableTransactionBuilder, balanceArgument *sui_types.Argument, coinType string) (*sui_types.Argument, error) {
	sui02Package, err := sui_types.NewObjectIdFromHex("0x0000000000000000000000000000000000000000000000000000000000000002")
	if err != nil {
		return nil, err
	}

	moduleName := "coin"
	functionName := "from_balance"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	coinTypeStructTag, err := GetStructTag(coinType)
	if err != nil {
		return nil, err
	}
	type1Tag := move_types.TypeTag{
		Struct: coinTypeStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, type1Tag)

	var arguments []sui_types.Argument

	arguments = append(arguments, *balanceArgument)

	command := ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *sui02Package,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)
	return &command, nil
}