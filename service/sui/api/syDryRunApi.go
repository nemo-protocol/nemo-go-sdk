package api

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coming-chat/go-sui/v2/account"
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/move_types"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/fardream/go-bcs/bcs"
	"github.com/shopspring/decimal"
	"math"
	"github.com/nemo-protocol/nemo-go-sdk/service/sui/common/constant"
	"github.com/nemo-protocol/nemo-go-sdk/service/sui/common/models"
	"strconv"
	"time"
)

type PoolRewarderInfo struct {
	TotalReward       string `json:"total_reward"`
	EndTime           string `json:"end_time"`
	LastRewardTime    string `json:"last_reward_time"`
	RewardHarvested   string `json:"reward_harvested"`
	RewardDebt        string `json:"reward_debt"`
	RewardToken       RewardToken `json:"reward_token"`
	AccPerShare       string `json:"acc_per_share"`
	Active            bool   `json:"active"`
	EmissionPerSecond string `json:"emission_per_second"`
	ID                ID     `json:"id"`
	Owner             string `json:"owner"`
	StartTime         string `json:"start_time"`
}

type RewardToken struct {
	Type   string `json:"type"`
	Fields struct {
		Name string `json:"name"`
	} `json:"fields"`
}

type ID struct {
	ID string `json:"id"`
}

type RawMarketState struct {
	TotalSy    string `json:"total_sy"`
	TotalPt    string `json:"total_pt"`
	LpSupply   string `json:"lp_supply"`
	MarketCap  string `json:"market_cap"`
	RewardPool struct {
		Fields struct {
			Rewarders struct {
				Fields struct {
					Contents []struct {
						Fields struct {
							Key   string `json:"key"`
							Value struct {
								Fields PoolRewarderInfo `json:"fields"`
							} `json:"value"`
						} `json:"fields"`
					} `json:"contents"`
				} `json:"fields"`
			} `json:"rewarders"`
		} `json:"fields"`
	} `json:"reward_pool"`
}

