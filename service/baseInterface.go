package service

type ContractInterface interface {
	MintPy(sourceCoin string, expectIn float64)(bool, error)
	RedeemPy(outCoin string, expectOut float64)(bool, error)
	AddLiquidity(sourceCoin string, expectIn float64)(bool, error)
	RedeemLiquidity(outCoin string, expectOut float64)(bool, error)
}
