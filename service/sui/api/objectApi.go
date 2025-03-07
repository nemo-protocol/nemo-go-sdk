package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	blockModels "github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
	"github.com/nemo-protocol/nemo-go-sdk/service/sui/common/models"
	"github.com/nemo-protocol/nemo-go-sdk/utils"
	"strings"
)

func GetObjectMetadata(client *client.Client, objectId string) (*types.SuiObjectResponse, error) {
	hexId, err := sui_types.NewObjectIdFromHex(objectId)
	if err != nil {
		fmt.Println("objectId invalid!")
		return nil, err
	}
	object, err := client.GetObject(context.Background(), *hexId, &types.SuiObjectDataOptions{
		ShowType: true, ShowContent: true, ShowBcs: true, ShowOwner: true, ShowPreviousTransaction: true, ShowStorageRebate: true, ShowDisplay: true,
	})
	if err != nil {
		return nil, err
	}
	return object, nil
}

func GetObjectMutable(client *client.Client, objectType, contractPackage, module, function string) bool{
	contractPackageAddr,err := sui_types.NewObjectIdFromHex(contractPackage)
	if err != nil{
		return false
	}
	objects,_ := client.GetObject(context.Background(), *contractPackageAddr, &types.SuiObjectDataOptions{
		ShowType: true, ShowContent: true, ShowBcs: true, ShowOwner: true, ShowPreviousTransaction: true, ShowStorageRebate: true, ShowDisplay: true,
	})
	if objects == nil || objects.Data == nil || objects.Data.Content == nil || objects.Data.Content.Data == nil{
		return false
	}
	marshal, err := json.Marshal(objects.Data.Content.Data)
	if err != nil{
		return false
	}
	packageObject := models.Object{}
	_ = json.Unmarshal(marshal, &packageObject)

	filterFunc := utils.FindFunctionInBytecode(packageObject.Package.Disassembled[module].(string), function)
	args := strings.Split(filterFunc, ",")
	typeList := strings.SplitN(objectType, "::", 3)
	fmt.Printf("\n==typeList:%v, args:%v==\n",typeList,args)
	if len(typeList) != 3{
		return false
	}
	argElement := ""
	for _,v := range args{
		if strings.HasSuffix(v, typeList[2]){
			argElement = v
			break
		}else if strings.Contains(v, "<") && strings.Contains(typeList[2], "<"){
			name := strings.SplitN(typeList[2], "<", 2)[0]
			if strings.Contains(v, fmt.Sprintf(" %v<",name)){
				argElement = v
			}
		}
	}

	//have &mut
	if strings.Contains(argElement, "&mut"){
		return true
	}
	return false
}

func GetOwnObjectsMap(client *sui.ISuiAPI, ownerAddress string) ([]map[string]interface{}, error) {
	if client == nil{
		return nil, errors.New("client is nil pointer！")
	}
	limit := uint64(50)
	objectFilter := blockModels.SuiXGetOwnedObjectsRequest{
		Address: ownerAddress,
		Query: blockModels.SuiObjectResponseQuery{
			Options: blockModels.SuiObjectDataOptions{
				ShowType: true, ShowContent: true, ShowBcs: true, ShowOwner: true, ShowPreviousTransaction: true, ShowStorageRebate: true, ShowDisplay: true,
			},
		},
		Limit: limit,
	}
	var resp blockModels.PaginatedObjectsResponse
	respList := make([]map[string]interface{}, 0)
	var err error
	hasNext := true
	for hasNext{
		resp, err = (*client).SuiXGetOwnedObjects(context.Background(), objectFilter)
		if err != nil{
			return nil, err
		}

		if resp.HasNextPage{
			hasNext = resp.HasNextPage
			cursor := resp.NextCursor
			objectFilter.Cursor = cursor
		}else {
			hasNext = false
		}


		b, _ := json.Marshal(resp.Data)
		data := make([]map[string]interface{}, 0)
		_ = json.Unmarshal(b, &data)
		respList = append(respList, data...)
	}

	return respList, nil
}

