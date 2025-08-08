package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coming-chat/go-sui/v2/account"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
	"github.com/fardream/go-bcs/bcs"
	"github.com/nemo-protocol/nemo-go-sdk/service/sui/api"
	"github.com/nemo-protocol/nemo-go-sdk/service/sui/common/constant"
	"github.com/nemo-protocol/nemo-go-sdk/service/sui/common/models"
	"github.com/nemo-protocol/nemo-go-sdk/service/sui/common/nemoError"
	"math"
	"strconv"
	"strings"
)

func (s *SuiService)AddLiquidity(amountFloat, slippage float64, sender *account.Account, amountInType string, nemoConfig *models.NemoConfig)(bool, error){
	amountSyIn := uint64(amountFloat * math.Pow(10, float64(nemoConfig.Decimal)))
	actualSyIn := uint64(amountFloat * math.Pow(10, float64(nemoConfig.Decimal)))
	// create trade builder
	ptb := sui_types.NewProgrammableTransactionBuilder()
	client := InitSuiService()

	if amountInType == nemoConfig.UnderlyingCoinType{
		conversionRate,err := api.DryRunConversionRate(s.SuiApi, nemoConfig, "0x1")
		if err != nil{
			return false, err
		}
		amountSyIn = uint64(float64(amountSyIn) / conversionRate)
	}

	fmt.Printf("\n===amountSyIn:%v===\n",amountSyIn)
	minLpOut,err := api.DryRunGetLpOutForSingleSyIn(s.SuiApi, nemoConfig, amountSyIn, sender)
	var errorMsg string
	if err != nil {
		errorMsg = fmt.Sprintf("%v",nemoError.ParseErrorMessage(err.Error()))
		if !strings.Contains(errorMsg, "Market proportion too high"){
			return false, errors.New(errorMsg)
		}
	}

	var ptRate float64
	if errorMsg != ""{
		ptRate, _, err = api.JudgePtSyRate(client.SuiApi, nemoConfig, float64(amountSyIn))
	}
	fmt.Printf("\n==ptRate:%v, err:%v==\n", ptRate, err)
	fmt.Printf("\n===minLpOut:%v===\n",minLpOut)
	minLpOut = minLpOut - uint64(float64(minLpOut) * slippage)

	ptValue,err := api.DryRunSingleLiquidityAddPtOut(s.SuiApi, nemoConfig, amountSyIn, sender)
	if err != nil{
		return false, errors.New(fmt.Sprintf("%v",nemoError.ParseErrorMessage(err.Error())))
	}
	ptValue = ptValue - uint64(float64(ptValue) * slippage)

	remainingCoins, gasCoin, err := api.RemainCoinAndGas(client.SuiApi, sender.Address, uint64(30000000), amountInType)
	if err != nil{
		return false, err
	}

	amountIn := amountSyIn
	if amountInType == nemoConfig.UnderlyingCoinType {
		amountIn = actualSyIn
	}
	splitResult ,_ ,err := api.SplitOrMergeCoin(ptb, client.SuiApi, remainingCoins, amountIn)
	if err != nil{
		return false, err
	}

	if !constant.IsGasCoinType(amountInType){
		_, gasCoin, err = api.RemainCoinAndGas(client.SuiApi, sender.Address, uint64(30000000), constant.GASCOINTYPE)
		if err != nil{
			return false, err
		}
	}
	if amountInType == nemoConfig.UnderlyingCoinType {
		argument,err := api.MintToSCoin(ptb, client.SuiApi, nemoConfig, &splitResult)
		if err != nil{
			return false, err
		}
		splitResult = *argument
	}

	var ptSplitResult *sui_types.Argument
	var ptAmount float64
	if ptRate > 0{
		ptAmount = float64(amountIn) * ptRate
		fmt.Printf("\n==ptAmount:%v==\b",ptAmount)
		result, err := api.SplitCoinFromMerged(ptb, splitResult, uint64(ptAmount))
		if err != nil {
			return false, fmt.Errorf("failed to split merged coin: %w", err)
		}
		ptSplitResult = &result
	}

	pyPosition,err := api.GetPyPosition(nemoConfig, sender.Address, client.SuiApi, client.BlockApi)
	if err != nil {
		return false, err
	}

	var pyPositionArgument *sui_types.Argument
	// transfer object
	transferArgs := make([]sui_types.Argument, 0)
	if pyPosition == ""{
		pyPositionArgument, err = api.InitPyPosition(ptb, client.SuiApi, nemoConfig)
		if err != nil{
			return false, err
		}
		transferArgs = append(transferArgs, *pyPositionArgument)
	}else {
		argument, err := api.GetObjectArgument(ptb, client.SuiApi, pyPosition, false, nemoConfig.NemoContract, "", "")
		if err != nil{
			return false, err
		}
		pyPositionArgument = &argument
	}

	depositArgument, err := api.Deposit(ptb, client.SuiApi, nemoConfig, &splitResult)
	if err != nil{
		return false, err
	}

	var ptDepositArgument *sui_types.Argument
	if ptSplitResult != nil{
		ptDepositArgument, err = api.Deposit(ptb, client.SuiApi, nemoConfig, ptSplitResult)
		if err != nil{
			return false, err
		}
	}

	oracleArgument, err := api.GetPriceVoucher(ptb, client.SuiApi, nemoConfig)
	if err != nil{
		return false, err
	}

	var marketPosition *sui_types.Argument
	if ptSplitResult == nil {
		marketPosition, err = api.AddLiquiditySingleSy(ptb, client.SuiApi, nemoConfig, minLpOut, ptValue, oracleArgument, pyPositionArgument, depositArgument)
		if err != nil{
			return false, err
		}
	}else {
		_,err = api.MintPy(ptb, client.SuiApi, nemoConfig, ptDepositArgument, oracleArgument, pyPositionArgument)
		if err != nil{
			return false, err
		}
		var remainCoinResult *sui_types.Argument
		remainCoinResult, marketPosition, err = api.MintLp(ptb, client.SuiApi, nemoConfig, depositArgument, pyPositionArgument, int64(ptAmount))
		if err != nil{
			return false, err
		}
		if amountInType != nemoConfig.CoinType{
			underlyingToken, err := api.SwapToUnderlyingCoin(ptb, client.SuiApi, nemoConfig, remainCoinResult)
			if err != nil{
				return false, err
			}
			transferArgs = append(transferArgs, *underlyingToken)
		}else {
			transferArgs = append(transferArgs, *remainCoinResult)
		}
	}

	previousMarketPosition,err := api.GetMarketPosition(client.BlockApi, client.SuiApi, nemoConfig, sender.Address)
	if err != nil {
		return false, err
	}
	fmt.Printf("previousMarketPosition:%v\n",previousMarketPosition)

	if previousMarketPosition != "" {
		previousMarketPositionArgument,err := api.GetObjectArgument(ptb, client.SuiApi, previousMarketPosition, false, nemoConfig.NemoContract, "market_position", "join")
		if err != nil {
			return false, err
		}
		_, err = api.MergeAllLpPositions(ptb, client.SuiApi, nemoConfig, &previousMarketPositionArgument, marketPosition)
		if err != nil {
			return false, err
		}
		marketPosition = &previousMarketPositionArgument
	}

	// change recipient address
	recipientAddr, err := sui_types.NewAddressFromHex(sender.Address)
	if err != nil {
		return false, err
	}


	recArg, err := ptb.Pure(*recipientAddr)
	if err != nil {
		return false, err
	}

	// transfer object
	transferArgs = append(transferArgs, *marketPosition)
	ptb.Command(
		sui_types.Command{
			TransferObjects: &struct {
				Arguments []sui_types.Argument
				Argument  sui_types.Argument
			}{
				Arguments: transferArgs,
				Argument:  recArg,
			},
		},
	)

	pt := ptb.Finish()

	gasPayment := []*sui_types.ObjectRef{gasCoin}

	senderAddr, err := sui_types.NewObjectIdFromHex(sender.Address)
	if err != nil {
		return false, fmt.Errorf("failed to convert sender address: %w", err)
	}

	tx := sui_types.NewProgrammable(
		*senderAddr,
		gasPayment,
		pt,
		30000000, // gasBudget
		1000,     // gasPrice
	)

	txBytes, err := bcs.Marshal(tx)
	if err != nil {
		return false, fmt.Errorf("failed to serialize transaction: %w", err)
	}

	// signature
	signature, err := sender.SignSecureWithoutEncode(txBytes, sui_types.DefaultIntent())
	if err != nil {
		return false, fmt.Errorf("failed to sign transaction: %w", err)
	}

	options := types.SuiTransactionBlockResponseOptions{
		ShowInput:          true,
		ShowEffects:        true,
		ShowEvents:         true,
		ShowObjectChanges:  true,
		ShowBalanceChanges: true,
	}

	resp, err := client.SuiApi.ExecuteTransactionBlock(
		context.Background(),
		txBytes,
		[]any{signature},
		&options,
		types.TxnRequestTypeWaitForLocalExecution,
	)
	if err != nil {
		return false, fmt.Errorf("failed to execute transaction: %w", err)
	}

	b,_ := json.Marshal(resp.Effects.Data)
	fmt.Printf("\n==response:%+v==\n",resp)
	errorMsg = nemoError.GetError(string(b))
	if errorMsg != ""{
		return false, errors.New(errorMsg)
	}

	return true, nil
}

