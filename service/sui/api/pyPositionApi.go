package api

import (
	"errors"
	"fmt"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/move_types"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/nemo-protocol/nemo-go-sdk/service/sui/common/constant"
	"github.com/nemo-protocol/nemo-go-sdk/service/sui/common/models"
	"strings"
)

func InitPyPosition(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig) (*sui_types.Argument,error) {
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.NemoContract)
	if err != nil {
		return nil, err
	}

	module := move_types.Identifier("py")
	function := move_types.Identifier("init_py_position")
	structTag, err := GetStructTag(nemoConfig.SyCoinType)
	if err != nil {
		return nil, err
	}
	typeTag := move_types.TypeTag{
		Struct: structTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, typeTag)

	versionCallArg,err := GetObjectArg(client, nemoConfig.Version, false, nemoConfig.NemoContract, "py", "init_py_position")
	if err != nil {
		return nil, err
	}

	pyStateCallArg,err := GetObjectArg(client, nemoConfig.PyState, false, nemoConfig.NemoContract, "py", "init_py_position")
	if err != nil {
		return nil, err
	}

	clockCallArg,err := GetObjectArg(client, constant.CLOCK, false, nemoConfig.NemoContract, "py", "init_py_position")
	if err != nil {
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs, sui_types.CallArg{Object: versionCallArg}, sui_types.CallArg{Object: pyStateCallArg}, sui_types.CallArg{Object: clockCallArg})
	var arguments []sui_types.Argument
	for _, v := range callArgs {
		argument, err := ptb.Input(v)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, argument)
	}
	command := ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *nemoPackageId,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)
	return &command, nil
}

func GetObjectArgument(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, shareObject string, isCoin bool, contractPackage, module, function string) (sui_types.Argument, error){
	arg, err := GetObjectArg(client, shareObject, isCoin, contractPackage, module, function)
	if err != nil{
		return sui_types.Argument{}, err
	}
	return ptb.Input(sui_types.CallArg{Object: arg})
}

// shareObjectsMap: key->objectId; value isCoin boolean value
func MultiGetObjectArg(client *client.Client, shareObjectsMap map[string]bool, contractPackage, module, function string, cacheContractPackageInfo ...string) (map[string]*sui_types.ObjectArg, error) {
	if len(shareObjectsMap) == 0{
		return nil, errors.New("share Object map is null")
	}
	shareObjectIdList := make([]string, 0)
	for shareObject,_ := range shareObjectsMap{
		shareObjectIdList = append(shareObjectIdList, shareObject)
	}

	fmt.Printf("\n==shareObjectIdList:%+v==\n",shareObjectIdList)
	objectMap, err := MultiGetObjectFieldByObjectId(client, shareObjectIdList)
	if err != nil{
		return nil, err
	}

	objectArgMap := make(map[string]*sui_types.ObjectArg, 0)
	for objectId, sourceObjectData := range objectMap {
		hexObject, _ := sui_types.NewObjectIdFromHex(objectId)
		isCoin := shareObjectsMap[objectId]

		var objectArg *sui_types.ObjectArg
		if !isCoin && sourceObjectData.Data.Owner.AddressOwner == nil{
			objectArg = &sui_types.ObjectArg{
				SharedObject: &struct {
					Id                   sui_types.ObjectID
					InitialSharedVersion sui_types.SequenceNumber
					Mutable              bool
				}{
					Id:                   *hexObject,
					InitialSharedVersion: *sourceObjectData.Data.Owner.Shared.InitialSharedVersion,
					Mutable:              GetObjectMutable(client, *sourceObjectData.Data.Type, contractPackage, module, function, cacheContractPackageInfo...),
				},
			}
		}else {
			objectArg = &sui_types.ObjectArg{
				ImmOrOwnedObject: &sui_types.ObjectRef{
					ObjectId: sourceObjectData.Data.ObjectId,
					Version: sourceObjectData.Data.Version.Uint64(),
					Digest: sourceObjectData.Data.Digest,
				},
			}
		}
		objectArgMap[objectId] = objectArg
	}

	return objectArgMap, nil
}

