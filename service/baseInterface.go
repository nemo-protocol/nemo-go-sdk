package service

import "github.com/coming-chat/go-sui/v2/account"

type ContractInterface interface {
	MintPy(sourceCoin string, expectIn float64, sender *account.Account)(bool, error)
	RedeemPy(outCoin string, expectOut float64, sender *account.Account)(bool, error)
	AddLiquidity(sourceCoin string, amountFloat float64, sender *account.Account)(bool, error)
	RedeemLiquidity(outCoin string, expectOut float64, sender *account.Account)(bool, error)
	SwapByPy(amountIn, slippage float64, amountInType, exactAmountOutType string, sender *account.Account)(bool, error)
	SwapToPy()(bool, error)
}
