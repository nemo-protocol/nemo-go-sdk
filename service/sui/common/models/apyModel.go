package models

type ApyModel struct {
	PoolApy             string       `json:"poolApy"`
	PtApy               string       `json:"ptApy"`
	YtApy               string       `json:"ytApy"`
	Incentive           string       `json:"incentive"`
	ScaledUnderlyingApy string       `json:"scaledUnderlyingApy"`
	UnderlyingApy       string       `json:"underlyingApy"`
	ScaledApy           string       `json:"scaledApy"`
	ScaledPtApy         string       `json:"scaledPtApy"`
	Tvl                 string       `json:"tvl"`
	PtTvl               string       `json:"ptTvl"`
	YtTvl               string       `json:"ytTvl"`
	SyTvl               string       `json:"syTvl"`
	PtPrice             string       `json:"ptPrice"`
	YtPrice             string       `json:"ytPrice"`
	SwapFeeApy          string       `json:"swapFeeApy"`
	LpPrice             string       `json:"lpPrice"`
	MarketState         MarketState  `json:"marketState"`
	Incentives          []Incentives `json:"incentives"`
}

type Incentives struct {
	Apy       string `json:"apy"`
	TokenType string `json:"tokenType"`
	TokenLogo string `json:"tokenLogo"`
}

type MarketState struct {
	MarketCap     string          `json:"marketCap"`
	TotalSy       string          `json:"totalSy"`
	LpSupply      string          `json:"lpSupply"`
	TotalPt       string          `json:"totalPt"`
	RewardMetrics []RewardMetrics `json:"rewardMetrics"`
}

type RewardMetrics struct {
	TokenType     string `json:"tokenType"`
	TokenLogo     string `json:"tokenLogo"`
	DailyEmission string `json:"dailyEmission"`
	TokenPrice    string `json:"tokenPrice"`
	TokenName     string `json:"tokenName"`
	Decimal       string `json:"decimal"`
}