func DryRunGetApproxPyOutForNetSyInInternal(client *client.Client, nemoConfig *models.NemoConfig, exactPyType string, netSyIn, minYtOut uint64, sender *account.Account) (approxPyOut uint64, netSyTokenization uint64, err error) {
	ptb := sui_types.NewProgrammableTransactionBuilder()

	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.NemoContract)
	if err != nil {
		return 0, 0, err
	}

	syStructTag, err := GetStructTag(nemoConfig.SyCoinType)
	if err != nil {
		return 0, 0, err
	}

	moduleName := "offchain"
	var functionName string
	if exactPyType == constant.PTTYPE{
		functionName = "get_approx_pt_out_for_net_sy_in_internal"
	}else if exactPyType == constant.YTTYPE{
		functionName = "get_approx_yt_out_for_net_sy_in_internal"
	}else {
		return 0, 0, errors.New("swap type error！")
	}
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	typeArguments := []move_types.TypeTag{
		{Struct: syStructTag},
	}

	netSyInArg := CreatePureU64CallArg(netSyIn)
	netSyInArgument, err := ptb.Input(netSyInArg)
	if err != nil {
		return 0, 0, err
	}

	minPyOutArg := CreatePureU64CallArg(minYtOut)
	minPyOutArgument, err := ptb.Input(minPyOutArg)
	if err != nil {
		return 0, 0, err
	}

	oracleArgument, err := GetPriceVoucher(ptb, client, nemoConfig)
	if err != nil{
		return 0, 0, err
	}

	ps, err := GetObjectArgument(ptb, client, nemoConfig.PyState, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return 0, 0, err
	}
	ms, err := GetObjectArgument(ptb, client, nemoConfig.MarketState, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return 0, 0, err
	}
	mgc, err := GetObjectArgument(ptb, client, nemoConfig.MarketFactoryConfig, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return 0, 0, err
	}
	c, err := GetObjectArgument(ptb, client, constant.CLOCK, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return 0, 0, err
	}

	arguments := []sui_types.Argument{
		netSyInArgument,
		minPyOutArgument,
		*oracleArgument,
		ps,
		ms,
		mgc,
		c,
	}
	if exactPyType == constant.PTTYPE{
		arguments = []sui_types.Argument{
			netSyInArgument,
			minPyOutArgument,
			*oracleArgument,
			ps,
			mgc,
			ms,
			c,
		}
	}

	ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *nemoPackageId,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)

	pt := ptb.Finish()

	txKind := sui_types.TransactionKind{
		ProgrammableTransaction: &pt,
	}

	txBytes, err := bcs.Marshal(txKind)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to serialize transaction: %w", err)
	}

	senderAddr, err := sui_types.NewAddressFromHex(sender.Address)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse sender address: %w", err)
	}

	result, err := client.DevInspectTransactionBlock(context.Background(), *senderAddr, txBytes, nil, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to inspect transaction: %w", err)
	}
	if result.Error != nil{
		return 0, 0, errors.New(fmt.Sprintf("%v", *result.Error))
	}
	if len(result.Results) == 0 {
		return 0, 0, fmt.Errorf("no results returned")
	}

	lastResult := result.Results[len(result.Results)-1]

	if len(lastResult.ReturnValues) > 0 {
		firstValue := lastResult.ReturnValues[0]
		if firstValueArray, ok := firstValue.([]interface{}); ok && len(firstValueArray) > 0 {
			if innerArray, ok := firstValueArray[0].([]interface{}); ok && len(innerArray) > 0 {
				byteSlice := make([]byte, len(innerArray))
				for i, v := range innerArray {
					if num, ok := v.(float64); ok {
						byteSlice[i] = byte(num)
					}
				}
				if len(byteSlice) >= 8 {
					approxPyOut = binary.LittleEndian.Uint64(byteSlice)
					fmt.Printf("Parsed approxYtOut: %d\n", approxPyOut)
				}
			}
		}
	}

	if len(lastResult.ReturnValues) > 1 {
		secondValue := lastResult.ReturnValues[1]
		if secondValueArray, ok := secondValue.([]interface{}); ok && len(secondValueArray) > 0 {
			if innerArray, ok := secondValueArray[0].([]interface{}); ok && len(innerArray) > 0 {
				byteSlice := make([]byte, len(innerArray))
				for i, v := range innerArray {
					if num, ok := v.(float64); ok {
						byteSlice[i] = byte(num)
					}
				}
				if len(byteSlice) >= 8 {
					netSyTokenization = binary.LittleEndian.Uint64(byteSlice)
					fmt.Printf("Parsed netSyTokenization: %d\n", netSyTokenization)
				}
			}
		}
	}

	fmt.Printf("\nFinal values - approxYtOut: %d, netSyTokenization: %d\n", approxPyOut, netSyTokenization)

	return approxPyOut, netSyTokenization, nil
}

func DryRunGetPyOutForExactSyInWithPriceVoucher(client *client.Client, nemoConfig *models.NemoConfig, exactPyType string, netSyIn uint64, sender *account.Account) (uint64, error){
	ptb := sui_types.NewProgrammableTransactionBuilder()

	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.NemoContract)
	if err != nil {
		return 0, err
	}

	syStructTag, err := GetStructTag(nemoConfig.SyCoinType)
	if err != nil {
		return 0, err
	}

	moduleName := "router"
	var functionName string
	if exactPyType == constant.PTTYPE{
		functionName = "get_pt_out_for_exact_sy_in_with_price_voucher"
	}else {
		functionName = "get_yt_out_for_exact_sy_in_with_price_voucher"
	}

	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	typeArguments := []move_types.TypeTag{
		{Struct: syStructTag},
	}

	netSyInArg := CreatePureU64CallArg(netSyIn)
	netSyInArgument, err := ptb.Input(netSyInArg)
	if err != nil {
		return 0, err
	}

	minPyOut := uint64(0)
	minPyOutArg := CreatePureU64CallArg(minPyOut)
	minPyOutArgument, err := ptb.Input(minPyOutArg)
	if err != nil {
		return 0, err
	}

	oracleArgument, err := GetPriceVoucher(ptb, client, nemoConfig)
	if err != nil{
		return 0, err
	}

	ps, err := GetObjectArgument(ptb, client, nemoConfig.PyState, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return 0, err
	}
	ms, err := GetObjectArgument(ptb, client, nemoConfig.MarketState, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return 0, err
	}
	mgc, err := GetObjectArgument(ptb, client, nemoConfig.MarketFactoryConfig, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return 0, err
	}
	c, err := GetObjectArgument(ptb, client, constant.CLOCK, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return 0, err
	}

	arguments := []sui_types.Argument{
		netSyInArgument,
		minPyOutArgument,
		*oracleArgument,
		ps,
		mgc,
		ms,
		c,
	}

	ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *nemoPackageId,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)

	pt := ptb.Finish()

	txKind := sui_types.TransactionKind{
		ProgrammableTransaction: &pt,
	}

	txBytes, err := bcs.Marshal(txKind)
	if err != nil {
		return 0, fmt.Errorf("failed to serialize transaction: %w", err)
	}

	senderAddr, err := sui_types.NewAddressFromHex(sender.Address)
	if err != nil {
		return 0, fmt.Errorf("failed to parse sender address: %w", err)
	}

	result, err := client.DevInspectTransactionBlock(context.Background(), *senderAddr, txBytes, nil, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to inspect transaction: %w", err)
	}
	if result.Error != nil{
		return 0, errors.New(fmt.Sprintf("%v", *result.Error))
	}
	if len(result.Results) == 0 {
		return 0, fmt.Errorf("no results returned")
	}

	lastResult := result.Results[len(result.Results)-1]

	if len(lastResult.ReturnValues) > 0 {
		firstValue := lastResult.ReturnValues[0]
		if firstValueArray, ok := firstValue.([]interface{}); ok && len(firstValueArray) > 0 {
			if innerArray, ok := firstValueArray[0].([]interface{}); ok && len(innerArray) > 0 {
				byteSlice := make([]byte, len(innerArray))
				for i, v := range innerArray {
					if num, ok := v.(float64); ok {
						byteSlice[i] = byte(num)
					}
				}
				if len(byteSlice) >= 8 {
					minPyOut = binary.LittleEndian.Uint64(byteSlice)
					fmt.Printf("Parsed minYtOut: %d\n", minPyOut)
				}
			}
		}
	}
	return minPyOut, nil
}