func (s *SuiService)RedeemLiquidity(amountIn, slippage float64, sender *account.Account, expectOutType string, nemoConfig *models.NemoConfig)(bool, error){
	amountLpIn := uint64(amountIn * math.Pow(10, float64(nemoConfig.Decimal)))
	// create trade builder
	ptb := sui_types.NewProgrammableTransactionBuilder()
	client := InitSuiService()

	pyPosition,err := api.GetPyPosition(nemoConfig, sender.Address, client.SuiApi, client.BlockApi)
	if err != nil {
		return false, err
	}

	if pyPosition == ""{
		return false, errors.New("pyPosition not found")
	}

	previousMarketPosition,err := api.GetMarketPosition(client.BlockApi, client.SuiApi, nemoConfig, sender.Address)
	if err != nil {
		return false, err
	}

	if previousMarketPosition == ""{
		return false, errors.New("marketPosition not found")
	}

	transferArgs := make([]sui_types.Argument, 0)
	syCoinArgument, err := api.BurnLp(ptb, client.SuiApi, nemoConfig, amountLpIn, pyPosition, previousMarketPosition)
	if err != nil {
		return false, err
	}

	yieldToken,err := api.SyRedeem(ptb, client.SuiApi, nemoConfig, syCoinArgument)
	if err != nil{
		return false, err
	}

	if expectOutType != nemoConfig.CoinType{
		underlyingToken, err := api.SwapToUnderlyingCoin(ptb, client.SuiApi, nemoConfig, yieldToken)
		if err != nil{
			return false, err
		}
		transferArgs = append(transferArgs, *underlyingToken)
	}else {
		transferArgs = append(transferArgs, *yieldToken)
	}


	// change recipient address
	recipientAddr, err := sui_types.NewAddressFromHex(sender.Address)
	if err != nil {
		return false, err
	}

	recArg, err := ptb.Pure(*recipientAddr)
	if err != nil {
		return false, err
	}

	ptb.Command(
		sui_types.Command{
			TransferObjects: &struct {
				Arguments []sui_types.Argument
				Argument  sui_types.Argument
			}{
				Arguments: transferArgs,
				Argument:  recArg,
			},
		},
	)

	pt := ptb.Finish()

	_, gasCoin, err := api.RemainCoinAndGas(client.SuiApi, sender.Address, uint64(30000000), constant.GASCOINTYPE)
	if err != nil{
		return false, err
	}

	gasPayment := []*sui_types.ObjectRef{gasCoin}

	senderAddr, err := sui_types.NewObjectIdFromHex(sender.Address)
	if err != nil {
		return false, fmt.Errorf("failed to convert sender address: %w", err)
	}

	tx := sui_types.NewProgrammable(
		*senderAddr,
		gasPayment,
		pt,
		30000000, // gasBudget
		1000,     // gasPrice
	)

	txBytes, err := bcs.Marshal(tx)
	if err != nil {
		return false, fmt.Errorf("failed to serialize transaction: %w", err)
	}

	// signature
	signature, err := sender.SignSecureWithoutEncode(txBytes, sui_types.DefaultIntent())
	if err != nil {
		return false, fmt.Errorf("failed to sign transaction: %w", err)
	}

	options := types.SuiTransactionBlockResponseOptions{
		ShowInput:          true,
		ShowEffects:        true,
		ShowEvents:         true,
		ShowObjectChanges:  true,
		ShowBalanceChanges: true,
	}

	resp, err := client.SuiApi.ExecuteTransactionBlock(
		context.Background(),
		txBytes,
		[]any{signature},
		&options,
		types.TxnRequestTypeWaitForLocalExecution,
	)
	if err != nil {
		return false, fmt.Errorf("failed to execute transaction: %w", err)
	}

	b,_ := json.Marshal(resp.Effects.Data)
	fmt.Printf("\n==response:%+v==\n",resp)
	errorMsg := nemoError.GetError(string(b))
	if errorMsg != ""{
		return false, errors.New(errorMsg)
	}

	return true, nil
}