func GetObjectFieldByObjectId(client *client.Client, objectId string) (map[string]interface{}, error){
	objectIdHex, err := sui_types.NewObjectIdFromHex(objectId)
	if err != nil{
		return nil, err
	}

	options := types.SuiObjectDataOptions{
		ShowType: true, ShowContent: true, ShowBcs: true, ShowOwner: true, ShowPreviousTransaction: true, ShowStorageRebate: true, ShowDisplay: true,
	}
	info, err := client.GetObject(context.Background(), *objectIdHex, &options)
	if err != nil{
		return nil, err
	}

	byteData,err := json.Marshal(info.Data.Content.Data)
	if err != nil{
		return nil, err
	}

	commonModels := models.CommonOnChainDataResp{}
	_ = json.Unmarshal(byteData, &commonModels)
	byteData, err = json.Marshal(commonModels.MoveObject.Fields)
	if err != nil{
		return nil, err
	}
	fields := make(map[string]interface{}, 0)
	err = json.Unmarshal(byteData, &fields)
	if err != nil{
		return nil, err
	}
	return fields, nil
}

func GetOwnerObjectByType(blockClient *sui.ISuiAPI, client *client.Client, objectsType []string, syType, maturity string, ownerAddress string) (string,error){
	objectsMap, err := GetOwnObjectsMap(blockClient, ownerAddress)
	if err != nil{
		return "", err
	}
	fmt.Printf("\n==objectsType list:%+v==\n",objectsType)
	for _,v := range objectsMap{
		objectType := v["data"].(map[string]interface{})["type"].(string)
		objectId := v["data"].(map[string]interface{})["objectId"].(string)
		fmt.Printf("\n==objectType:%v,objectId:%v==\n",objectType,objectId)
		if !utils.Contains(objectsType, objectType){
			continue
		}

		fields,err := GetObjectFieldByObjectId(client, objectId)
		if err != nil{
			continue
		}

		objectMaturity := fields["expiry"]
		pyStateId := fields["py_state_id"]
		if objectMaturity == "" || pyStateId == ""{
			continue
		}

		pyStateFields,err := GetObjectFieldByObjectId(client, pyStateId.(string))
		if err != nil{
			continue
		}
		objectSyType := pyStateFields["name"].(string)
		if !strings.HasPrefix(objectSyType, "0x"){
			objectSyType = fmt.Sprintf("0x%v",objectSyType)
		}

		if objectMaturity.(string) == maturity && objectSyType == syType{
			return objectId, nil
		}
	}
	return "", nil
}

func GetOwnerMarketPositionByType(blockClient *sui.ISuiAPI, client *client.Client, objectsType []string, syType, maturity string, ownerAddress string) (string,error){
	objectsMap, err := GetOwnObjectsMap(blockClient, ownerAddress)
	if err != nil{
		return "", err
	}
	fmt.Printf("\n==objectsType list:%+v==\n",objectsType)
	for _,v := range objectsMap{
		objectType := v["data"].(map[string]interface{})["type"].(string)
		objectId := v["data"].(map[string]interface{})["objectId"].(string)
		fmt.Printf("\n==objectType:%v,objectId:%v==\n",objectType,objectId)
		if !utils.Contains(objectsType, objectType){
			continue
		}

		fields,err := GetObjectFieldByObjectId(client, objectId)
		if err != nil{
			continue
		}

		objectMaturity := fields["expiry"]
		marketStateId := fields["market_state_id"]
		if objectMaturity == "" || marketStateId == ""{
			continue
		}

		marketStateFields,err := GetObjectFieldByObjectId(client, marketStateId.(string))
		if err != nil{
			continue
		}
		pyStateId := marketStateFields["py_state_id"]
		fmt.Printf("\npyStateId:%v\n",pyStateId)

		pyStateFields,err := GetObjectFieldByObjectId(client, pyStateId.(string))
		if err != nil{
			continue
		}
		objectSyType := pyStateFields["name"].(string)
		if !strings.HasPrefix(objectSyType, "0x"){
			objectSyType = fmt.Sprintf("0x%v",objectSyType)
		}

		if objectMaturity.(string) == maturity && objectSyType == syType{
			return objectId, nil
		}
	}
	return "", nil
}