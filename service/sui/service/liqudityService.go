package service

import (
	"context"
	"encoding/json"
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

func (s *SuiService)AddLiquidity(amountFloat float64, sender *account.Account, amountInType string, nemoConfig *models.NemoConfig)(bool, error){
	amountSyIn := uint64(amountFloat * math.Pow(10, float64(nemoConfig.Decimal)))
	// create trade builder
	ptb := sui_types.NewProgrammableTransactionBuilder()
	client := InitSuiService()

	if amountInType == nemoConfig.UnderlyingCoinType{
		conversionRate,err := strconv.ParseFloat(nemoConfig.ConversionRate, 10)
		if err != nil{
			return false, err
		}
		amountSyIn = uint64(float64(amountSyIn) / conversionRate)
	}

	fmt.Printf("\n===amountSyIn:%v===\n",amountSyIn)
	minLpOut,err := api.DryRunGetLpOutForSingleSyIn(s.SuiApi, models.InitConfig(), amountSyIn, sender)
	if err != nil{
		return false, err
	}
	fmt.Printf("\n===minLpOut:%v===\n",minLpOut)

	ptValue,err := api.DryRunSingleLiquidityAddPtOut(s.SuiApi, models.InitConfig(), amountSyIn, sender)
	if err != nil{
		return false, err
	}

	remainingCoins, gasCoin, err := api.RemainCoinAndGas(client.SuiApi, sender.Address, uint64(30000000), amountInType)
	if err != nil{
		return false, err
	}

	splitResult ,_ ,err := api.SplitOrMergeCoin(ptb, client.SuiApi, remainingCoins, amountSyIn)
	if err != nil{
		return false, err
	}

	if !constant.IsGasCoinType(amountInType){
		_, gasCoin, err = api.RemainCoinAndGas(client.SuiApi, sender.Address, uint64(30000000), constant.GASCOINTYPE)
		if err != nil{
			return false, err
		}
	}else {
		argument,err := api.MintToSCoin(ptb, client.SuiApi, nemoConfig, &splitResult)
		if err != nil{
			return false, err
		}
		splitResult = *argument
	}

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

	oracleArgument, err := api.GetPriceVoucher(ptb, client.SuiApi, nemoConfig)
	if err != nil{
		return false, err
	}

	marketPosition, err := api.AddLiquiditySingleSy(ptb, client.SuiApi, nemoConfig, minLpOut, ptValue, oracleArgument, pyPositionArgument, depositArgument)
	if err != nil{
		return false, err
	}

	expectMarketPositionTypeList := make([]string, 0)
	for _, pkg := range nemoConfig.NemoContractList{
		expectMarketPositionTypeList = append(expectMarketPositionTypeList, fmt.Sprintf("%v::market_position::MarketPosition", pkg))
	}

	previousMarketPosition,err := api.GetOwnerMarketPositionByType(client.BlockApi, client.SuiApi, expectMarketPositionTypeList, nemoConfig.SyCoinType, maturity, sender.Address)
	if err != nil {
		return false, err
	}
	fmt.Printf("previousMarketPosition:%v\n",previousMarketPosition)

	if previousMarketPosition != "" {
		marketPosition,err = api.MergeAllLpPositions(ptb, client.SuiApi, nemoConfig, previousMarketPosition, marketPosition)
		if err != nil {
			return false, err
		}
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

	b,_ := json.Marshal(resp)
	fmt.Printf("\n==resp:%+v==\n",string(b))

	return true, nil
}

func (s *SuiService)RedeemLiquidity(expectOut float64, sender *account.Account, amountInType string, nemoConfig *models.NemoConfig)(bool, error){
	return false, nil
}