func (s *SuiService)DryRunYtReward(nemoConfig *models.NemoConfig, sender *account.Account) (float64, error){
	ptb := sui_types.NewProgrammableTransactionBuilder()
	client := InitSuiService()

	pyPosition, err := api.GetPyPosition(nemoConfig, sender.Address, client.SuiApi, client.BlockApi)
	if err != nil{
		return 0, err
	}
	fmt.Printf("pyposition:%v",pyPosition)

	if pyPosition == ""{
		return 0, errors.New("pyPosition not found")
	}

	oracleArgument, err := api.GetPriceVoucher(ptb, client.SuiApi, nemoConfig)
	if err != nil{
		return 0, err
	}

	syArgument, err := api.RedeemDueInterest(ptb, client.SuiApi, nemoConfig, pyPosition, oracleArgument)
	if err != nil{
		return 0, err
	}

	_, err = api.SyRedeem(ptb, client.SuiApi, nemoConfig, syArgument)

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

	result, err := client.SuiApi.DevInspectTransactionBlock(context.Background(), *senderAddr, txBytes, nil, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to inspect transaction: %w", err)
	}
	if result.Error != nil{
		return 0, errors.New(fmt.Sprintf("%v", *result.Error))
	}
	if len(result.Results) == 0 {
		return 0, fmt.Errorf("no results returned")
	}

	if len(result.Results) == 0 {
		return 0, fmt.Errorf("empty results")
	}
	lastResult := result.Results[len(result.Results)-1]

	if len(lastResult.ReturnValues) == 0 {
		return 0, fmt.Errorf("no return values")
	}

	fmt.Printf("lastResult:%+v",lastResult)

	firstValue := lastResult.ReturnValues[0]
	fmt.Printf("\n==firstValue:%v==\n",firstValue)
	var coin models.Coin
	if firstValueArray, ok := firstValue.([]interface{}); ok && len(firstValueArray) > 0 {
		if innerArray, ok := firstValueArray[0].([]interface{}); ok && len(innerArray) > 0 {
			byteSlice := make([]byte, len(innerArray))
			for i, v := range innerArray {
				if num, ok := v.(float64); ok {
					byteSlice[i] = byte(num)
				}
			}
			_,err = bcs.Unmarshal(byteSlice, &coin)
			if err != nil {
				return 0, err
			}
			fmt.Printf("Coin Value: %d\n", coin.Value)
		}
	}

	decimalPow := math.Pow(10, float64(nemoConfig.Decimal))

	return float64(coin.Value) / decimalPow, nil
}

