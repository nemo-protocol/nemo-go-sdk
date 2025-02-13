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
	SYTYPE = "0x36a4c63cd17d48d33e32c3796245b2d1ebe50c2898ee80e682b787fb9b6519d5::sSUI::SSUI"
)

func (s *SuiService)MintPy(sourceCoin string, amountFloat float64, sender *account.Account) (bool, error){
	// create trade builder
	ptb := sui_types.NewProgrammableTransactionBuilder()
	client := InitSuiService()

	arg1, err := api.InitPyPosition(ptb, client.suiApi, NEMOPACKAGE, SYTYPE)
	if err != nil{
		return false, err
	}

	amountIn := uint64(amountFloat * 10000000)
	arg2, coins, err := api.SplitCoin(ptb, client.suiApi, amountIn, sender.Address, uint64(100000000))
	if err != nil{
		return false, err
	}

	// 转换接收地址
	recipientAddr, err := sui_types.NewAddressFromHex("0x1cbee4287aa17ff4d3ecfc2b5ee3e4cca27ec356040518c161a34101dd9ed491")
	if err != nil {
		return false, err
	}

	// 准备转移参数
	recArg, err := ptb.Pure(*recipientAddr)
	if err != nil {
		return false, err
	}

	// 转移两个对象
	transferArgs := []sui_types.Argument{*arg1, *arg2[0]}

	// 执行转移
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

	gasPayment := []*sui_types.ObjectRef{coins}

	// 转换发送者地址
	senderAddr, err := sui_types.NewObjectIdFromHex(sender.Address)
	if err != nil {
		return false, fmt.Errorf("failed to convert sender address: %w", err)
	}

	// 创建交易
	tx := sui_types.NewProgrammable(
		*senderAddr,
		gasPayment,
		pt,
		10000000, // gasBudget
		1000,     // gasPrice
	)


	// 序列化交易
	txBytes, err := bcs.Marshal(tx)
	if err != nil {
		return false, fmt.Errorf("failed to serialize transaction: %w", err)
	}

	// 签名交易
	signature, err := sender.SignSecureWithoutEncode(txBytes, sui_types.DefaultIntent())
	if err != nil {
		return false, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// 设置交易选项
	options := types.SuiTransactionBlockResponseOptions{
		ShowInput:          true,
		ShowEffects:        true,
		ShowEvents:         true,
		ShowObjectChanges:  true,
		ShowBalanceChanges: true,
	}

	// 执行交易
	resp, err := client.suiApi.ExecuteTransactionBlock(
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
	// 检查交易是否成功
	//if resp.Effects.Data {
	//	return false, fmt.Errorf("transaction failed: %s", resp.Effects.Status.Error)
	//}

	return false, nil
}

func (s *SuiService)RedeemPy(outCoin string, expectOut float64) (bool, error){
	return false, nil
}
