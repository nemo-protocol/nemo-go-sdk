package service

import (
	"github.com/nemo-protocol/nemo-go-sdk/service/sui/api"
	"github.com/nemo-protocol/nemo-go-sdk/service/sui/common/models"
	"math"
	"strconv"
)

func (s *SuiService)QueryAsset(nemoConfig *models.NemoConfig, address string) (*models.AssetModel, error){
	client := InitSuiService()
	assetModel := &models.AssetModel{}
	pyPositionList,err := api.GetPyPositionList(nemoConfig, address, client.SuiApi, client.BlockApi)
	if err != nil {
		return assetModel, err
	}

	var lpBalanceFloat, ytBalanceFloat, ptBalanceFloat float64
	decimalPow := math.Pow(10, float64(nemoConfig.Decimal))
	for _,pyPositionId := range pyPositionList{
		pyPosition, err := api.GetObjectFieldByObjectId(client.SuiApi, pyPositionId)
		if err != nil{
			return nil, err
		}
		ytBalanceStr := pyPosition["yt_balance"].(string)
		ptBalanceStr := pyPosition["pt_balance"].(string)
		ytBalance,_ := strconv.ParseFloat(ytBalanceStr, 64)
		ptBalance,_ := strconv.ParseFloat(ptBalanceStr, 64)
		if ytBalance > 0{
			ytBalanceFloat += ytBalance / decimalPow
		}
		if ptBalance > 0{
			ptBalanceFloat += ptBalance / decimalPow
		}
	}
	assetModel.YtBalance = ytBalanceFloat
	assetModel.PtBalance = ptBalanceFloat

	marketPositionList,err := api.GetMarketPositionList(client.BlockApi, client.SuiApi, nemoConfig, address)
	if err != nil {
		return assetModel, err
	}
	marketPositionFieldList,err :=  api.MultiGetObjectFieldByObjectId(client.SuiApi, marketPositionList)
	if err != nil {
		return assetModel, err
	}
	for _,marketPositionField := range marketPositionFieldList{
		marketPosition,err := api.SuiResponseToMap(marketPositionField)
		if err != nil{
			continue
		}
		lpBalanceStr := marketPosition["lp_amount"].(string)
		lpBalance,_ := strconv.ParseFloat(lpBalanceStr, 64)
		if lpBalance > 0{
			lpBalanceFloat += lpBalance / decimalPow
		}
	}
	assetModel.LpBalance = lpBalanceFloat

	return assetModel, nil
}