func (s *SuiService)ClaimYtReward(nemoConfig *models.NemoConfig, sender *account.Account) (bool, error){
	ptb := sui_types.NewProgrammableTransactionBuilder()
	client := InitSuiService()

	pyPosition, err := api.GetPyPosition(nemoConfig, sender.Address, client.SuiApi, client.BlockApi)
	if err != nil{
		return false, err
	}
	fmt.Printf("pyposition:%v",pyPosition)

	if pyPosition == ""{
		return false, errors.New("pyPosition not found")
	}

	oracleArgument, err := api.GetPriceVoucher(ptb, client.SuiApi, nemoConfig)
	if err != nil{
		return false, err
	}

	syCoinArgument, err := api.RedeemDueInterest(ptb, client.SuiApi, nemoConfig, pyPosition, oracleArgument)
	if err != nil{
		return false, err
	}

	sCoinArgument,err := api.SyRedeem(ptb, client.SuiApi, nemoConfig, syCoinArgument)
	if err != nil{
		return false, err
	}

	transferArgs := make([]sui_types.Argument, 0)
	if constant.IsScallopCoin(nemoConfig.CoinType) || nemoConfig.ProviderProtocol == constant.SCALLOP{
		underlyingCoinArgument, err := api.BurnSCoin(ptb, client.SuiApi, nemoConfig.CoinType, nemoConfig.UnderlyingCoinType, sCoinArgument)
		if err != nil{
			return false, err
		}

		transferArgs = append(transferArgs, *underlyingCoinArgument)
	}else {
		transferArgs = append(transferArgs, *sCoinArgument)
	}

	// change recipient address
	recipientAddr, err := sui_types.NewAddressFromHex(sender.Address)
	if err != nil {
		return false, err
	}

	recArg, err := ptb.Pure(*recipientAddr)
	if err != nil {
		return false, err
	}

	ptb.Command(
		sui_types.Command{
			TransferObjects: &struct {
				Arguments []sui_types.Argument
				Argument  sui_types.Argument
			}{
				Arguments: transferArgs,
				Argument:  recArg,
			},
		},
	)

	pt := ptb.Finish()

	_, gasCoin, err := api.RemainCoinAndGas(client.SuiApi, sender.Address, uint64(30000000), constant.GASCOINTYPE)
	if err != nil{
		return false, err
	}

	gasPayment := []*sui_types.ObjectRef{gasCoin}

	senderAddr, err := sui_types.NewObjectIdFromHex(sender.Address)
	if err != nil {
		return false, fmt.Errorf("failed to convert sender address: %w", err)
	}

	tx := sui_types.NewProgrammable(
		*senderAddr,
		gasPayment,
		pt,
		30000000, // gasBudget
		1000,     // gasPrice
	)

	txBytes, err := bcs.Marshal(tx)
	if err != nil {
		return false, fmt.Errorf("failed to serialize transaction: %w", err)
	}

	// signature
	signature, err := sender.SignSecureWithoutEncode(txBytes, sui_types.DefaultIntent())
	if err != nil {
		return false, fmt.Errorf("failed to sign transaction: %w", err)
	}

	options := types.SuiTransactionBlockResponseOptions{
		ShowInput:          true,
		ShowEffects:        true,
		ShowEvents:         true,
		ShowObjectChanges:  true,
		ShowBalanceChanges: true,
	}

	resp, err := client.SuiApi.ExecuteTransactionBlock(
		context.Background(),
		txBytes,
		[]any{signature},
		&options,
		types.TxnRequestTypeWaitForLocalExecution,
	)
	if err != nil {
		return false, fmt.Errorf("failed to execute transaction: %w", err)
	}

	b,_ := json.Marshal(resp.Effects.Data)
	fmt.Printf("\n==response:%+v==\n",resp)
	errorMsg := nemoError.GetError(string(b))
	if errorMsg != ""{
		return false, errors.New(errorMsg)
	}

	return true, nil
}

