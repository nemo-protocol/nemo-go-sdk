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
	"nemo-go-sdk/service/sui/common/models"
)

func (s *SuiService)MintPy(amountIn float64, sender *account.Account, nemoConfig *models.NemoConfig) (bool, error){
	netSyIn := uint64(amountIn*1000000000)
	// create trade builder
	ptb := sui_types.NewProgrammableTransactionBuilder()
	client := InitSuiService()

	pyStateInfo, err := api.GetObjectFieldByObjectId(client.SuiApi, nemoConfig.PyState)
	if err != nil{
		return false, err
	}
	maturity := pyStateInfo["expiry"].(string)

	expectPyPositionTypeList := make([]string, 0)
	for _, pkg := range nemoConfig.NemoContractList{
		expectPyPositionTypeList = append(expectPyPositionTypeList, fmt.Sprintf("%v::py_position::PyPosition", pkg))
	}
	pyPosition,err := api.GetOwnerObjectByType(client.BlockApi, client.SuiApi, expectPyPositionTypeList, nemoConfig.SyCoinType, maturity, sender.Address)
	if err != nil {
		return false, err
	}

	var pyPositionArgument *sui_types.Argument
	if pyPosition == ""{
		pyPositionArgument, err = api.InitPyPosition(ptb, client.SuiApi, nemoConfig)
		if err != nil{
			return false, err
		}
	}else {
		argument, err := api.GetObjectArgument(ptb, client.SuiApi, pyPosition, false, nemoConfig.NemoContract, "yield_factory", "mint_py")
		if err != nil{
			return false, err
		}
		pyPositionArgument = &argument
	}

	remainingCoins, gasCoin, err := api.RemainCoinAndGas(client.SuiApi, sender.Address, uint64(10000000), nemoConfig.CoinType)
	if err != nil{
		return false, err
	}

	if !constant.IsGasCoinType(nemoConfig.CoinType){
		_, gasCoin, err = api.RemainCoinAndGas(client.SuiApi, sender.Address, uint64(10000000), constant.GASCOINTYPE)
		if err != nil{
			return false, err
		}
	}

	mergeCoinArgument, remainingCoins, err := api.MergeCoin(ptb, client.SuiApi, remainingCoins, netSyIn)
	if err != nil{
		return false, err
	}

	splitResult,_,err := api.SplitCoinFromMerged(ptb, *mergeCoinArgument[0], netSyIn)
	if err != nil{
		return false, err
	}

	depositArgument, err := api.Deposit(ptb, client.SuiApi, nemoConfig, &splitResult)

	oracleArgument, err := api.GetPriceVoucherFromXOracle(ptb, client.SuiApi, nemoConfig)
	if err != nil{
		return false, err
	}

	_, err = api.MintPy(ptb, client.SuiApi, nemoConfig, depositArgument, oracleArgument, pyPositionArgument)
	if err != nil{
		return false, err
	}
	//sCoin, err := api.MintSCoin(ptb, client.SuiApi, COINTYPE, UNDERLYINGCOINTYPE, arg2[0])
	//if err != nil{
	//	return false, err
	//}

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

	if pyPosition == ""{
		transferArgs := []sui_types.Argument{*pyPositionArgument}
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

	b,_ := json.Marshal(resp)
	fmt.Printf("\n==resp:%+v==\n",string(b))

	return true, nil
}

func (s *SuiService)RedeemPy(amountIn float64, sender *account.Account, nemoConfig *models.NemoConfig)(bool, error){
	netAmountIn := uint64(amountIn*1000000000)
	// create trade builder
	ptb := sui_types.NewProgrammableTransactionBuilder()
	client := InitSuiService()

	pyStateInfo, err := api.GetObjectFieldByObjectId(client.SuiApi, nemoConfig.PyState)
	if err != nil{
		return false, err
	}
	maturity := pyStateInfo["expiry"].(string)
	
	expectPyPositionTypeList := make([]string, 0)
	for _, pkg := range nemoConfig.NemoContractList{
		expectPyPositionTypeList = append(expectPyPositionTypeList, fmt.Sprintf("%v::py_position::PyPosition", pkg))
	}
	pyPosition,err := api.GetOwnerObjectByType(client.BlockApi, client.SuiApi, expectPyPositionTypeList, nemoConfig.SyCoinType, maturity, sender.Address)
	if err != nil {
		return false, err
	}
	if pyPosition == ""{
		return false, errors.New("pyPosition not exist！")
	}

	pyPositionArgument, err := api.GetObjectArgument(ptb, client.SuiApi, pyPosition, false, nemoConfig.NemoContract, "yield_factory", "redeem_py")
	if err != nil{
		return false, err
	}

	oracleArgument, err := api.GetPriceVoucherFromXOracle(ptb, client.SuiApi, nemoConfig)
	if err != nil{
		return false, err
	}

	_, gasCoin, err := api.RemainCoinAndGas(client.SuiApi, sender.Address, uint64(10000000), constant.GASCOINTYPE)
	if err != nil{
		return false, err
	}

	redeemPyResult, err := api.RedeemPy(ptb, client.SuiApi, nemoConfig, netAmountIn, oracleArgument, &pyPositionArgument)
	if err != nil{
		return false, err
	}

	syRedeemResult, err := api.SyRedeem(ptb, client.SuiApi, nemoConfig, redeemPyResult)
	if err != nil{
		return false, err
	}

	transferArgs := []sui_types.Argument{*syRedeemResult}

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

	b,_ := json.Marshal(resp)
	fmt.Printf("\n==resp:%+v==\n",string(b))

	return false, nil
}
