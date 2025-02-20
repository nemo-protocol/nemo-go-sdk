package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"nemo-go-sdk/service/sui/api"
	"nemo-go-sdk/service/sui/common/constant"
	"nemo-go-sdk/service/sui/common/models"

	"github.com/coming-chat/go-sui/v2/account"
	"github.com/coming-chat/go-sui/v2/move_types"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
	"github.com/fardream/go-bcs/bcs"
)

func (s *SuiService) SwapByPy(amountIn, slippage float64, amountInType, exactAmountOutType string, sender *account.Account, nemoConfig *models.NemoConfig) (bool, error) {
	ptb := sui_types.NewProgrammableTransactionBuilder()
	suiService := InitSuiService()

	oracleArgument, err := api.GetPriceVoucher(ptb, suiService.SuiApi, nemoConfig)
	if err != nil {
		return false, err
	}

	var swapArgument *sui_types.Argument
	if amountInType == constant.PTTYPE {
		swapArgument, err = api.SwapExactPtForSy(ptb, suiService.BlockApi, suiService.SuiApi, nemoConfig, sender.Address, oracleArgument)
		if err != nil {
			return false, err
		}
	} else if amountInType == constant.YTTYPE {
		swapArgument, err = api.SwapExactYtForSy(ptb, suiService.BlockApi, suiService.SuiApi, nemoConfig, sender.Address, oracleArgument)
		if err != nil {
			return false, err
		}
	} else {
		return false, errors.New("swap type error！")
	}

	syRedeemResult, err := api.SyRedeem(ptb, suiService.SuiApi, nemoConfig, swapArgument)
	if err != nil {
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
	if err != nil {
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

	b, _ := json.Marshal(resp)
	fmt.Printf("\n==resp:%+v==\n", string(b))

	return false, nil
}

func (s *SuiService) SwapToPy(amountIn, slippage float64, amountInType, exactAmountOutType string, sender *account.Account, nemoConfig *models.NemoConfig) (bool, error) {
	client := InitSuiService()
	netSyIn := uint64(amountIn * 1000000000)

	minYtOut, err := api.DryRunGetPyOutForExactSyInWithPriceVoucher(client.SuiApi, nemoConfig, exactAmountOutType, netSyIn, nemoConfig.PriceOracle, sender)
	if err != nil {
		return false, err
	}
	minYtOut = minYtOut - uint64(float64(minYtOut)*slippage)

	approxPyOut, netSyTokenization, err := api.DryRunGetApproxPyOutForNetSyInInternal(client.SuiApi, nemoConfig, exactAmountOutType, netSyIn, minYtOut, sender)
	if err != nil {
		return false, err
	}

	ptb := sui_types.NewProgrammableTransactionBuilder()

	remainingCoins, gasCoin, err := api.RemainCoinAndGas(client.SuiApi, sender.Address, uint64(10000000), amountInType)
	if err != nil {
		return false, err
	}

	if !constant.IsGasCoinType(nemoConfig.CoinType) {
		_, gasCoin, err = api.RemainCoinAndGas(client.SuiApi, sender.Address, uint64(10000000), constant.GASCOINTYPE)
		if err != nil {
			return false, err
		}
	}

	mergeCoinArgument, remainingCoins, err := api.MergeCoin(ptb, client.SuiApi, remainingCoins, netSyIn)
	if err != nil {
		return false, err
	}

	splitResult, _, err := api.SplitCoinFromMerged(ptb, *mergeCoinArgument[0], netSyIn)
	if err != nil {
		return false, err
	}

	depositArgument, err := api.Deposit(ptb, client.SuiApi, nemoConfig, &splitResult)
	if err != nil {
		return false, err
	}

	oracleArgument, err := api.GetPriceVoucher(ptb, client.SuiApi, nemoConfig)
	if err != nil {
		return false, err
	}

	if exactAmountOutType == constant.YTTYPE {
		_, err = api.SwapExactSyForYt(ptb, client.BlockApi, client.SuiApi, nemoConfig, sender.Address, approxPyOut, netSyTokenization, minYtOut, oracleArgument, depositArgument)
		if err != nil {
			return false, err
		}
	} else if exactAmountOutType == constant.PTTYPE {
		_, err = api.SwapExactSyForPt(ptb, client.BlockApi, client.SuiApi, nemoConfig, sender.Address, approxPyOut, minYtOut, oracleArgument, depositArgument)
		if err != nil {
			return false, err
		}
	} else {
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

	b, _ := json.Marshal(resp)
	fmt.Printf("\n==resp:%+v==\n", string(b))

	return false, nil
}

func (s *SuiService) SwapCetusB2A(amountIn float64, sender *account.Account) (bool, error) {
	netAmountIn := uint64(amountIn * 1000000000)
	// create trade builder
	ptb := sui_types.NewProgrammableTransactionBuilder()
	client := InitSuiService()

	remainingCoins, gasCoin, err := api.RemainCoinAndGas(client.SuiApi, sender.Address, uint64(10000000), "0x2::sui::SUI")
	if err != nil {
		return false, err
	}

	if constant.IsGasCoinType("0x2::sui::SUI") {
		_, gasCoin, err = api.RemainCoinAndGas(client.SuiApi, sender.Address, uint64(10000000), constant.GASCOINTYPE)
		if err != nil {
			return false, err
		}
	}

	mergeCoinArgument, remainingCoins, err := api.MergeCoin(ptb, client.SuiApi, remainingCoins, netAmountIn)
	if err != nil {
		return false, err
	}

	splitResult, _, err := api.SplitCoinFromMerged(ptb, *mergeCoinArgument[0], netAmountIn)
	if err != nil {
		return false, err
	}

	// Get object arguments for the swap
	globalConfigArg, err := api.GetObjectArgument(ptb, client.SuiApi, "0xdaa46292632c3c4d8f31f23ea0f9b36a28ff3677e9684980e4438403a67a3d8f", false, "0x1eabed72c53feb3805120a081dc15963c204dc8d091542592abaf7a35689b2fb", "config", "swap_b2a")
	if err != nil {
		return false, err
	}

	poolArg, err := api.GetObjectArgument(ptb, client.SuiApi, "0x6c545e78638c8c1db7a48b282bb8ca79da107993fcb185f75cedc1f5adb2f535", true, "0x1eabed72c53feb3805120a081dc15963c204dc8d091542592abaf7a35689b2fb", "pool", "swap_b2a")
	if err != nil {
		return false, err
	}

	partnerArg, err := api.GetObjectArgument(ptb, client.SuiApi, "0xeb863165a109f7791a3182be08aff1438ab2a429314fc135ae19d953afe1edd6", true, "0x1eabed72c53feb3805120a081dc15963c204dc8d091542592abaf7a35689b2fb", "partner", "swap_b2a")
	if err != nil {
		return false, err
	}

	clockArg, err := api.GetObjectArgument(ptb, client.SuiApi, "0x0000000000000000000000000000000000000000000000000000000000000006", false, "0x2", "clock", "swap_b2a")
	if err != nil {
		return false, err
	}

	// Convert contract address to ObjectID
	contractAddr, err := sui_types.NewObjectIdFromHex("0x2485feb9d42c7c3bcb8ecde555ad40f1b073d9fb4faf354fa2d30a0b183a23ce")
	if err != nil {
		return false, err
	}

	// Create type arguments
	usdcStructTag, err := api.GetStructTag("0xdba34672e30cb065b1f93e3ab55318768fd6fef66c15942c9f7cb846e2f900e7::usdc::USDC")
	if err != nil {
		return false, err
	}
	usdcTypeTag := move_types.TypeTag{
		Struct: usdcStructTag,
	}

	suiStructTag, err := api.GetStructTag("0x2::sui::SUI")
	if err != nil {
		return false, err
	}
	suiTypeTag := move_types.TypeTag{
		Struct: suiStructTag,
	}

	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, usdcTypeTag) // T0
	typeArguments = append(typeArguments, suiTypeTag)  // T1

	// Call swap_b2a
	swapResult := ptb.Command(sui_types.Command{
		MoveCall: &sui_types.ProgrammableMoveCall{
			Package:       *contractAddr,
			Module:        "cetus",
			Function:      "swap_b2a",
			TypeArguments: typeArguments,
			Arguments: []sui_types.Argument{
				globalConfigArg,
				poolArg,
				partnerArg,
				splitResult,
				clockArg,
			},
		},
	})

	// Transfer the result to sender
	recipientAddr, err := sui_types.NewAddressFromHex(sender.Address)
	if err != nil {
		return false, err
	}

	recArg, err := ptb.Pure(*recipientAddr)
	if err != nil {
		return false, err
	}

	transferArgs := []sui_types.Argument{swapResult}
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

	b, _ := json.Marshal(resp)
	fmt.Printf("\n==resp:%+v==\n", string(b))

	return true, nil
}