func (s *SuiService)DryRunLpReward(nemoConfig *models.NemoConfig, sender *account.Account) ([]models.CoinInfo, error){
	ptb := sui_types.NewProgrammableTransactionBuilder()
	client := InitSuiService()
	innerPriceInfoMap := api.GetCoinPriceInfo()

	coinInfoList := make([]models.CoinInfo, 0)
	for _,v := range nemoConfig.Incentives {
		marketPosition, err := api.GetMarketPosition(client.BlockApi, client.SuiApi, nemoConfig, sender.Address)
		if err != nil {
			return coinInfoList, err
		}

		if marketPosition == "" {
			return coinInfoList, errors.New("marketPosition not found")
		}

		rewardArgument, err := api.ClaimReward(ptb, client.SuiApi, nemoConfig, marketPosition, v.TokenType)
		if err != nil {
			return coinInfoList, err
		}

		// change recipient address
		recipientAddr, err := sui_types.NewAddressFromHex(sender.Address)
		if err != nil {
			return coinInfoList, err
		}

		recArg, err := ptb.Pure(*recipientAddr)
		if err != nil {
			return coinInfoList, err
		}

		transferArgs := make([]sui_types.Argument, 0)
		transferArgs = append(transferArgs, *rewardArgument)

		ptb.Command(
			sui_types.Command{
				TransferObjects: &struct {
					Arguments []sui_types.Argument
					Argument  sui_types.Argument
				}{
					Arguments: transferArgs,
					Argument:  recArg,
				},
			},
		)

		pt := ptb.Finish()
		txKind := sui_types.TransactionKind{
			ProgrammableTransaction: &pt,
		}

		txBytes, err := bcs.Marshal(txKind)
		if err != nil {
			return coinInfoList, fmt.Errorf("failed to serialize transaction: %w", err)
		}

		senderAddr, err := sui_types.NewAddressFromHex(sender.Address)
		if err != nil {
			return coinInfoList, fmt.Errorf("failed to parse sender address: %w", err)
		}

		result, err := client.SuiApi.DevInspectTransactionBlock(context.Background(), *senderAddr, txBytes, nil, nil)
		if err != nil {
			return coinInfoList, fmt.Errorf("failed to inspect transaction: %w", err)
		}
		if result.Error != nil {
			return coinInfoList, errors.New(fmt.Sprintf("%v", *result.Error))
		}
		if len(result.Results) == 0 {
			return coinInfoList, fmt.Errorf("no results returned")
		}

		if len(result.Results) == 0 {
			return coinInfoList, fmt.Errorf("empty results")
		}
		lastResult := result.Results[0]

		if len(lastResult.ReturnValues) == 0 {
			return coinInfoList, fmt.Errorf("no return values")
		}

		fmt.Printf("lastResult:%+v", lastResult)

		firstValue := lastResult.ReturnValues[0]
		fmt.Printf("\n==firstValue:%v==\n", firstValue)
		var coin models.Coin
		var coinInfo models.CoinInfo
		if firstValueArray, ok := firstValue.([]interface{}); ok && len(firstValueArray) > 0 {
			if innerArray, ok := firstValueArray[0].([]interface{}); ok && len(innerArray) > 0 {
				byteSlice := make([]byte, len(innerArray))
				for i, v := range innerArray {
					if num, ok := v.(float64); ok {
						byteSlice[i] = byte(num)
					}
				}
				_, err = bcs.Unmarshal(byteSlice, &coin)
				if err != nil {
					return coinInfoList, err
				}
				fmt.Printf("Coin Value: %d\n", coin.Value)
				coinInfo.CoinType = v.TokenType
				decimal := int64(0)
				for coinType, coinPriceInfo := range innerPriceInfoMap{
					if coinType == v.TokenType{
						decimal,_ = strconv.ParseInt(coinPriceInfo.Decimal, 10, 64)
						break
					}
				}
				coinInfo.Amount = float64(coin.Value) / math.Pow(10, float64(decimal))
				coinInfoList = append(coinInfoList, coinInfo)
			}
		}
	}

	return coinInfoList, nil
}

