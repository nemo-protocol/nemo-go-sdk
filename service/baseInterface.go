package service

import (
	"github.com/coming-chat/go-sui/v2/account"
	"nemo-go-sdk/service/sui/common/models"
)

type ContractInterface interface {
	MintPy(amountIn float64, sender *account.Account, nemoConfig *models.NemoConfig)(bool, error)
	RedeemPy(amountIn float64, sender *account.Account, nemoConfig *models.NemoConfig)(bool, error)
	AddLiquidity(amountFloat float64, sender *account.Account, nemoConfig *models.NemoConfig)(bool, error)
	RedeemLiquidity(expectOut float64, sender *account.Account, nemoConfig *models.NemoConfig)(bool, error)
	SwapByPy(amountIn, slippage float64, amountInType, exactAmountOutType string, sender *account.Account, nemoConfig *models.NemoConfig)(bool, error)
	SwapToPy(amountIn, slippage float64, amountInType, exactAmountOutType string, sender *account.Account, nemoConfig *models.NemoConfig)(bool, error)
}
