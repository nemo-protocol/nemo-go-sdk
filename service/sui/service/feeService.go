package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/nemo-protocol/nemo-go-sdk/service/sui/api"
	"github.com/nemo-protocol/nemo-go-sdk/service/sui/common/models"
	"math"
	"strconv"
	"strings"
	"time"
)

func (s *SuiService)QueryFee(nemoConfig *models.NemoConfig) (*models.FeeModel, error){
	yieldFactory,err := api.GetObjectFieldByObjectId(s.SuiApi, nemoConfig.YieldFactoryConfig)
	if err != nil{
		return nil, err
	}
	vaultBagId := yieldFactory["vault_bag"].(map[string]interface{})["fields"].(map[string]interface{})["bag"].(map[string]interface{})["fields"].(map[string]interface{})["id"].(map[string]interface{})["id"].(string)
	fmt.Printf("\n==vaultBagId:%v==\n",vaultBagId)

	parentId,err := sui_types.NewObjectIdFromHex(vaultBagId)
	if err != nil{
		return nil, err
	}

	var allData []interface{}
	var cursor *sui_types.ObjectID
	hasNextPage := true
	limit := uint(100)
	syCoinType := nemoConfig.SyCoinType

	for hasNextPage {
		respData, err := s.SuiApi.GetDynamicFields(context.Background(), *parentId, cursor, &limit)
		if err != nil {
			fmt.Printf("Error getting dynamic fields: %v\n", err)
			return nil, err
		}

		var dataMap map[string]interface{}
		jsonData, err := json.Marshal(respData)
		if err != nil {
			return nil, fmt.Errorf("marshal error: %v", err)
		}

		if err := json.Unmarshal(jsonData, &dataMap); err != nil {
			return nil, fmt.Errorf("unmarshal error: %v", err)
		}

		if data, ok := dataMap["data"].([]interface{}); ok {
			allData = append(allData, data...)
		}

		hasNextPage = false
		if nextCursor, ok := dataMap["nextCursor"].(string); ok && nextCursor != "" {
			cursor,err = sui_types.NewObjectIdFromHex(nextCursor)
			if err != nil{
				return nil, err
			}
			if hasNextPageValue, ok := dataMap["hasNextPage"].(bool); ok && hasNextPageValue {
				hasNextPage = true
			}
		}

		time.Sleep(time.Millisecond * 100)
	}

	syDynamicId := ""
	for _,data := range allData{
		objectType := data.(map[string]interface{})["objectType"].(string)
		objectId := data.(map[string]interface{})["objectId"].(string)
		if strings.Contains(objectType, syCoinType){
			syDynamicId = objectId
			break
		}
	}

	syDynamicFields,err := api.GetObjectFieldByObjectId(s.SuiApi, syDynamicId)
	if err != nil{
		return nil, err
	}

	feeValue := syDynamicFields["value"].(string)
	feeValueFloat,err := strconv.ParseFloat(feeValue, 64)
	if err != nil{
		return nil, err
	}

	feeModel := &models.FeeModel{
		YtClaimFee: feeValueFloat / math.Pow(10, float64(nemoConfig.Decimal)),
	}

	return feeModel, err
}