func GetObjectArg(client *client.Client, shareObject string, isCoin bool, contractPackage, module, function string) (*sui_types.ObjectArg, error) {
	if shareObject == ""{
		return nil, errors.New("share Object is null")
	}
	hexObject, err := sui_types.NewObjectIdFromHex(shareObject)
	if err != nil {
		return nil, err
	}
	sourceObjectData, err := GetObjectMetadata(client, shareObject)
	if err != nil {
		return nil, err
	}
	if sourceObjectData.Data == nil{
		return nil, errors.New("get share Object fail")
	}
	if !isCoin && sourceObjectData.Data.Owner.AddressOwner == nil{
		return &sui_types.ObjectArg{
			SharedObject: &struct {
				Id                   sui_types.ObjectID
				InitialSharedVersion sui_types.SequenceNumber
				Mutable              bool
			}{
				Id:                   *hexObject,
				InitialSharedVersion: *sourceObjectData.Data.Owner.Shared.InitialSharedVersion,
				Mutable:              GetObjectMutable(client, *sourceObjectData.Data.Type, contractPackage, module, function),
			},
		}, nil
	}
	return &sui_types.ObjectArg{
		ImmOrOwnedObject: &sui_types.ObjectRef{
			ObjectId: sourceObjectData.Data.ObjectId,
			Version: sourceObjectData.Data.Version.Uint64(),
			Digest: sourceObjectData.Data.Digest,
		},
	}, nil
}

func GetStructTag(syType string) (*move_types.StructTag, error) {
	elements := strings.Split(syType, "::")
	if len(elements) != 3 {
		return nil, errors.New("syType invalid！")
	}

	addr, err := sui_types.NewAddressFromHex(elements[0])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("init syType address error: %v", err))
	}

	structTag := &move_types.StructTag{
		Address:    *addr,
		Module:     move_types.Identifier(elements[1]),
		Name:       move_types.Identifier(elements[2]),
		TypeParams: []move_types.TypeTag{},
	}

	return structTag, nil
}

func GetPyPosition(nemoConfig *models.NemoConfig, address string, client *client.Client, blockApi *sui.ISuiAPI) (string, error){
	pyStateInfo, err := GetObjectFieldByObjectId(client, nemoConfig.PyState)
	if err != nil{
		return "", err
	}
	maturity := pyStateInfo["expiry"].(string)

	expectPyPositionTypeList := make([]string, 0)
	for _, pkg := range nemoConfig.NemoContractList{
		expectPyPositionTypeList = append(expectPyPositionTypeList, fmt.Sprintf("%v::py_position::PyPosition", pkg))
	}

	pyPosition,err := GetOwnerObjectByType(blockApi, client, expectPyPositionTypeList, nemoConfig.SyCoinType, maturity, address)
	if err != nil {
		return "", err
	}
	return pyPosition, nil
}

func GetPyPositionList(nemoConfig *models.NemoConfig, address string, client *client.Client, blockApi *sui.ISuiAPI) ([]string, error){
	pyStateInfo, err := GetObjectFieldByObjectId(client, nemoConfig.PyState)
	if err != nil{
		return nil, err
	}
	maturity := pyStateInfo["expiry"].(string)

	expectPyPositionTypeList := make([]string, 0)
	for _, pkg := range nemoConfig.NemoContractList{
		expectPyPositionTypeList = append(expectPyPositionTypeList, fmt.Sprintf("%v::py_position::PyPosition", pkg))
	}

	pyPositionList,err := GetOwnerObjectListByType(blockApi, client, expectPyPositionTypeList, nemoConfig.SyCoinType, maturity, address)
	if err != nil {
		return pyPositionList, err
	}
	return pyPositionList, nil
}
