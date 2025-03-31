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
)

func (s *SuiService)SwapByPy(amountIn, slippage float64, amountInType, exactAmountOutType string, sender *account.Account, nemoConfig *models.NemoConfig) (bool, error){
	netPyIn := uint64(amountIn * math.Pow(10, float64(nemoConfig.Decimal)))
	suiService := InitSuiService()
	minPyOut, err := api.DryRunGetPyInForExactSyOutWithPriceVoucher(suiService.SuiApi, nemoConfig, amountInType, netPyIn, sender.Address)
	if err != nil{
		return false, errors.New(fmt.Sprintf("%v",nemoError.ParseErrorMessage(err.Error())))
	}
	minPyOut = minPyOut - uint64(float64(minPyOut) * slippage)

	ptb := sui_types.NewProgrammableTransactionBuilder()
	oracleArgument, err := api.GetPriceVoucher(ptb, suiService.SuiApi, nemoConfig)
	if err != nil{
		return false, err
	}

	var swapArgument *sui_types.Argument
	if amountInType == constant.PTTYPE{
		swapArgument, err = api.SwapExactPtForSy(ptb, suiService.BlockApi, suiService.SuiApi, nemoConfig, netPyIn, minPyOut, sender.Address, oracleArgument)
		if err != nil{
			return false, err
		}
	} else if amountInType == constant.YTTYPE{
		swapArgument, err = api.SwapExactYtForSy(ptb, suiService.BlockApi, suiService.SuiApi, nemoConfig, netPyIn, minPyOut, sender.Address, oracleArgument)
		if err != nil{
			return false, err
		}
	} else {
		return false, errors.New("swap type error！")
	}

	syRedeemResult, err := api.SyRedeem(ptb, suiService.SuiApi, nemoConfig, swapArgument)
	if err != nil{
		return false, err
	}

	if exactAmountOutType != nemoConfig.CoinType{
		syRedeemResult, err = api.SwapToUnderlyingCoin(ptb, suiService.SuiApi, nemoConfig, syRedeemResult)
		if err != nil{
			return false, err
		}
	}

	recipientAddr, err := sui_types.NewAddressFromHex(sender.Address)
	if err != nil {
		return false, err
	}

	recArg, err := ptb.Pure(recipientAddr)
	if err != nil {
		return false, err
	}

	transferArgs := make([]sui_types.Argument, 0)
	// transfer object
	resultArg := &sui_types.Argument{
		NestedResult: &struct {
			Result1 uint16
			Result2 uint16
		}{Result1: *syRedeemResult.Result, Result2: 0},
	}
	transferArgs = append(transferArgs, *resultArg)

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

	_, gasCoin, err := api.RemainCoinAndGas(suiService.SuiApi, sender.Address, uint64(10000000), constant.GASCOINTYPE)
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
		10000000, // gasBudget
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

	resp, err := suiService.SuiApi.ExecuteTransactionBlock(
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

	return false, nil
}

func (s *SuiService)SwapToPy(amountIn, slippage float64, amountInType, exactAmountOutType string, sender *account.Account, nemoConfig *models.NemoConfig) (bool, error){
	client := InitSuiService()
	actualSyIn := uint64(amountIn * math.Pow(10, float64(nemoConfig.Decimal)))
	netSyIn := actualSyIn
	if amountInType == nemoConfig.UnderlyingCoinType{
		conversionRate,err := api.DryRunConversionRate(s.SuiApi, nemoConfig, "0x1")
		if err != nil{
			return false, err
		}
		netSyIn = uint64(float64(netSyIn) / conversionRate)
	}

	minPyOut, err := api.DryRunGetPyOutForExactSyInWithPriceVoucher(client.SuiApi, nemoConfig, exactAmountOutType, netSyIn, sender)
	if err != nil{
		return false, errors.New(fmt.Sprintf("%v",nemoError.ParseErrorMessage(err.Error())))
	}
	minPyOut = minPyOut - uint64(float64(minPyOut) * slippage)

	approxPyOut, netSyTokenization, err := api.DryRunGetApproxPyOutForNetSyInInternal(client.SuiApi, nemoConfig, exactAmountOutType, netSyIn, minPyOut, sender)
	if err != nil{
		return false, errors.New(fmt.Sprintf("%v",nemoError.ParseErrorMessage(err.Error())))
	}

	ptb := sui_types.NewProgrammableTransactionBuilder()

	remainingCoins, gasCoin, err := api.RemainCoinAndGas(client.SuiApi, sender.Address, uint64(10000000), amountInType)
	if err != nil{
		return false, err
	}

	splitResult ,_ ,err := api.SplitOrMergeCoin(ptb, client.SuiApi, remainingCoins, actualSyIn)
	if err != nil{
		return false, err
	}

	if !constant.IsGasCoinType(amountInType){
		_, gasCoin, err = api.RemainCoinAndGas(client.SuiApi, sender.Address, uint64(10000000), constant.GASCOINTYPE)
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

	depositArgument, err := api.Deposit(ptb, client.SuiApi, nemoConfig, &splitResult)
	if err != nil{
		return false, err
	}
	
	oracleArgument, err := api.GetPriceVoucher(ptb, client.SuiApi, nemoConfig)
	if err != nil{
		return false, err
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

	if exactAmountOutType == constant.YTTYPE{
		_, err = api.SwapExactSyForYt(ptb, client.BlockApi, client.SuiApi, nemoConfig, sender.Address, approxPyOut, netSyTokenization, minPyOut, oracleArgument, depositArgument, pyPositionArgument)
		if err != nil{
			return false, err
		}
	}else if exactAmountOutType == constant.PTTYPE{
		_, err = api.SwapExactSyForPt(ptb, client.BlockApi, client.SuiApi, nemoConfig, sender.Address, approxPyOut, minPyOut, oracleArgument, depositArgument, pyPositionArgument)
		if err != nil{
			return false, err
		}
	}else{
		return false, errors.New("swap type error！")
	}

	senderAddr, err := sui_types.NewObjectIdFromHex(sender.Address)
	if err != nil {
		return false, fmt.Errorf("failed to convert sender address: %w", err)
	}

	recArg, err := ptb.Pure(senderAddr)
	if err != nil {
		return false, err
	}

	if len(transferArgs) > 0{
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
	}

	pt := ptb.Finish()

	gasPayment := []*sui_types.ObjectRef{gasCoin}



	tx := sui_types.NewProgrammable(
		*senderAddr,
		gasPayment,
		pt,
		10000000, // gasBudget
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

	return false, nil
}