func DryRunGetPyInForExactSyOutWithPriceVoucher(client *client.Client, nemoConfig *models.NemoConfig, exactPyType string, pyInAmount uint64, address string) (uint64, error){
	ptb := sui_types.NewProgrammableTransactionBuilder()

	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.NemoContract)
	if err != nil {
		return 0, err
	}

	syStructTag, err := GetStructTag(nemoConfig.SyCoinType)
	if err != nil {
		return 0, err
	}

	moduleName := "router"
	var functionName string
	if exactPyType == constant.PTTYPE{
		functionName = "get_sy_amount_out_for_exact_pt_in_with_price_voucher"
	}else {
		functionName = "get_sy_amount_out_for_exact_yt_in_with_price_voucher"
	}

	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	typeArguments := []move_types.TypeTag{
		{Struct: syStructTag},
	}

	exactPtInArg := CreatePureU64CallArg(pyInAmount)
	exactPtInArgument, err := ptb.Input(exactPtInArg)
	if err != nil {
		return 0, err
	}

	oracleArgument, err := GetPriceVoucher(ptb, client, nemoConfig)
	if err != nil{
		return 0, err
	}

	ps, err := GetObjectArgument(ptb, client, nemoConfig.PyState, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return 0, err
	}
	mgc, err := GetObjectArgument(ptb, client, nemoConfig.MarketFactoryConfig, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return 0, err
	}
	ms, err := GetObjectArgument(ptb, client, nemoConfig.MarketState, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return 0, err
	}
	c, err := GetObjectArgument(ptb, client, constant.CLOCK, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return 0, err
	}

	arguments := []sui_types.Argument{
		exactPtInArgument,
		*oracleArgument,
		ps,
		mgc,
		ms,
		c,
	}

	ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *nemoPackageId,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)

	pt := ptb.Finish()

	txKind := sui_types.TransactionKind{
		ProgrammableTransaction: &pt,
	}

	txBytes, err := bcs.Marshal(txKind)
	if err != nil {
		return 0, fmt.Errorf("failed to serialize transaction: %w", err)
	}

	senderAddr, err := sui_types.NewAddressFromHex(address)
	if err != nil {
		return 0, fmt.Errorf("failed to parse sender address: %w", err)
	}

	result, err := client.DevInspectTransactionBlock(context.Background(), *senderAddr, txBytes, nil, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to inspect transaction: %w", err)
	}
	if result.Error != nil{
		return 0, errors.New(fmt.Sprintf("%v", *result.Error))
	}
	if len(result.Results) == 0 {
		return 0, fmt.Errorf("no results returned")
	}

	lastResult := result.Results[len(result.Results)-1]

	var minSyOut uint64
	if len(lastResult.ReturnValues) > 0 {
		firstValue := lastResult.ReturnValues[0]
		if firstValueArray, ok := firstValue.([]interface{}); ok && len(firstValueArray) > 0 {
			if innerArray, ok := firstValueArray[0].([]interface{}); ok && len(innerArray) > 0 {
				byteSlice := make([]byte, len(innerArray))
				for i, v := range innerArray {
					if num, ok := v.(float64); ok {
						byteSlice[i] = byte(num)
					}
				}
				if len(byteSlice) >= 8 {
					minSyOut = binary.LittleEndian.Uint64(byteSlice)
					fmt.Printf("Parsed minSyOut: %d\n", minSyOut)
				}
			}
		}
	}
	return minSyOut, nil
}

