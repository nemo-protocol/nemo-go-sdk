package service

import (
	"github.com/coming-chat/go-sui/v2/account"
	"github.com/nemo-protocol/nemo-go-sdk/service/sui/api"
	"github.com/nemo-protocol/nemo-go-sdk/service/sui/common/models"
)

type ContractInterface interface {
	MintPy(amountIn float64, sender *account.Account, nemoConfig *models.NemoConfig)(bool, error)
	RedeemPy(amountIn float64, sender *account.Account, nemoConfig *models.NemoConfig)(bool, error)
	AddLiquidity(amountIn, slippage float64, sender *account.Account, amountInType string, nemoConfig *models.NemoConfig)(bool, error)
	RedeemLiquidity(amountOut, slippage float64, sender *account.Account, amountOutType string, nemoConfig *models.NemoConfig)(bool, error)
	SwapByPy(amountIn, slippage float64, amountInType, amountOutType string, sender *account.Account, nemoConfig *models.NemoConfig)(bool, error)
	SwapToPy(amountIn, slippage float64, amountInType, amountOutType string, sender *account.Account, nemoConfig *models.NemoConfig)(bool, error)
	ClaimYtReward(nemoConfig *models.NemoConfig, sender *account.Account) (bool, error)
	ClaimLpReward(nemoConfig *models.NemoConfig, sender *account.Account) (bool, error)
	QueryPoolApy(nemoConfig *models.NemoConfig, priceInfoMap ...map[string]api.PriceInfo) (*models.ApyModel, error)
}