func (s *SuiService)ClaimLpReward(nemoConfig *models.NemoConfig, sender *account.Account) (bool, error){
	client := InitSuiService()

	for _, incentive := range nemoConfig.Incentives{
		ptb := sui_types.NewProgrammableTransactionBuilder()
		marketPosition, err := api.GetMarketPosition(client.BlockApi, client.SuiApi, nemoConfig, sender.Address)
		if err != nil {
			return false, err
		}

		if marketPosition == "" {
			return false, errors.New("marketPosition not found")
		}

		rewardArgument, err := api.ClaimReward(ptb, client.SuiApi, nemoConfig, marketPosition, incentive.TokenType)
		if err != nil {
			return false, err
		}

		// change recipient address
		recipientAddr, err := sui_types.NewAddressFromHex(sender.Address)
		if err != nil {
			return false, err
		}

		recArg, err := ptb.Pure(*recipientAddr)
		if err != nil {
			return false, err
		}

		transferArgs := make([]sui_types.Argument, 0)
		transferArgs = append(transferArgs, *rewardArgument)

		ptb.Command(
			sui_types.Command{
				TransferObjects: &struct {
					Arguments []sui_types.Argument
					Argument  sui_types.Argument
				}{
					Arguments: transferArgs,
					Argument:  recArg,
				},
			},
		)

		pt := ptb.Finish()

		_, gasCoin, err := api.RemainCoinAndGas(client.SuiApi, sender.Address, uint64(30000000), constant.GASCOINTYPE)
		if err != nil{
			return false, err
		}

		gasPayment := []*sui_types.ObjectRef{gasCoin}

		senderAddr, err := sui_types.NewObjectIdFromHex(sender.Address)
		if err != nil {
			return false, fmt.Errorf("failed to convert sender address: %w", err)
		}

		tx := sui_types.NewProgrammable(
			*senderAddr,
			gasPayment,
			pt,
			30000000, // gasBudget
			1000,     // gasPrice
		)

		txBytes, err := bcs.Marshal(tx)
		if err != nil {
			return false, fmt.Errorf("failed to serialize transaction: %w", err)
		}

		// signature
		signature, err := sender.SignSecureWithoutEncode(txBytes, sui_types.DefaultIntent())
		if err != nil {
			return false, fmt.Errorf("failed to sign transaction: %w", err)
		}

		options := types.SuiTransactionBlockResponseOptions{
			ShowInput:          true,
			ShowEffects:        true,
			ShowEvents:         true,
			ShowObjectChanges:  true,
			ShowBalanceChanges: true,
		}

		resp, err := client.SuiApi.ExecuteTransactionBlock(
			context.Background(),
			txBytes,
			[]any{signature},
			&options,
			types.TxnRequestTypeWaitForLocalExecution,
		)
		if err != nil {
			return false, fmt.Errorf("failed to execute transaction: %w", err)
		}

		b,_ := json.Marshal(resp.Effects.Data)
		fmt.Printf("\n==response:%+v==\n",resp)
		errorMsg := nemoError.GetError(string(b))
		if errorMsg != ""{
			return false, errors.New(errorMsg)
		}
	}

	return true, nil
}

