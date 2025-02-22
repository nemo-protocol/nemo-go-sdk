package api

import (
	"errors"
	"fmt"
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/move_types"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"nemo-go-sdk/service/sui/common/constant"
	"nemo-go-sdk/service/sui/common/models"
)

var (
	SCALLOP_PACKAGE = "0x80ca577876dec91ae6d22090e56c39bc60dce9086ab0729930c6900bc4162b4c"
	SCOIN_TREASURY  = "0x5c1678c8261ac9eec024d4d630006a9f55c80dc0b1aa38a003fcb1d425818c6b"

	SCALLOP_VERSION       = "0x07871c4b3c847a0f674510d4978d5cf6f960452795e8ff6f189fd2088a3f6ac7"
	SCALLOP_MARKET_OBJECT = "0xa757975255146dc9686aa823b7838b507f315d704f428cbadad2f4ea061939d9"
	SCALLOP_MINT_PACKAGE  = "0x3fc1f14ca1017cff1df9cd053ce1f55251e9df3019d728c7265f028bb87f0f97"

	AFTERMATH_PACKAGE = "0x7f6ce7ade63857c4fd16ef7783fed2dfc4d7fb7e40615abdb653030b76aef0c6"
	STAKED_SUI_VAULT  = "0x2f8f6d5da7f13ea37daa397724280483ed062769813b6f31e9788e59cc88994d"
	SYSTEM_STATE      = "0x5"
	REFERRAL_VAULT    = "0x4ce9a19b594599536c53edb25d22532f82f18038dc8ef618afd00fbbfb9845ef"
	MYSTEN_2          = "0xcb7efe4253a0fe58df608d8a2d3c0eea94b4b40a8738c8daae4eb77830c16cd7"
	SAFE              = "0xeb685899830dd5837b47007809c76d91a098d52aabbf61e8ac467c59e5cc4610"

	SPRING_PACKAGE      = "0x82e6f4f75441eae97d2d5850f41a09d28c7b64a05b067d37748d471f43aaf3f7"
	LIQUID_STAKING_INFO = "0x15eda7330c8f99c30e430b4d82fd7ab2af3ead4ae17046fcb224aa9bad394f6b"

	VOLO_PACKAGE = "0x549e8b69270defbfafd4f94e17ec44cdbdd99820b33bda2278dea3b9a32d3f55"
	NATIVE_POOL  = "0x7fa2faa111b8c65bea48a23049bfd81ca8f971a262d981dcd9a17c3825cb5baf"
	METEDATA     = "0x680cd26af32b2bde8d3361e804c53ec1d1cfe24c7f039eb7f549e8dfde389a60"

	HAEDAL_PACKAGE = "0x3f45767c1aa95b25422f675800f02d8a813ec793a00b60667d071a77ba7178a2"
	HAEDAL_STAKING = "0x47b224762220393057ebf4f70501b6e657c3e56684737568439a04f80849b2ca"
)

func MintSCoin(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, coinType, underlyingCoinType string, marketCoin *sui_types.Argument) (*sui_types.Argument, error) {
	scallopMintSPackage, err := sui_types.NewObjectIdFromHex(SCALLOP_PACKAGE)
	if err != nil {
		return nil, err
	}

	module := move_types.Identifier("s_coin_converter")
	function := move_types.Identifier("mint_s_coin")
	sCoinStructTag, err := GetStructTag(coinType)
	if err != nil {
		return nil, err
	}
	type1Tag := move_types.TypeTag{
		Struct: sCoinStructTag,
	}
	underlyingCoinStructTag, err := GetStructTag(underlyingCoinType)
	if err != nil {
		return nil, err
	}
	type2Tag := move_types.TypeTag{
		Struct: underlyingCoinStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, type1Tag, type2Tag)

	scaTreasuryCallArg, err := GetObjectArg(client, SCOIN_TREASURY, false, SCALLOP_PACKAGE, "s_coin_converter", "mint_s_coin")
	if err != nil {
		return nil, err
	}
	scaTreasuryArgument, err := ptb.Input(sui_types.CallArg{Object: scaTreasuryCallArg})
	if err != nil {
		return nil, err
	}

	var arguments []sui_types.Argument

	arguments = append(arguments, scaTreasuryArgument, *marketCoin)
	command := ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *scallopMintSPackage,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)
	return &command, nil
}

