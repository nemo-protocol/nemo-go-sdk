package api

import (
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/move_types"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"nemo-go-sdk/service/sui/common/constant"
	"nemo-go-sdk/service/sui/common/models"
)

func GetPriceVoucherFromXOracle(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig) (*sui_types.Argument,error) {
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.NemoContract)
	if err != nil {
		return nil, err
	}

	module := move_types.Identifier("oracle")
	function := move_types.Identifier("get_price_voucher_from_x_oracle")
	syStructTag, err := GetStructTag(nemoConfig.SyCoinType)
	if err != nil {
		return nil, err
	}
	syTypeTag := move_types.TypeTag{
		Struct: syStructTag,
	}

	structTag, err := GetStructTag(nemoConfig.UnderlyingCoinType)
	if err != nil {
		return nil, err
	}
	typeTag := move_types.TypeTag{
		Struct: structTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, syTypeTag, typeTag)

	priceOracleCallArg,err := GetObjectArg(client, nemoConfig.PriceOracle, false, nemoConfig.NemoContract, "oracle", "get_price_voucher_from_x_oracle")
	if err != nil {
		return nil, err
	}

	scallopVersionCallArg,err := GetObjectArg(client, SCALLOPVERSION, false, nemoConfig.NemoContract, "oracle", "get_price_voucher_from_x_oracle")
	if err != nil {
		return nil, err
	}

	scallopMarketCallArg,err := GetObjectArg(client, SCALLOPMARKETOBJECT, false, nemoConfig.NemoContract, "oracle", "get_price_voucher_from_x_oracle")
	if err != nil {
		return nil, err
	}

	syStateCallArg,err := GetObjectArg(client, nemoConfig.SyState, false, nemoConfig.NemoContract, "oracle", "get_price_voucher_from_x_oracle")
	if err != nil {
		return nil, err
	}

	clockCallArg,err := GetObjectArg(client, constant.CLOCK, false, nemoConfig.NemoContract, "oracle", "get_price_voucher_from_x_oracle")
	if err != nil {
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs, sui_types.CallArg{Object: priceOracleCallArg}, sui_types.CallArg{Object: scallopVersionCallArg}, sui_types.CallArg{Object: scallopMarketCallArg}, sui_types.CallArg{Object: syStateCallArg}, sui_types.CallArg{Object: clockCallArg})
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