func (s *SuiService)QueryPoolApy(nemoConfig *models.NemoConfig, priceInfoMap ...map[string]api.PriceInfo) (*models.ApyModel, error){

	conversionRate,err := api.DryRunConversionRate(s.SuiApi, nemoConfig, "0x1")
	if err != nil{
		return nil, errors.New(fmt.Sprintf("%v",nemoError.ParseErrorMessage(err.Error())))
	}

	coinPrice,err := strconv.ParseFloat(nemoConfig.CoinPrice, 64)
	if err != nil{
		return nil, err
	}
	underlyingPrice := coinPrice / conversionRate
	nemoConfig.UnderlyingCoinPrice = fmt.Sprintf("%0.10f",underlyingPrice)

	syIn, pyOut, err := api.GetSyInAndPyOut(s.SuiApi, nemoConfig, "0x1", api.GetYtInitInAmount(nemoConfig.CoinType))
	if err != nil{
		return nil, errors.New(fmt.Sprintf("%v",nemoError.ParseErrorMessage(err.Error())))
	}
	fmt.Printf("ytin:%v, syout:%v",pyOut, syIn)

	objectList := []string{
		nemoConfig.PyState,
		nemoConfig.MarketState,
	}
	objectMap, err := api.MultiGetObjectFieldByObjectId(s.SuiApi, objectList)
	if err != nil{
		return nil, err
	}

	pyStateInfo, err := api.SuiResponseToMap(objectMap[nemoConfig.PyState])
	if err != nil{
		return nil, err
	}
	maturity := pyStateInfo["expiry"].(string)

	marketStateInfo, err := api.SuiResponseToMap(objectMap[nemoConfig.MarketState])
	if err != nil{
		return nil, err
	}

	innerPriceInfoMap := map[string]api.PriceInfo{}
	if len(priceInfoMap) == 0 || len(priceInfoMap[0]) == 0{
		innerPriceInfoMap = api.GetCoinPriceInfo()
	}else {
		innerPriceInfoMap = priceInfoMap[0]
	}

	marketState := api.MarketState{}
	api.GetRewarders(marketStateInfo, int(nemoConfig.Decimal), &marketState, innerPriceInfoMap)

	coinInfo := api.CoinInfo{}
	coinInfo.CoinPrice = coinPrice
	coinInfo.Decimal = nemoConfig.Decimal
	coinInfo.UnderlyingPrice = underlyingPrice
	coinInfo.UnderlyingApy,_ = strconv.ParseFloat(nemoConfig.UnderlyingApy, 64)
	coinInfo.Maturity,_ = strconv.ParseInt(maturity, 10, 64)
	coinInfo.SwapFeeForLpHolder,_ = strconv.ParseFloat(nemoConfig.SwapFeeForLpHolder, 64)
	coinInfo.SwapFeeForLpHolderAvg,_ = strconv.ParseFloat(nemoConfig.SwapFeeForLpHolderAvg, 64)

	marketState.TotalPt = marketStateInfo["total_pt"].(string)
	marketState.LpSupply = marketStateInfo["lp_supply"].(string)
	marketState.TotalSy = marketStateInfo["total_sy"].(string)
	marketState.MarketCap = marketStateInfo["market_cap"].(string)
	response := api.CalculatePoolApy(coinInfo, marketState, int64(pyOut), int64(syIn))

	return response, nil
}