func Mint(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, underlyingCoinType string, coinArgument *sui_types.Argument) (*sui_types.Argument, error) {
	scallopMintPackage, err := sui_types.NewObjectIdFromHex(SCALLOP_MINT_PACKAGE)
	if err != nil {
		return nil, err
	}
	module := move_types.Identifier("mint")
	function := move_types.Identifier("mint")

	underlyingCoinStructTag, err := GetStructTag(underlyingCoinType)
	if err != nil {
		return nil, err
	}
	type1Tag := move_types.TypeTag{
		Struct: underlyingCoinStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, type1Tag)

	versionCallArg, err := GetObjectArg(client, SCALLOP_VERSION, false, SCALLOP_MINT_PACKAGE, "mint", "mint")
	if err != nil {
		return nil, err
	}

	marketObjectCallArg, err := GetObjectArg(client, SCALLOP_MARKET_OBJECT, false, SCALLOP_MINT_PACKAGE, "mint", "mint")
	if err != nil {
		return nil, err
	}

	clockCallArg, err := GetObjectArg(client, constant.CLOCK, false, SCALLOP_MINT_PACKAGE, "mint", "mint")
	if err != nil {
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs, sui_types.CallArg{Object: versionCallArg}, sui_types.CallArg{Object: marketObjectCallArg})
	var arguments []sui_types.Argument
	for _, v := range callArgs {
		argument, err := ptb.Input(v)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, argument)
	}

	arguments = append(arguments, *coinArgument)
	clockArgument, err := ptb.Input(sui_types.CallArg{Object: clockCallArg})
	arguments = append(arguments, clockArgument)

	command := ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *scallopMintPackage,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)
	return &command, nil
}

func BurnSCoin(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, coinType, underlyingCoinType string, coinArgument *sui_types.Argument) (*sui_types.Argument, error) {
	scallopMintSPackage, err := sui_types.NewObjectIdFromHex(SCALLOP_PACKAGE)
	if err != nil {
		return nil, err
	}

	module := move_types.Identifier("s_coin_converter")
	function := move_types.Identifier("burn_s_coin")
	sCoinStructTag, err := GetStructTag(coinType)
	if err != nil {
		return nil, err
	}
	type1Tag := move_types.TypeTag{
		Struct: sCoinStructTag,
	}
	underlyingCoinStructTag, err := GetStructTag(underlyingCoinType)
	if err != nil {
		return nil, err
	}
	type2Tag := move_types.TypeTag{
		Struct: underlyingCoinStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, type1Tag, type2Tag)

	scaTreasuryCallArg, err := GetObjectArg(client, SCOIN_TREASURY, false, SCALLOP_PACKAGE, "s_coin_converter", "mint_s_coin")
	if err != nil {
		return nil, err
	}
	scaTreasuryArgument, err := ptb.Input(sui_types.CallArg{Object: scaTreasuryCallArg})
	if err != nil {
		return nil, err
	}

	var arguments []sui_types.Argument

	arguments = append(arguments, scaTreasuryArgument, *coinArgument)
	command := ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *scallopMintSPackage,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)

	marketCoin, err := Redeem(ptb, client, underlyingCoinType, &command)
	if err != nil {
		return nil, err
	}
	return marketCoin, nil
}

