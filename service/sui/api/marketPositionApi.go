package api

import (
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/move_types"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"nemo-go-sdk/service/sui/common/constant"
	"nemo-go-sdk/service/sui/common/models"
)

func MergeAllLpPositions(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, previousMarketPosition string, marketPositionArgument *sui_types.Argument) (*sui_types.Argument,error) {
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.NemoContract)
	if err != nil {
		return nil, err
	}

	moduleName := "market_position"
	functionName := "join"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)
	syStructTag, err := GetStructTag(nemoConfig.SyCoinType)
	if err != nil {
		return nil, err
	}
	syTypeTag := move_types.TypeTag{
		Struct: syStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, syTypeTag)

	previousMarketPositionArgument,err := GetObjectArgument(ptb, client, previousMarketPosition, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	clockArgument,err := GetObjectArgument(ptb, client, constant.CLOCK, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	var arguments []sui_types.Argument
	arguments = append(arguments, previousMarketPositionArgument, *marketPositionArgument, clockArgument)
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