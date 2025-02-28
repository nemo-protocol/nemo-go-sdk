package api

import (
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/move_types"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"nemo-go-sdk/service/sui/common/constant"
	"nemo-go-sdk/service/sui/common/models"
)

func MergeAllLpPositions(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, previousMarketPositionArgument, marketPositionArgument *sui_types.Argument) (*sui_types.Argument,error) {
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.NemoContract)
	if err != nil {
		return nil, err
	}

	moduleName := "market_position"
	functionName := "join"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	marketPositionNestResult := sui_types.Argument{
		NestedResult: &struct {
			Result1 uint16
			Result2 uint16
		}{
			Result1: *marketPositionArgument.Result,
			Result2: 0,  // 第一个分割结果
		},
	}

	clockArgument,err := GetObjectArgument(ptb, client, constant.CLOCK, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	var arguments []sui_types.Argument
	arguments = append(arguments, *previousMarketPositionArgument, marketPositionNestResult, clockArgument)
	command := ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *nemoPackageId,
				Module:        module,
				Function:      function,
				Arguments:     arguments,
			},
		},
	)
	return &command, nil
}