func Redeem(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, underlyingCoinType string, coinArgument *sui_types.Argument) (*sui_types.Argument, error) {
	scallopMintPackage, err := sui_types.NewObjectIdFromHex(SCALLOP_MINT_PACKAGE)
	if err != nil {
		return nil, err
	}
	moduleName := "redeem"
	functionName := "redeem"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	underlyingCoinStructTag, err := GetStructTag(underlyingCoinType)
	if err != nil {
		return nil, err
	}
	type1Tag := move_types.TypeTag{
		Struct: underlyingCoinStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, type1Tag)

	versionCallArg, err := GetObjectArg(client, SCALLOP_VERSION, false, SCALLOP_MINT_PACKAGE, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	marketObjectCallArg, err := GetObjectArg(client, SCALLOP_MARKET_OBJECT, false, SCALLOP_MINT_PACKAGE, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	clockCallArg, err := GetObjectArg(client, constant.CLOCK, false, SCALLOP_MINT_PACKAGE, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs, sui_types.CallArg{Object: versionCallArg}, sui_types.CallArg{Object: marketObjectCallArg})
	var arguments []sui_types.Argument
	for _, v := range callArgs {
		argument, err := ptb.Input(v)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, argument)
	}

	arguments = append(arguments, *coinArgument)
	clockArgument, err := ptb.Input(sui_types.CallArg{Object: clockCallArg})
	arguments = append(arguments, clockArgument)

	command := ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *scallopMintPackage,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)
	return &command, nil
}

func SplitCoinFromMerged(ptb *sui_types.ProgrammableTransactionBuilder, mergeCoinArgument sui_types.Argument, netSyIn uint64) (splitCoin, remainingCoin sui_types.Argument, err error) {
	splitCoinArgument, err := ptb.Pure(netSyIn)
	if err != nil {
		return sui_types.Argument{}, sui_types.Argument{}, fmt.Errorf("failed to create split coin argument: %w", err)
	}

	// 执行 SplitCoins 操作
	splitResult := ptb.Command(sui_types.Command{
		SplitCoins: &struct {
			Argument  sui_types.Argument
			Arguments []sui_types.Argument
		}{
			Argument:  mergeCoinArgument,                       // 源 coin
			Arguments: []sui_types.Argument{splitCoinArgument}, // 要分割出的数量
		},
	})

	return splitResult, remainingCoin, nil
}

func MintToSCoin(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, sCoinArgument *sui_types.Argument) (underlyingCoinArgument *sui_types.Argument, err error) {
	if constant.IsScallopCoin(nemoConfig.CoinType) {
		marketCoinArgument, err := Mint(ptb, client, nemoConfig.UnderlyingCoinType, sCoinArgument)
		if err != nil {
			return nil, err
		}
		argument, err := MintSCoin(ptb, client, nemoConfig.CoinType, nemoConfig.UnderlyingCoinType, marketCoinArgument)
		if err != nil {
			return nil, err
		}
		return argument, nil
	} else if constant.IsAfSui(nemoConfig.CoinType) {
		return MintAftermathCoin(ptb, client, nemoConfig, sCoinArgument)
	} else if constant.IsSpringSui(nemoConfig.CoinType) {
		return MintSpringCoin(ptb, client, nemoConfig, sCoinArgument)
	} else if constant.IsVSui(nemoConfig.CoinType) {
		return MintVoloCoin(ptb, client, nemoConfig, sCoinArgument)
	} else if constant.IsHaSui(nemoConfig.CoinType){
		return MintHaedalCoin(ptb, client, nemoConfig, sCoinArgument)
	}
	return nil, errors.New("coin not support！")
}

