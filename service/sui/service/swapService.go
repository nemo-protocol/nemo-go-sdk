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
	"nemo-go-sdk/service/sui/api"
	"nemo-go-sdk/service/sui/common/constant"
)

func (s *SuiService)SwapByPy(amountIn, slippage float64, coinType, amountInType, exactAmountOutType string, sender *account.Account) (bool, error){
	ptb := sui_types.NewProgrammableTransactionBuilder()
	suiService := InitSuiService()

	oracleArgument, err := api.GetPriceVoucherFromXOracle(ptb, suiService.SuiApi, NEMOPACKAGE, SYTYPE, UNDERLYINGCOINTYPE)
	if err != nil{
		return false, err
	}

	pyState := "0x60422aa99f040c7ac8d0071a3bfd5431bd05b3ad82c77636761eab2709681fde"
	nemoPackageList := []string{"0xbde9dd9441697413cf312a2d4e37721f38814b96d037cb90d5af10b79de1d446", NEMOPACKAGE}

	var swapArgument *sui_types.Argument
	if amountInType == constant.PTTYPE{
		swapArgument, err = api.SwapExactPtForSy(ptb, suiService.BlockApi, suiService.SuiApi, NEMOPACKAGE, pyState, SYTYPE, sender.Address, nemoPackageList, oracleArgument)
		if err != nil{
			return false, err
		}
	} else if amountInType == constant.YTTYPE{
		swapArgument, err = api.SwapExactYtForSy(ptb, suiService.BlockApi, suiService.SuiApi, NEMOPACKAGE, pyState, SYTYPE, sender.Address, nemoPackageList, oracleArgument)
		if err != nil{
			return false, err
		}
	} else {
		return false, errors.New("swap type error！")
	}


	syRedeemResult, err := api.SyRedeem(ptb, suiService.SuiApi, NEMOPACKAGE, COINTYPE, SYTYPE, swapArgument)
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

	_, gasCoin, err := api.RemainCoinAndGas(suiService.SuiApi, sender.Address, uint64(10000000), coinType)
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

func (s *SuiService)SwapToPy(amountIn, slippage float64, coinType, amountInType, exactAmountOutType string, sender *account.Account) (bool, error){
	client := InitSuiService()
	netSyIn := uint64(amountIn*1000000000)

	minYtOut, err := api.GetYtOutForExactSyInWithPriceVoucher(client.SuiApi, NEMOPACKAGE, SYTYPE, api.PYSTATE, api.MARKETGLOBALCONFIG, api.MARKETSTATECONFIG,netSyIn, api.PRICEORACLECONFIG, sender)
	if err != nil{
		return false, err
	}
	minYtOut = minYtOut - uint64(float64(minYtOut) * slippage)

	approxYtOut, netSyTokenization, err := api.DryRunGetApproxYtOutForNetSyInInternal(client.SuiApi, NEMOPACKAGE, SYTYPE, api.PYSTATE, api.MARKETGLOBALCONFIG, api.MARKETSTATECONFIG, netSyIn, minYtOut, sender)
	if err != nil{
		return false, err
	}

	ptb := sui_types.NewProgrammableTransactionBuilder()
	
	remainingCoins, gasCoin, err := api.RemainCoinAndGas(client.SuiApi, sender.Address, uint64(10000000), coinType)
	if err != nil{
		return false, err
	}
	
	if !constant.IsGasCoinType(coinType){
		_, gasCoin, err = api.RemainCoinAndGas(client.SuiApi, sender.Address, uint64(10000000), constant.GASCOINTYPE)
		if err != nil{
			return false, err
		}
	}

	mergeCoinArgument, remainingCoins, err := api.MergeCoin(ptb, client.SuiApi, remainingCoins, netSyIn)
	if err != nil{
		return false, err
	}

	// 执行 SplitCoin 操作
	splitResult,_,err := api.SplitCoinFromMerged(ptb, *mergeCoinArgument[0], netSyIn)
	if err != nil{
		return false, err
	}

	depositArgument, err := api.Deposit(ptb, client.SuiApi, NEMOPACKAGE, coinType, SYTYPE, &splitResult)
	if err != nil{
		return false, err
	}
	
	oracleArgument, err := api.GetPriceVoucherFromXOracle(ptb, client.SuiApi, NEMOPACKAGE, SYTYPE, constant.GASCOINTYPE)
	if err != nil{
		return false, err
	}

	nemoPackageList := []string{"0xbde9dd9441697413cf312a2d4e37721f38814b96d037cb90d5af10b79de1d446", NEMOPACKAGE}
	_, err = api.SwapExactSyForYt(ptb, client.BlockApi, client.SuiApi, NEMOPACKAGE, api.PYSTATE, SYTYPE, sender.Address, nemoPackageList, approxYtOut, netSyTokenization, minYtOut, oracleArgument, depositArgument)
	if err != nil{
		return false, err
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