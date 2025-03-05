package api

import (
	"fmt"
	"math"
	"time"

	"github.com/shopspring/decimal"
)

type CoinInfo struct {
	CoinPrice          string
	UnderlyingPrice    string
	UnderlyingApy      string
	SwapFeeForLpHolder string
	Decimal            int
	Maturity           int64
}

type MarketState struct {
	TotalPt       string
	TotalSy       string
	LpSupply      string
	RewardMetrics []RewardMetric
}

type RewardMetric struct {
	TokenPrice    string
	DailyEmission string
}

func safeDivide(numerator, denominator decimal.Decimal) decimal.Decimal {
	if denominator.IsZero() {
		return decimal.Zero
	}
	return numerator.Div(denominator)
}

func calculatePtApy(underlyingPrice, ptPrice, daysToExpiry decimal.Decimal) decimal.Decimal{
	fmt.Printf("\n==underlyingPrice:%v, ptPrice:%v, daysToExpiry:%v==\n",underlyingPrice, ptPrice, daysToExpiry)
	ratio := underlyingPrice.Div(ptPrice)
	exponent := decimal.NewFromFloat(365).Div(daysToExpiry)
	rf,_ := ratio.Float64()
	ef,_ := exponent.Float64()
	return decimal.NewFromFloat(math.Pow(rf,ef)).Sub(decimal.NewFromInt(1)).Mul(decimal.NewFromInt(100))
}

func CalculatePoolApy(coinInfo CoinInfo, marketState MarketState, ytIn, syOut int64) string {
	// Convert strings to decimals
	coinPrice, _ := decimal.NewFromString(coinInfo.CoinPrice)
	underlyingPrice, _ := decimal.NewFromString(coinInfo.UnderlyingPrice)
	underlyingApy, _ := decimal.NewFromString(coinInfo.UnderlyingApy)
	swapFeeForLpHolder, _ := decimal.NewFromString(coinInfo.SwapFeeForLpHolder)
	totalPt, _ := decimal.NewFromString(marketState.TotalPt)
	totalSy, _ := decimal.NewFromString(marketState.TotalSy)

	// Calculate days to expiry
	daysToExpiry := decimal.NewFromFloat(float64(coinInfo.Maturity/1000 - time.Now().Unix()) / float64(86400))

	// Calculate TVL
	ytPrice := safeDivide(coinPrice.Mul(decimal.NewFromInt(syOut)), decimal.NewFromInt(ytIn))
	ptPrice := underlyingPrice.Sub(ytPrice)
	fmt.Printf("\nsyOut:%v,ytIn:%v,ptPrice:%v\n",syOut,ytIn,ptPrice)
	ptTvl := totalPt.Mul(ptPrice).Div(decimal.NewFromInt(int64(math.Pow(10, float64(coinInfo.Decimal)))))
	syTvl := totalSy.Mul(coinPrice).Div(decimal.NewFromInt(int64(math.Pow(10, float64(coinInfo.Decimal)))))
	tvl := syTvl.Add(ptTvl)

	// Calculate scaled APYs
	rSy := safeDivide(totalSy, totalSy.Add(totalPt))
	rPt := safeDivide(totalPt, totalSy.Add(totalPt))
	ptApy := calculatePtApy(underlyingPrice, ptPrice, daysToExpiry)
	scaledUnderlyingApy := rSy.Mul(underlyingApy).Mul(decimal.NewFromInt(100))
	scaledPtApy := rPt.Mul(ptApy)
	fmt.Printf("\n==scaledPtApy:%v,ptApy:%v，scaledUnderlyingApy：%v==\n", scaledPtApy,ptApy,scaledUnderlyingApy)

	// Calculate swap fee APY
	swapFeeRateForLpHolder,_ := safeDivide(swapFeeForLpHolder.Mul(coinPrice), tvl).Float64()
	expiryRate,_ := safeDivide(decimal.NewFromFloat(365), daysToExpiry).Float64()
	swapFeeApy := decimal.NewFromFloat(math.Pow(swapFeeRateForLpHolder + 1,expiryRate)).Sub(decimal.NewFromInt(1)).Mul(decimal.NewFromInt(100))
	fmt.Printf("\nswapFeeApy:%v\n",swapFeeApy)

	// Calculate incentive APY
	incentiveApy := decimal.Zero
	for _, reward := range marketState.RewardMetrics {
		tokenPrice, _ := decimal.NewFromString(reward.TokenPrice)
		dailyEmission, _ := decimal.NewFromString(reward.DailyEmission)
		apy := safeDivide(tokenPrice.Mul(dailyEmission), tvl).Add(decimal.NewFromInt(1)).Pow(decimal.NewFromInt(365)).Sub(decimal.NewFromInt(1)).Mul(decimal.NewFromInt(100))
		incentiveApy = incentiveApy.Add(apy)
	}

	// Calculate pool APY
	poolApy := scaledUnderlyingApy.Add(scaledPtApy).Add(swapFeeApy).Add(incentiveApy)

	return poolApy.String()
}
