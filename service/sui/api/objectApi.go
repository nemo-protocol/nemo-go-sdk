package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
	"nemo-go-sdk/service/sui/common/models"
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

func GetObjectMutable(sourceData *types.SuiObjectResponse) bool{
	marshal,_ := json.Marshal(sourceData.Data.Bcs.Data)
	bcsData := models.BcsData{}
	_ = json.Unmarshal(marshal, &bcsData)

	return bcsData.MoveObject.Version == sourceData.Data.Version.Int64() && bcsData.MoveObject.HasPublicTransfer
}