func DryRunGetLpOutForSingleSyIn(client *client.Client, nemoConfig *models.NemoConfig, syInAmount uint64, sender *account.Account) (uint64, error){
	ptb := sui_types.NewProgrammableTransactionBuilder()

	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.NemoContract)
	if err != nil {
		return 0, err
	}

	syStructTag, err := GetStructTag(nemoConfig.SyCoinType)
	if err != nil {
		return 0, err
	}

	moduleName := "router"
	functionName := "get_lp_out_for_single_sy_in"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	typeArguments := []move_types.TypeTag{
		{Struct: syStructTag},
	}

	syInArg := CreatePureU64CallArg(syInAmount)
	syInArgument, err := ptb.Input(syInArg)
	if err != nil {
		return 0, err
	}

	oracleArgument, err := GetPriceVoucher(ptb, client, nemoConfig)
	if err != nil{
		return 0, err
	}

	ps, err := GetObjectArgument(ptb, client, nemoConfig.PyState, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return 0, err
	}
	mfc, err := GetObjectArgument(ptb, client, nemoConfig.MarketFactoryConfig, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return 0, err
	}
	ms, err := GetObjectArgument(ptb, client, nemoConfig.MarketState, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return 0, err
	}
	c, err := GetObjectArgument(ptb, client, constant.CLOCK, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return 0, err
	}

	arguments := []sui_types.Argument{
		syInArgument,
		*oracleArgument,
		ps,
		mfc,
		ms,
		c,
	}

	ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *nemoPackageId,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)

	pt := ptb.Finish()

	txKind := sui_types.TransactionKind{
		ProgrammableTransaction: &pt,
	}

	txBytes, err := bcs.Marshal(txKind)
	if err != nil {
		return 0, fmt.Errorf("failed to serialize transaction: %w", err)
	}

	senderAddr, err := sui_types.NewAddressFromHex(sender.Address)
	if err != nil {
		return 0, fmt.Errorf("failed to parse sender address: %w", err)
	}

	result, err := client.DevInspectTransactionBlock(context.Background(), *senderAddr, txBytes, nil, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to inspect transaction: %w", err)
	}
	if result.Error != nil{
		return 0, errors.New(fmt.Sprintf("%v", *result.Error))
	}
	if len(result.Results) == 0 {
		return 0, fmt.Errorf("no results returned")
	}

	lastResult := result.Results[len(result.Results)-1]

	var lpOut uint64
	if len(lastResult.ReturnValues) > 0 {
		firstValue := lastResult.ReturnValues[0]
		if firstValueArray, ok := firstValue.([]interface{}); ok && len(firstValueArray) > 0 {
			if innerArray, ok := firstValueArray[0].([]interface{}); ok && len(innerArray) > 0 {
				byteSlice := make([]byte, len(innerArray))
				for i, v := range innerArray {
					if num, ok := v.(float64); ok {
						byteSlice[i] = byte(num)
					}
				}
				if len(byteSlice) >= 8 {
					lpOut = binary.LittleEndian.Uint64(byteSlice)
					fmt.Printf("Parsed lpOut: %d\n", lpOut)
				}
			}
		}
	}
	return lpOut, nil
}