func MintAftermathCoin(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, sCoinArgument *sui_types.Argument) (underlyingCoinArgument *sui_types.Argument, err error) {
	scallopMintSPackage, err := sui_types.NewObjectIdFromHex(AFTERMATH_PACKAGE)
	if err != nil {
		return nil, err
	}

	moduleName := "staked_sui_vault"
	functionName := "request_stake"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)
	typeArguments := make([]move_types.TypeTag, 0)

	stakedSuiVaultArgument, err := GetObjectArgument(ptb, client, STAKED_SUI_VAULT, false, AFTERMATH_PACKAGE, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	safeArgument, err := GetObjectArgument(ptb, client, SAFE, false, AFTERMATH_PACKAGE, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	systemStateArgument, err := GetObjectArgument(ptb, client, SYSTEM_STATE, false, AFTERMATH_PACKAGE, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	referralVaultArgument, err := GetObjectArgument(ptb, client, REFERRAL_VAULT, false, AFTERMATH_PACKAGE, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	mystenArgument, err := GetObjectArgument(ptb, client, MYSTEN_2, false, AFTERMATH_PACKAGE, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	var arguments []sui_types.Argument

	arguments = append(arguments, stakedSuiVaultArgument, safeArgument, systemStateArgument, referralVaultArgument, *sCoinArgument, mystenArgument)
	command := ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:   *scallopMintSPackage,
				Module:    module,
				Function:  function,
				TypeArguments: typeArguments,
				Arguments: arguments,
			},
		},
	)
	return &command, nil
}

func MintSpringCoin(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, sCoinArgument *sui_types.Argument) (underlyingCoinArgument *sui_types.Argument, err error) {
	scallopMintSPackage, err := sui_types.NewObjectIdFromHex(SPRING_PACKAGE)
	if err != nil {
		return nil, err
	}

	moduleName := "liquid_staking"
	functionName := "mint"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	coinTypeStructTag, err := GetStructTag(nemoConfig.CoinType)
	if err != nil {
		return nil, err
	}
	type1Tag := move_types.TypeTag{
		Struct: coinTypeStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, type1Tag)

	liquidStakingInfoArgument, err := GetObjectArgument(ptb, client, LIQUID_STAKING_INFO, false, AFTERMATH_PACKAGE, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	systemStateArgument, err := GetObjectArgument(ptb, client, SYSTEM_STATE, false, AFTERMATH_PACKAGE, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	var arguments []sui_types.Argument

	arguments = append(arguments, liquidStakingInfoArgument, systemStateArgument, *sCoinArgument)
	command := ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *scallopMintSPackage,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)
	return &command, nil
}

func MintVoloCoin(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, coinArgument *sui_types.Argument) (underlyingCoinArgument *sui_types.Argument, err error) {
	scallopMintSPackage, err := sui_types.NewObjectIdFromHex(VOLO_PACKAGE)
	if err != nil {
		return nil, err
	}

	moduleName := "native_pool"
	functionName := "stake_non_entry"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	typeArguments := make([]move_types.TypeTag, 0)

	nativePoolArgument, err := GetObjectArgument(ptb, client, NATIVE_POOL, false, VOLO_PACKAGE, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	metadataArgument, err := GetObjectArgument(ptb, client, METEDATA, false, VOLO_PACKAGE, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	systemStateArgument, err := GetObjectArgument(ptb, client, SYSTEM_STATE, false, VOLO_PACKAGE, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	var arguments []sui_types.Argument

	arguments = append(arguments, nativePoolArgument, metadataArgument, systemStateArgument, *coinArgument)
	command := ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *scallopMintSPackage,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)
	return &command, nil
}

func MintHaedalCoin(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, coinArgument *sui_types.Argument) (underlyingCoinArgument *sui_types.Argument, err error) {
	scallopMintSPackage, err := sui_types.NewObjectIdFromHex(HAEDAL_PACKAGE)
	if err != nil {
		return nil, err
	}

	moduleName := "staking"
	functionName := "request_stake_coin"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	typeArguments := make([]move_types.TypeTag, 0)

	systemStateArgument, err := GetObjectArgument(ptb, client, SYSTEM_STATE, false, HAEDAL_PACKAGE, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	haedalStakingArgument, err := GetObjectArgument(ptb, client, HAEDAL_STAKING, false, HAEDAL_PACKAGE, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	address := "0x0000000000000000000000000000000000000000000000000000000000000000"
	addressArgument, err := GetObjectArgument(ptb, client, address, false, HAEDAL_PACKAGE, moduleName, functionName)
	if err != nil {
		return nil, err
	}



	var arguments []sui_types.Argument

	arguments = append(arguments, systemStateArgument, haedalStakingArgument, *coinArgument,addressArgument)
	command := ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *scallopMintSPackage,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)
	return &command, nil
}
