package api

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/coming-chat/go-sui/v2/account"
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/move_types"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/fardream/go-bcs/bcs"
	"nemo-go-sdk/service/sui/common/constant"
	"nemo-go-sdk/service/sui/common/models"
)

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

func DryRunGetPyInForExactSyOutWithPriceVoucher(client *client.Client, nemoConfig *models.NemoConfig, exactPyType string, pyInAmount uint64, sender *account.Account) (uint64, error){
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