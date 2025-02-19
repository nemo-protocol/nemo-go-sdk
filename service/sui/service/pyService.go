package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coming-chat/go-sui/v2/account"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
	"github.com/fardream/go-bcs/bcs"
	"nemo-go-sdk/service/sui/api"
)

var (
	NEMOPACKAGE = "0xa035d268323e40ab99ce8e4b12353bd89a63270935b4969d5bba87aa850c2b19"
	SYTYPE = "0x9b545bff00534f06d4f826802c2cc727c3827ac9a659ceeb117940b6c234dda7::sSCA::SSCA"
	COINTYPE = "0xaafc4f740de0dd0dde642a31148fb94517087052f19afb0f7bed1dc41a50c77b::scallop_sui::SCALLOP_SUI"
	UNDERLYINGCOINTYPE = "0x2::sui::SUI"
)

func (s *SuiService)MintPy(coinType string, amountFloat float64, sender *account.Account) (bool, error){
	// create trade builder
	ptb := sui_types.NewProgrammableTransactionBuilder()
	client := InitSuiService()

	arg1, err := api.InitPyPosition(ptb, client.SuiApi, NEMOPACKAGE, SYTYPE)
	if err != nil{
		return false, err
	}

	amountIn := uint64(amountFloat * 1000000000)
	remainingCoins, gasCoin, err := api.RemainCoinAndGas(client.SuiApi, sender.Address, uint64(10000000), coinType)
	if err != nil{
		return false, err
	}

	arg2, remainingCoins, err := api.MergeCoin(ptb, client.SuiApi, remainingCoins, amountIn)
	if err != nil{
		return false, err
	}


	sCoin, err := api.MintSCoin(ptb, client.SuiApi, COINTYPE, UNDERLYINGCOINTYPE, arg2[0])
	if err != nil{
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

	// transfer object
	transferArgs := []sui_types.Argument{*arg1, *sCoin}

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

	return true, nil
}

func (s *SuiService)RedeemPy(outCoin string, expectOut float64, sender *account.Account)(bool, error){
	return false, nil
}
