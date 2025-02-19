package service

import "github.com/coming-chat/go-sui/v2/account"

type ContractInterface interface {
	MintPy(coinType string, amountIn float64, sender *account.Account)(bool, error)
	RedeemPy(coinType string, amountIn float64, sender *account.Account)(bool, error)
	AddLiquidity(sourceCoin string, amountFloat float64, sender *account.Account)(bool, error)
	RedeemLiquidity(outCoin string, expectOut float64, sender *account.Account)(bool, error)
	SwapByPy(amountIn, slippage float64, coinType , amountInType, exactAmountOutType string, sender *account.Account)(bool, error)
	SwapToPy(amountIn, slippage float64, coinType , amountInType, exactAmountOutType string, sender *account.Account)(bool, error)
}
