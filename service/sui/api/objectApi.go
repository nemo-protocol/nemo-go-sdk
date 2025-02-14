package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
	"nemo-go-sdk/service/sui/common/models"
	"nemo-go-sdk/utils"
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
	po,_ := sui_types.NewObjectIdFromHex(contractPackage)
	o,_ := client.GetObject(context.Background(), *po, &types.SuiObjectDataOptions{
		ShowType: true, ShowContent: true, ShowBcs: true, ShowOwner: true, ShowPreviousTransaction: true, ShowStorageRebate: true, ShowDisplay: true,
	})
	marshal, _ := json.Marshal(o.Data.Content.Data)
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