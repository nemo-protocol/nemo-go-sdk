package api

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/coming-chat/go-sui/v2/account"
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/move_types"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/fardream/go-bcs/bcs"
	"nemo-go-sdk/service/sui/common/constant"
)

func DryRunGetApproxYtOutForNetSyInInternal(client *client.Client, nemoPackage, syType, pyState, marketGlobalConfig, marketState string, netSyIn, minYtOut uint64, sender *account.Account) (approxYtOut uint64, netSyTokenization uint64, err error) {
	ptb := sui_types.NewProgrammableTransactionBuilder()

	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoPackage)
	if err != nil {
		return 0, 0, err
	}

	syStructTag, err := GetStructTag(syType)
	if err != nil {
		return 0, 0, err
	}

	moduleName := "offchain"
	functionName := "get_approx_yt_out_for_net_sy_in_internal"
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

	minYtOutArg := CreatePureU64CallArg(minYtOut)
	minYtOutArgument, err := ptb.Input(minYtOutArg)
	if err != nil {
		return 0, 0, err
	}

	oracleArgument, err := GetPriceVoucherFromXOracle(ptb, client, nemoPackage, syType, constant.GASCOINTYPE)
	if err != nil{
		return 0, 0, err
	}

	ps, err := GetObjectArgument(ptb, client, pyState, false, nemoPackage, moduleName, functionName)
	if err != nil {
		return 0, 0, err
	}
	ms, err := GetObjectArgument(ptb, client, marketState, false, nemoPackage, moduleName, functionName)
	if err != nil {
		return 0, 0, err
	}
	mgc, err := GetObjectArgument(ptb, client, marketGlobalConfig, false, nemoPackage, moduleName, functionName)
	if err != nil {
		return 0, 0, err
	}
	c, err := GetObjectArgument(ptb, client, constant.CLOCK, false, nemoPackage, moduleName, functionName)
	if err != nil {
		return 0, 0, err
	}

	arguments := []sui_types.Argument{
		netSyInArgument,
		minYtOutArgument,
		*oracleArgument,
		ps,
		ms,
		mgc,
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
					approxYtOut = binary.LittleEndian.Uint64(byteSlice)
					fmt.Printf("Parsed approxYtOut: %d\n", approxYtOut)
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

	fmt.Printf("\nFinal values - approxYtOut: %d, netSyTokenization: %d\n", approxYtOut, netSyTokenization)

	return approxYtOut, netSyTokenization, nil
}

func GetYtOutForExactSyInWithPriceVoucher(client *client.Client, nemoPackage, syType, pyState, marketGlobalConfig, marketState string, netSyIn uint64, priceOracle string, sender *account.Account) (uint64, error){
	ptb := sui_types.NewProgrammableTransactionBuilder()

	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoPackage)
	if err != nil {
		return 0, err
	}

	syStructTag, err := GetStructTag(syType)
	if err != nil {
		return 0, err
	}

	moduleName := "router"
	functionName := "get_yt_out_for_exact_sy_in_with_price_voucher"
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

	minYtOut := uint64(0)
	minYtOutArg := CreatePureU64CallArg(minYtOut)
	minYtOutArgument, err := ptb.Input(minYtOutArg)
	if err != nil {
		return 0, err
	}

	oracleArgument, err := GetPriceVoucherFromXOracle(ptb, client, nemoPackage, syType, constant.GASCOINTYPE)
	if err != nil{
		return 0, err
	}

	ps, err := GetObjectArgument(ptb, client, pyState, false, nemoPackage, moduleName, functionName)
	if err != nil {
		return 0, err
	}
	ms, err := GetObjectArgument(ptb, client, marketState, false, nemoPackage, moduleName, functionName)
	if err != nil {
		return 0, err
	}
	mgc, err := GetObjectArgument(ptb, client, marketGlobalConfig, false, nemoPackage, moduleName, functionName)
	if err != nil {
		return 0, err
	}
	c, err := GetObjectArgument(ptb, client, constant.CLOCK, false, nemoPackage, moduleName, functionName)
	if err != nil {
		return 0, err
	}

	arguments := []sui_types.Argument{
		netSyInArgument,
		minYtOutArgument,
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
					minYtOut = binary.LittleEndian.Uint64(byteSlice)
					fmt.Printf("Parsed minYtOut: %d\n", minYtOut)
				}
			}
		}
	}
	return minYtOut, nil
}