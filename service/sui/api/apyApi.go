package api

import (
	"fmt"
	"github.com/nemo-protocol/nemo-go-sdk/service/sui/common/models"
	"math"
	"time"

	"github.com/shopspring/decimal"
)

type CoinInfo struct {
	CoinPrice          float64
	UnderlyingPrice    float64
	UnderlyingApy      float64
	SwapFeeForLpHolder float64
	Decimal            uint64
	Maturity           int64
}

type MarketState struct {
	MarketCap     string
	TotalPt       string
	TotalSy       string
	LpSupply      string
	RewardMetrics []RewardMetric
}

type RewardMetric struct {
	TokenPrice    string
	TokenLogo     string
	DailyEmission string
	CoinType      string
	CoinName      string
	Decimal       string
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

func CalculateYtAPY(underlyingInterestApyDec, ytPriceInAsset, yearsToExpiryDec decimal.Decimal) decimal.Decimal {
	yearsToExpiryResult,_ := yearsToExpiryDec.Float64()
	if yearsToExpiryResult <= 0 {
		return decimal.NewFromInt(0)
	}

	addResult,_ := underlyingInterestApyDec.Add(decimal.NewFromInt(1)).Float64()
	yearsToExpiryDecResult,_ := yearsToExpiryDec.Float64()

	interestReturns := decimal.NewFromFloat(math.Pow(addResult, yearsToExpiryDecResult)).
		Sub(decimal.NewFromInt(1))

	rewardsReturns := decimal.NewFromInt(0)

	ytReturns := interestReturns.Add(rewardsReturns)

	ytReturnsAfterFee := ytReturns.Mul(decimal.NewFromFloat(0.965))

	divResult1,_ := safeDivide(ytReturnsAfterFee, ytPriceInAsset).Float64()
	divResult2,_ := decimal.NewFromInt(1).Div(yearsToExpiryDec).Float64()

	ytApy := decimal.NewFromFloat(math.Pow(divResult1, divResult2)).
		Sub(decimal.NewFromInt(1)).
		Mul(decimal.NewFromInt(100))

	return ytApy
}

func CalculatePoolApy(coinInfo CoinInfo, marketState MarketState, ytIn, syOut int64) *models.ApyModel{
	response := &models.ApyModel{}
	// Convert strings to decimals
	coinPrice := decimal.NewFromFloat(coinInfo.CoinPrice)
	underlyingPrice := decimal.NewFromFloat(coinInfo.UnderlyingPrice)
	underlyingApy := decimal.NewFromFloat(coinInfo.UnderlyingApy)
	swapFeeForLpHolder := decimal.NewFromFloat(coinInfo.SwapFeeForLpHolder)
	totalPt, _ := decimal.NewFromString(marketState.TotalPt)
	totalSy, _ := decimal.NewFromString(marketState.TotalSy)
	lpSupply, _ := decimal.NewFromString(marketState.LpSupply)

	// Calculate days to expiry
	daysToExpiry := decimal.NewFromFloat(float64(coinInfo.Maturity/1000 - time.Now().Unix()) / float64(86400))
	yearToExpiry := decimal.NewFromFloat(float64(coinInfo.Maturity/1000 - time.Now().Unix()) / float64(31536000))

	// Calculate TVL
	ytPrice := safeDivide(coinPrice.Mul(decimal.NewFromInt(syOut)), decimal.NewFromInt(ytIn))
	ptPrice := underlyingPrice.Sub(ytPrice)
	fmt.Printf("\nsyOut:%v,ytIn:%v,ptPrice:%v\n",syOut,ytIn,ptPrice)
	ptTvl := totalPt.Mul(ptPrice).Div(decimal.NewFromInt(int64(math.Pow(10, float64(coinInfo.Decimal)))))
	ytTvl := totalPt.Mul(ytPrice).Div(decimal.NewFromInt(int64(math.Pow(10, float64(coinInfo.Decimal)))))
	syTvl := totalSy.Mul(coinPrice).Div(decimal.NewFromInt(int64(math.Pow(10, float64(coinInfo.Decimal)))))
	tvl := syTvl.Add(ptTvl)
	lpPrice := tvl.Div(lpSupply).Mul(decimal.NewFromInt(int64(math.Pow(10, float64(coinInfo.Decimal)))))

	// Calculate scaled APYs
	rSy := safeDivide(totalSy, totalSy.Add(totalPt))
	rPt := safeDivide(totalPt, totalSy.Add(totalPt))
	ptApy := calculatePtApy(underlyingPrice, ptPrice, daysToExpiry)
	fmt.Printf("yearToExpiry:%v\n",yearToExpiry)
	ytApy := CalculateYtAPY(underlyingApy, safeDivide(ytPrice, underlyingPrice), yearToExpiry)
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
	marketStateInfo := models.MarketState{
		MarketCap: marketState.MarketCap,
		TotalSy: marketState.TotalSy,
		TotalPt: marketState.TotalPt,
		LpSupply: marketState.LpSupply,
		RewardMetrics: make([]models.RewardMetrics, 0),
	}
	for _, reward := range marketState.RewardMetrics {
		tokenPrice, _ := decimal.NewFromString(reward.TokenPrice)
		dailyEmission, _ := decimal.NewFromString(reward.DailyEmission)
		divResult,_ := safeDivide(tokenPrice.Mul(dailyEmission), tvl).Add(decimal.NewFromInt(1)).Float64()
		apy := decimal.NewFromFloat(math.Pow(divResult, 365)).Sub(decimal.NewFromInt(1)).Mul(decimal.NewFromInt(100))
		incentiveApy = incentiveApy.Add(apy)

		incentives := models.Incentives{
			Apy: apy.String(),
			TokenType: reward.CoinType,
			TokenLogo: reward.TokenLogo,
		}
		response.Incentives = append(response.Incentives, incentives)

		rm := models.RewardMetrics{
			TokenType: reward.CoinType,
			TokenLogo: reward.TokenLogo,
			DailyEmission: dailyEmission.String(),
			TokenPrice: tokenPrice.String(),
			TokenName: reward.CoinName,
			Decimal: reward.Decimal,
		}
		marketStateInfo.RewardMetrics = append(marketStateInfo.RewardMetrics, rm)
	}

	// Calculate pool APY
	poolApy := scaledUnderlyingApy.Add(scaledPtApy).Add(swapFeeApy).Add(incentiveApy)

	response.PoolApy = poolApy.String()
	response.PtApy = ptApy.String()
	response.YtApy = ytApy.String()
	response.Incentive = incentiveApy.String()
	response.ScaledUnderlyingApy = scaledUnderlyingApy.String()
	response.ScaledPtApy = scaledPtApy.String()
	response.ScaledApy = scaledPtApy.String()
	response.Tvl = tvl.String()
	response.PtTvl = ptTvl.String()
	response.YtTvl = ytTvl.String()
	response.SyTvl = syTvl.String()
	response.PtPrice = ptPrice.String()
	response.YtPrice = ytPrice.String()
	response.SwapFeeApy = swapFeeApy.String()
	response.LpPrice = lpPrice.String()
	response.MarketState = marketStateInfo

	return response
}
