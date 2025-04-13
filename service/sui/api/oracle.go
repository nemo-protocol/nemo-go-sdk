package api

import (
	"errors"
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/move_types"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/nemo-protocol/nemo-go-sdk/service/sui/common/constant"
	"github.com/nemo-protocol/nemo-go-sdk/service/sui/common/models"
)

func GetPriceVoucherFromXOracle(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig) (*sui_types.Argument,error) {
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.OraclePackage)
	if err != nil {
		return nil, err
	}

	moduleName := "scallop"
	functionName := "get_price_voucher_from_x_oracle"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)
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

	priceOracleCallArg,err := GetObjectArg(client, nemoConfig.PriceOracle, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	oracleTicketCallArg,err := GetObjectArg(client, nemoConfig.OracleTicket, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	scallopVersionCallArg,err := GetObjectArg(client, SCALLOP_VERSION, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	scallopMarketCallArg,err := GetObjectArg(client, SCALLOP_MARKET_OBJECT, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	syStateCallArg,err := GetObjectArg(client, nemoConfig.SyState, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	clockCallArg,err := GetObjectArg(client, constant.CLOCK, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs, sui_types.CallArg{Object: priceOracleCallArg}, sui_types.CallArg{Object: oracleTicketCallArg}, sui_types.CallArg{Object: scallopVersionCallArg}, sui_types.CallArg{Object: scallopMarketCallArg}, sui_types.CallArg{Object: syStateCallArg}, sui_types.CallArg{Object: clockCallArg})
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

func GetPriceVoucherFromVolo(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig) (*sui_types.Argument,error) {
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.OraclePackage)
	if err != nil {
		return nil, err
	}

	moduleName := "volo"
	functionName := "get_price_voucher_from_volo"
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

	priceOracleCallArg,err := GetObjectArg(client, nemoConfig.PriceOracle, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	oracleTicketCallArg,err := GetObjectArg(client, nemoConfig.OracleTicket, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	nativePoolCallArg,err := GetObjectArg(client, NATIVE_POOL, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	metadataCallArg,err := GetObjectArg(client, METADATA, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	syStateCallArg,err := GetObjectArg(client, nemoConfig.SyState, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs, sui_types.CallArg{Object: priceOracleCallArg}, sui_types.CallArg{Object: oracleTicketCallArg}, sui_types.CallArg{Object: nativePoolCallArg}, sui_types.CallArg{Object: metadataCallArg}, sui_types.CallArg{Object: syStateCallArg})
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

func GetPriceVoucherFromSpring(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, lstInfo string, moduleName string) (*sui_types.Argument,error) {
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.OraclePackage)
	if err != nil {
		return nil, err
	}

	functionName := "get_price_voucher_from_spring"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)
	syStructTag, err := GetStructTag(nemoConfig.SyCoinType)
	if err != nil {
		return nil, err
	}
	coinTypeStructTag, err := GetStructTag(nemoConfig.CoinType)
	if err != nil {
		return nil, err
	}
	syTypeTag := move_types.TypeTag{
		Struct: syStructTag,
	}
	coinTypeTypeTag := move_types.TypeTag{
		Struct: coinTypeStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, syTypeTag, coinTypeTypeTag)

	priceOracleCallArg,err := GetObjectArg(client, nemoConfig.PriceOracle, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	oracleTicketCallArg,err := GetObjectArg(client, nemoConfig.OracleTicket, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	lstInfoCallArg,err := GetObjectArg(client, lstInfo, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	syStateCallArg,err := GetObjectArg(client, nemoConfig.SyState, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs, sui_types.CallArg{Object: priceOracleCallArg}, sui_types.CallArg{Object: oracleTicketCallArg}, sui_types.CallArg{Object: lstInfoCallArg}, sui_types.CallArg{Object: syStateCallArg})
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

func GetPriceVoucherFromAftermath(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig) (*sui_types.Argument,error) {
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.OraclePackage)
	if err != nil {
		return nil, err
	}

	moduleName := "aftermath"
	functionName := "get_price_voucher_from_aftermath"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)
	syStructTag, err := GetStructTag(nemoConfig.SyCoinType)
	if err != nil {
		return nil, err
	}
	syTypeTag := move_types.TypeTag{
		Struct: syStructTag,
	}
	coinTypeStructTag, err := GetStructTag(nemoConfig.CoinType)
	if err != nil {
		return nil, err
	}
	coinTypeTag := move_types.TypeTag{
		Struct: coinTypeStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, syTypeTag, coinTypeTag)

	priceOracleCallArg,err := GetObjectArg(client, nemoConfig.PriceOracle, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	oracleTicketCallArg,err := GetObjectArg(client, nemoConfig.OracleTicket, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	stakeSuiValutCallArg,err := GetObjectArg(client, STAKED_SUI_VAULT, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	safeCallArg,err := GetObjectArg(client, SAFE, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	syStateCallArg,err := GetObjectArg(client, nemoConfig.SyState, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs, sui_types.CallArg{Object: priceOracleCallArg}, sui_types.CallArg{Object: oracleTicketCallArg}, sui_types.CallArg{Object: stakeSuiValutCallArg}, sui_types.CallArg{Object: safeCallArg}, sui_types.CallArg{Object: syStateCallArg})
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

func GetPriceVoucherFromHasui(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig) (*sui_types.Argument,error) {
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.OraclePackage)
	if err != nil {
		return nil, err
	}

	moduleName := "haedal"
	functionName := "get_price_voucher_from_haSui"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)
	syStructTag, err := GetStructTag(nemoConfig.SyCoinType)
	if err != nil {
		return nil, err
	}
	syTypeTag := move_types.TypeTag{
		Struct: syStructTag,
	}
	coinTypeStructTag, err := GetStructTag(nemoConfig.CoinType)
	if err != nil {
		return nil, err
	}
	coinTypeTag := move_types.TypeTag{
		Struct: coinTypeStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, syTypeTag, coinTypeTag)

	priceOracleCallArg,err := GetObjectArg(client, nemoConfig.PriceOracle, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	oracleTicketCallArg,err := GetObjectArg(client, nemoConfig.OracleTicket, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	haedalStakeingCallArg,err := GetObjectArg(client, HAEDAL_STAKING, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	syStateCallArg,err := GetObjectArg(client, nemoConfig.SyState, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs, sui_types.CallArg{Object: priceOracleCallArg}, sui_types.CallArg{Object: oracleTicketCallArg}, sui_types.CallArg{Object: haedalStakeingCallArg}, sui_types.CallArg{Object: syStateCallArg})
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

func GetPriceVoucherFromLpToken(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, lpVault, lpPool, moduleName string) (*sui_types.Argument,error) {
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.OraclePackage)
	if err != nil {
		return nil, err
	}

	functionName := "get_price_voucher_from_cetus_vault"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)
	syStructTag, err := GetStructTag(nemoConfig.SyCoinType)
	if err != nil {
		return nil, err
	}
	yieldTokenStructTag, err := GetStructTag(nemoConfig.YieldTokenType)
	if err != nil {
		return nil, err
	}
	coinTypeStructTag, err := GetStructTag(nemoConfig.CoinType)
	if err != nil {
		return nil, err
	}
	syTypeTag := move_types.TypeTag{
		Struct: syStructTag,
	}
	YieldTokenTypeTag := move_types.TypeTag{
		Struct: yieldTokenStructTag,
	}
	coinTypeTypeTag := move_types.TypeTag{
		Struct: coinTypeStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, syTypeTag, YieldTokenTypeTag, coinTypeTypeTag)

	priceOracleCallArg,err := GetObjectArg(client, nemoConfig.PriceOracle, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	oracleTicketCallArg,err := GetObjectArg(client, nemoConfig.OracleTicket, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	stakeingCallArg,err := GetObjectArg(client, HAEDAL_STAKING, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	vaultCallArg,err := GetObjectArg(client, lpVault, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	poolCallArg,err := GetObjectArg(client, lpPool, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	syStateCallArg,err := GetObjectArg(client, nemoConfig.SyState, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	//aftermath need
	stakeVaultCallArg,err := GetObjectArg(client, STAKED_SUI_VAULT, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	//aftermath need
	safeCallArg,err := GetObjectArg(client, SAFE, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	//vsui need
	nativePoolCallArg,err := GetObjectArg(client, NATIVE_POOL, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	//vsui need
	metadataCallArg,err := GetObjectArg(client, METADATA, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	//default hasui
	if constant.IsLpTokenAfSui(nemoConfig.CoinType){
		callArgs = append(callArgs, sui_types.CallArg{Object: priceOracleCallArg}, sui_types.CallArg{Object: oracleTicketCallArg}, sui_types.CallArg{Object: stakeVaultCallArg}, sui_types.CallArg{Object: safeCallArg}, sui_types.CallArg{Object: vaultCallArg}, sui_types.CallArg{Object: poolCallArg}, sui_types.CallArg{Object: syStateCallArg})
	} else if constant.IsLpTokenVSui(nemoConfig.CoinType){
		callArgs = append(callArgs, sui_types.CallArg{Object: priceOracleCallArg}, sui_types.CallArg{Object: oracleTicketCallArg}, sui_types.CallArg{Object: nativePoolCallArg}, sui_types.CallArg{Object: metadataCallArg}, sui_types.CallArg{Object: vaultCallArg}, sui_types.CallArg{Object: poolCallArg}, sui_types.CallArg{Object: syStateCallArg})
	} else if constant.IsLpTokenHaSui(nemoConfig.CoinType){
		callArgs = append(callArgs, sui_types.CallArg{Object: priceOracleCallArg}, sui_types.CallArg{Object: oracleTicketCallArg}, sui_types.CallArg{Object: stakeingCallArg}, sui_types.CallArg{Object: vaultCallArg}, sui_types.CallArg{Object: poolCallArg}, sui_types.CallArg{Object: syStateCallArg})
	}else {
		return nil, errors.New("coinType not has support oracle")
	}

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

func GetPriceVoucherFromBuck(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig) (*sui_types.Argument,error) {
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.OraclePackage)
	if err != nil {
		return nil, err
	}

	moduleName := "buck"
	functionName := "get_price_voucher_from_ssbuck"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)
	syStructTag, err := GetStructTag(nemoConfig.SyCoinType)
	if err != nil {
		return nil, err
	}
	coinTypeStructTag, err := GetStructTag(nemoConfig.CoinType)
	if err != nil {
		return nil, err
	}
	syTypeTag := move_types.TypeTag{
		Struct: syStructTag,
	}
	coinTypeTypeTag := move_types.TypeTag{
		Struct: coinTypeStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, syTypeTag, coinTypeTypeTag)

	priceOracleCallArg,err := GetObjectArg(client, nemoConfig.PriceOracle, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	oracleTicketCallArg,err := GetObjectArg(client, nemoConfig.OracleTicket, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	vaultCallArg,err := GetObjectArg(client, STSBUCK_VAULT, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	clockCallArg,err := GetObjectArg(client, constant.CLOCK, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs, sui_types.CallArg{Object: priceOracleCallArg}, sui_types.CallArg{Object: oracleTicketCallArg}, sui_types.CallArg{Object: vaultCallArg}, sui_types.CallArg{Object: clockCallArg})

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

func GetPriceVoucherFromMsTable(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig) (*sui_types.Argument,error) {
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.OraclePackage)
	if err != nil {
		return nil, err
	}

	moduleName := "aftermath"
	functionName := "get_meta_coin_price_voucher"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)
	syStructTag, err := GetStructTag(nemoConfig.SyCoinType)
	if err != nil {
		return nil, err
	}
	coinTypeStructTag, err := GetStructTag(nemoConfig.CoinType)
	if err != nil {
		return nil, err
	}
	syTypeTag := move_types.TypeTag{
		Struct: syStructTag,
	}
	coinTypeTypeTag := move_types.TypeTag{
		Struct: coinTypeStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, syTypeTag, coinTypeTypeTag)

	priceOracleCallArg,err := GetObjectArg(client, nemoConfig.PriceOracle, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	oracleTicketCallArg,err := GetObjectArg(client, nemoConfig.OracleTicket, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	registryCallArg,err := GetObjectArg(client, MSTABLE_REGISTRY, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	vaultCallArg,err := GetObjectArg(client, MSTABLE_VAULT, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	syStateCallArg,err := GetObjectArg(client, nemoConfig.SyState, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs, sui_types.CallArg{Object: priceOracleCallArg}, sui_types.CallArg{Object: oracleTicketCallArg}, sui_types.CallArg{Object: registryCallArg}, sui_types.CallArg{Object: vaultCallArg}, sui_types.CallArg{Object: syStateCallArg})

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

func GetPriceVoucherFromHaWal(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig) (*sui_types.Argument,error) {
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.OraclePackage)
	if err != nil {
		return nil, err
	}

	moduleName := "haedal"
	functionName := "get_haWAL_price_voucher"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)
	syStructTag, err := GetStructTag(nemoConfig.SyCoinType)
	if err != nil {
		return nil, err
	}
	coinTypeStructTag, err := GetStructTag(nemoConfig.CoinType)
	if err != nil {
		return nil, err
	}
	syTypeTag := move_types.TypeTag{
		Struct: syStructTag,
	}
	coinTypeTypeTag := move_types.TypeTag{
		Struct: coinTypeStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, syTypeTag, coinTypeTypeTag)

	priceOracleCallArg,err := GetObjectArg(client, nemoConfig.PriceOracle, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	oracleTicketCallArg,err := GetObjectArg(client, nemoConfig.OracleTicket, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	hawalStakingCallArg,err := GetObjectArg(client, HAWAL_STAKING, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	syStateCallArg,err := GetObjectArg(client, nemoConfig.SyState, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs, sui_types.CallArg{Object: priceOracleCallArg}, sui_types.CallArg{Object: oracleTicketCallArg}, sui_types.CallArg{Object: hawalStakingCallArg}, sui_types.CallArg{Object: syStateCallArg})

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

func GetPriceVoucherFromWWal(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig) (*sui_types.Argument,error) {
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.OraclePackage)
	if err != nil {
		return nil, err
	}

	moduleName := "haedal"
	functionName := "get_price_voucher_from_blizzard"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)
	syStructTag, err := GetStructTag(nemoConfig.SyCoinType)
	if err != nil {
		return nil, err
	}
	coinTypeStructTag, err := GetStructTag(nemoConfig.CoinType)
	if err != nil {
		return nil, err
	}
	syTypeTag := move_types.TypeTag{
		Struct: syStructTag,
	}
	coinTypeTypeTag := move_types.TypeTag{
		Struct: coinTypeStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, syTypeTag, coinTypeTypeTag)

	priceOracleCallArg,err := GetObjectArg(client, nemoConfig.PriceOracle, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	oracleTicketCallArg,err := GetObjectArg(client, nemoConfig.OracleTicket, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	wWalBlizzardStakingCallArg,err := GetObjectArg(client, WWAL_BLIZZARD_STAKING, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	walrusStakingCallArg,err := GetObjectArg(client, WALRUS_STAKING, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	syStateCallArg,err := GetObjectArg(client, nemoConfig.SyState, false, nemoConfig.OraclePackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs, sui_types.CallArg{Object: priceOracleCallArg}, sui_types.CallArg{Object: oracleTicketCallArg}, sui_types.CallArg{Object: wWalBlizzardStakingCallArg}, sui_types.CallArg{Object: walrusStakingCallArg}, sui_types.CallArg{Object: syStateCallArg})

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

func GetPriceVoucher(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig) (*sui_types.Argument,error){
	if constant.IsScallopCoin(nemoConfig.CoinType) || nemoConfig.ProviderProtocol == constant.SCALLOP{
		return GetPriceVoucherFromXOracle(ptb, client, nemoConfig)
	}else if constant.IsVSui(nemoConfig.CoinType){
		return GetPriceVoucherFromVolo(ptb, client, nemoConfig)
	}else if constant.IsSpringSui(nemoConfig.CoinType){
		return GetPriceVoucherFromSpring(ptb, client, nemoConfig, constant.SPRINGLSTINFO, "spring")
	}else if constant.IsAfSui(nemoConfig.CoinType) {
		return GetPriceVoucherFromAftermath(ptb, client, nemoConfig)
	}else if constant.IsHaSui(nemoConfig.CoinType) {
		return GetPriceVoucherFromHasui(ptb, client, nemoConfig)
	}else if constant.IsStSui(nemoConfig.CoinType){
		return GetPriceVoucherFromSpring(ptb, client, nemoConfig, constant.ALPHAFILSTINFO, "alphafi")
	}else if constant.IsLpTokenHaSui(nemoConfig.CoinType){
		return GetPriceVoucherFromLpToken(ptb, client, nemoConfig, LP_HASUI_VAULT, LP_HASUI_POOL,"haedal")
	}else if constant.IsLpTokenAfSui(nemoConfig.CoinType){
		return GetPriceVoucherFromLpToken(ptb, client, nemoConfig, LP_AFSUI_VAULT, LP_AFSUI_POOL,"aftermath")
	}else if constant.IsLpTokenVSui(nemoConfig.CoinType){
		return GetPriceVoucherFromLpToken(ptb, client, nemoConfig, LP_VSUI_VAULT, LP_VSUI_POOL,"volo")
	}else if constant.IsStsBuck(nemoConfig.CoinType){
		return GetPriceVoucherFromBuck(ptb, client, nemoConfig)
	}else if constant.IsSuperSui(nemoConfig.CoinType){
		return GetPriceVoucherFromMsTable(ptb, client, nemoConfig)
	}else if constant.IsHaWal(nemoConfig.CoinType){
		return GetPriceVoucherFromHaWal(ptb, client, nemoConfig)
	}else if constant.IsWWal(nemoConfig.CoinType){
		return GetPriceVoucherFromWWal(ptb, client, nemoConfig)
	}
	return nil, errors.New("coinType oracle not support！")
}