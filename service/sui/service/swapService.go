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
	"math"
	"nemo-go-sdk/service/sui/api"
	"nemo-go-sdk/service/sui/common/constant"
	"nemo-go-sdk/service/sui/common/models"
	"strconv"
)

func (s *SuiService)SwapByPy(amountIn, slippage float64, amountInType, exactAmountOutType string, sender *account.Account, nemoConfig *models.NemoConfig) (bool, error){
	ptb := sui_types.NewProgrammableTransactionBuilder()
	suiService := InitSuiService()

	oracleArgument, err := api.GetPriceVoucher(ptb, suiService.SuiApi, nemoConfig)
	if err != nil{
		return false, err
	}

	var swapArgument *sui_types.Argument
	if amountInType == constant.PTTYPE{
		swapArgument, err = api.SwapExactPtForSy(ptb, suiService.BlockApi, suiService.SuiApi, nemoConfig, sender.Address, oracleArgument)
		if err != nil{
			return false, err
		}
	} else if amountInType == constant.YTTYPE{
		swapArgument, err = api.SwapExactYtForSy(ptb, suiService.BlockApi, suiService.SuiApi, nemoConfig, sender.Address, oracleArgument)
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

	//coin, err := api.BurnSCoin(ptb, suiService.SuiApi, COINTYPE, UNDERLYINGCOINTYPE, syRedeemResult)
	//if err != nil{
	//	return false, err
	//}
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

	b,_ := json.Marshal(resp)
	fmt.Printf("\n==resp:%+v==\n",string(b))

	return false, nil
}

func (s *SuiService)SwapToPy(amountIn, slippage float64, amountInType, exactAmountOutType string, sender *account.Account, nemoConfig *models.NemoConfig) (bool, error){
	client := InitSuiService()
	actualSyIn := uint64(amountIn * math.Pow(10, float64(nemoConfig.Decimal)))
	netSyIn := actualSyIn
	if amountInType == nemoConfig.UnderlyingCoinType{
		conversionRate,err := strconv.ParseFloat(nemoConfig.ConversionRate, 10)
		if err != nil{
			return false, err
		}
		netSyIn = uint64(float64(netSyIn) / conversionRate)
	}

	minPyOut, err := api.DryRunGetPyOutForExactSyInWithPriceVoucher(client.SuiApi, nemoConfig, exactAmountOutType, netSyIn, nemoConfig.PriceOracle, sender)
	if err != nil{
		return false, err
	}
	minPyOut = minPyOut - uint64(float64(minPyOut) * slippage)

	approxPyOut, netSyTokenization, err := api.DryRunGetApproxPyOutForNetSyInInternal(client.SuiApi, nemoConfig, exactAmountOutType, netSyIn, minPyOut, sender)
	if err != nil{
		return false, err
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
	}else {
		if constant.IsScallopCoin(nemoConfig.CoinType){
			marketCoinArgument,err := api.Mint(ptb, client.SuiApi, nemoConfig.UnderlyingCoinType, &splitResult)
			if err != nil{
				return false, err
			}
			argument,err := api.MintSCoin(ptb, client.SuiApi, nemoConfig.CoinType, nemoConfig.UnderlyingCoinType, marketCoinArgument)
			if err != nil{
				return false, err
			}
			splitResult = *argument
		}
	}

	depositArgument, err := api.Deposit(ptb, client.SuiApi, nemoConfig, &splitResult)
	if err != nil{
		return false, err
	}
	
	oracleArgument, err := api.GetPriceVoucher(ptb, client.SuiApi, nemoConfig)
	if err != nil{
		return false, err
	}

	if exactAmountOutType == constant.YTTYPE{
		_, err = api.SwapExactSyForYt(ptb, client.BlockApi, client.SuiApi, nemoConfig, sender.Address, approxPyOut, netSyTokenization, minPyOut, oracleArgument, depositArgument)
		if err != nil{
			return false, err
		}
	}else if exactAmountOutType == constant.PTTYPE{
		_, err = api.SwapExactSyForPt(ptb, client.BlockApi, client.SuiApi, nemoConfig, sender.Address, approxPyOut, minPyOut, oracleArgument, depositArgument)
		if err != nil{
			return false, err
		}
	}else{
		return false, errors.New("swap type error！")
	}


	// transfer object
	//transferArgs := []sui_types.Argument{remainMergeCoinArgument}

	senderAddr, err := sui_types.NewObjectIdFromHex(sender.Address)
	if err != nil {
		return false, fmt.Errorf("failed to convert sender address: %w", err)
	}

	//recArg, err := ptb.Pure(senderAddr)
	//if err != nil {
	//	return false, err
	//}

	//ptb.Command(
	//	sui_types.Command{
	//		TransferObjects: &struct {
	//			Arguments []sui_types.Argument
	//			Argument  sui_types.Argument
	//		}{
	//			Arguments: transferArgs,
	//			Argument:  recArg,
	//		},
	//	},
	//)

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

	b,_ := json.Marshal(resp)
	fmt.Printf("\n==resp:%+v==\n",string(b))

	return false, nil
}