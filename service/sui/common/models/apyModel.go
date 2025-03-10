package models

type ApyModel struct {
	PoolApy             string                 `json:"poolApy"`
	PtApy               string                 `json:"ptApy"`
	YtApy               string                 `json:"ytApy"`
	IncentiveApy        string                 `json:"incentiveApy"`
	ScaledUnderlyingApy string                 `json:"scaledUnderlyingApy"`
	ScaledPtApy         string                 `json:"scaledPtApy"`
	Tvl                 string                 `json:"tvl"`
	PtTvl               string                 `json:"ptTvl"`
	YtTvl               string                 `json:"ytTvl"`
	SyTvl               string                 `json:"syTvl"`
	PtPrice             string                 `json:"ptPrice"`
	YtPrice             string                 `json:"ytPrice"`
	SwapFeeApy          string                 `json:"swapFeeApy"`
	LpPrice             string                 `json:"lpPrice"`
	MarketState         map[string]interface{} `json:"marketState"`
	Incentives          []Incentives           `json:"incentives"`
}

type Incentives struct {
	Apy       string `json:"apy"`
	TokenType string `json:"tokenType"`
	TokenLogo string `json:"tokenLogo"`
}