func DryRunSingleLiquidityAddPtOut(client *client.Client, nemoConfig *models.NemoConfig, syInAmount uint64, sender *account.Account) (uint64, error){
	ptb := sui_types.NewProgrammableTransactionBuilder()

	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.NemoContract)
	if err != nil {
		return 0, err
	}

	syStructTag, err := GetStructTag(nemoConfig.SyCoinType)
	if err != nil {
		return 0, err
	}

	moduleName := "offchain"
	functionName := "single_liquidity_add_pt_out"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	typeArguments := []move_types.TypeTag{
		{Struct: syStructTag},
	}

	syInArg := CreatePureU64CallArg(syInAmount)
	syInArgument, err := ptb.Input(syInArg)
	if err != nil {
		return 0, err
	}

	oracleArgument, err := GetPriceVoucher(ptb, client, nemoConfig)
	if err != nil{
		return 0, err
	}

	mfc, err := GetObjectArgument(ptb, client, nemoConfig.MarketFactoryConfig, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return 0, err
	}

	ps, err := GetObjectArgument(ptb, client, nemoConfig.PyState, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return 0, err
	}

	ms, err := GetObjectArgument(ptb, client, nemoConfig.MarketState, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return 0, err
	}

	c, err := GetObjectArgument(ptb, client, constant.CLOCK, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return 0, err
	}

	arguments := []sui_types.Argument{
		syInArgument,
		*oracleArgument,
		mfc,
		ps,
		ms,
		c,
	}

	ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *nemoPackageId,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)

	pt := ptb.Finish()

	txKind := sui_types.TransactionKind{
		ProgrammableTransaction: &pt,
	}

	txBytes, err := bcs.Marshal(txKind)
	if err != nil {
		return 0, fmt.Errorf("failed to serialize transaction: %w", err)
	}

	senderAddr, err := sui_types.NewAddressFromHex(sender.Address)
	if err != nil {
		return 0, fmt.Errorf("failed to parse sender address: %w", err)
	}

	result, err := client.DevInspectTransactionBlock(context.Background(), *senderAddr, txBytes, nil, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to inspect transaction: %w", err)
	}
	if result.Error != nil{
		return 0, errors.New(fmt.Sprintf("%v", *result.Error))
	}
	if len(result.Results) == 0 {
		return 0, fmt.Errorf("no results returned")
	}

	lastResult := result.Results[len(result.Results)-1]

	var ptValue uint64
	if len(lastResult.ReturnValues) > 0 {
		firstValue := lastResult.ReturnValues[0]
		if firstValueArray, ok := firstValue.([]interface{}); ok && len(firstValueArray) > 0 {
			if innerArray, ok := firstValueArray[0].([]interface{}); ok && len(innerArray) > 0 {
				byteSlice := make([]byte, len(innerArray))
				for i, v := range innerArray {
					if num, ok := v.(float64); ok {
						byteSlice[i] = byte(num)
					}
				}
				if len(byteSlice) >= 8 {
					ptValue = binary.LittleEndian.Uint64(byteSlice)
					fmt.Printf("Parsed ptValue: %d\n", ptValue)
				}
			}
		}
	}
	return ptValue, nil
}

func DryRunConversionRate(client *client.Client, nemoConfig *models.NemoConfig, address string) (float64, error){
	ptb := sui_types.NewProgrammableTransactionBuilder()

	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.OracleVoucherPackage)
	if err != nil {
		return 0, err
	}

	syStructTag, err := GetStructTag(nemoConfig.SyCoinType)
	if err != nil {
		return 0, err
	}

	moduleName := "oracle_voucher"
	functionName := "get_price"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	typeArguments := []move_types.TypeTag{
		{Struct: syStructTag},
	}

	oracleArgument, err := GetPriceVoucher(ptb, client, nemoConfig)
	if err != nil{
		return 0, err
	}

	arguments := []sui_types.Argument{
		*oracleArgument,
	}

	ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *nemoPackageId,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)

	pt := ptb.Finish()

	txKind := sui_types.TransactionKind{
		ProgrammableTransaction: &pt,
	}

	txBytes, err := bcs.Marshal(txKind)
	if err != nil {
		return 0, fmt.Errorf("failed to serialize transaction: %w", err)
	}

	senderAddr, err := sui_types.NewAddressFromHex(address)
	if err != nil {
		return 0, fmt.Errorf("failed to parse sender address: %w", err)
	}

	result, err := client.DevInspectTransactionBlock(context.Background(), *senderAddr, txBytes, nil, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to inspect transaction: %w", err)
	}
	if result.Error != nil{
		return 0, errors.New(fmt.Sprintf("%v", *result.Error))
	}
	if len(result.Results) == 0 {
		return 0, fmt.Errorf("no results returned")
	}

	lastResult := result.Results[len(result.Results)-1]

	var ptValue uint64
	if len(lastResult.ReturnValues) > 0 {
		firstValue := lastResult.ReturnValues[0]
		if firstValueArray, ok := firstValue.([]interface{}); ok && len(firstValueArray) > 0 {
			if innerArray, ok := firstValueArray[0].([]interface{}); ok && len(innerArray) > 0 {
				byteSlice := make([]byte, len(innerArray))
				for i, v := range innerArray {
					if num, ok := v.(float64); ok {
						byteSlice[i] = byte(num)
					}
				}
				if len(byteSlice) >= 8 {
					ptValue = binary.LittleEndian.Uint64(byteSlice)
					fmt.Printf("Parsed ptValue: %d\n", ptValue)
				}
			}
		}
	}


	return float64(ptValue) / math.Pow(2, 64) + 1, nil
}

func GetYtInAndSyOut(client *client.Client, nemoConfig *models.NemoConfig, address string, ytIn, retryTime uint64) (uint64, uint64, error){
	syOut, err := DryRunGetPyInForExactSyOutWithPriceVoucher(client, nemoConfig, constant.YTTYPE, ytIn, address)
	if err != nil{
		if retryTime > 3{
			return 0 , 0, err
		}
		return GetYtInAndSyOut(client, nemoConfig, address, ytIn / 100, retryTime + 1)
	}
	return ytIn, syOut, nil
}
func CalculateDailyEmission(emissionPerSecond, tokenType string, decimalPlaces int) float64 {
	emissionPerSecondDec, err := decimal.NewFromString(emissionPerSecond)
	if err != nil {
		fmt.Println("Invalid emissionPerSecond:", err)
		return 0
	}

	dailyEmission := emissionPerSecondDec.
		Mul(decimal.NewFromInt(60 * 60 * 24)).
		Div(decimal.NewFromFloat(math.Pow10(decimalPlaces)))

	dailyEmissionFloat, _ := dailyEmission.Float64()
	return dailyEmissionFloat
}

func GetRewarders(marketStateInfo map[string]interface{}, decimal int, sourceMarketState *MarketState, priceInfoMap map[string]PriceInfo) {
	byteBody,err := json.Marshal(marketStateInfo)
	if err != nil {
		return
	}
	var marketState RawMarketState
	err = json.Unmarshal(byteBody, &marketState)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	for _, content := range marketState.RewardPool.Fields.Rewarders.Fields.Contents {
		endtime,_ := strconv.ParseInt(content.Fields.Value.Fields.EndTime, 10 ,64)
		if endtime < time.Now().Unix() * 1000{
			continue
		}
		rewarder := content.Fields.Value.Fields
		dailyEmission := CalculateDailyEmission(rewarder.EmissionPerSecond, rewarder.RewardToken.Fields.Name, decimal)
		rewardName := fmt.Sprintf("0x%v",rewarder.RewardToken.Fields.Name)
		priceInfo, ok := priceInfoMap[rewardName]
		if !ok{
			continue
		}
		emissionActualDecimal,_ := strconv.ParseInt(priceInfo.Decimal, 10 ,64)
		if emissionActualDecimal == 0{
			continue
		}
		if emissionActualDecimal != int64(decimal) {
			dailyEmission = dailyEmission * math.Pow(10, float64(6)) / math.Pow(10, float64(emissionActualDecimal))
		}
		sourceMarketState.RewardMetrics = append(sourceMarketState.RewardMetrics, RewardMetric{
			TokenPrice: priceInfo.Price,
			TokenLogo: priceInfo.Logo,
			DailyEmission: fmt.Sprintf("%0.10f",dailyEmission),
			CoinType: rewardName,
			CoinName: priceInfo.Name,
			Decimal:  priceInfo.Decimal,
		})
	}
}

func GetYtInitInAmount(coinType string) uint64{
	switch coinType{
	case constant.SCALLOPSSUI:
		return 1000000
	case constant.SCALLOPDEEP:
		return 1000000
	case constant.SCALLOPUSDC:
		return 1000000
	case constant.SCALLOPSCA:
		return 1000000
	case constant.SCALLOPSBUSDT:
		return 1000000
	case constant.SCALLOPSBETH:
		return 100
	case constant.STSUI:
		return 1000000
	default:
		return 1000000
	}
}