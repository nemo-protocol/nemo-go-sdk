package api

import (
	"errors"
	"fmt"
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/move_types"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/nemo-protocol/nemo-go-sdk/service/sui/common/constant"
	"github.com/nemo-protocol/nemo-go-sdk/service/sui/common/models"
)

var (
	SCALLOP_PACKAGE = "0x80ca577876dec91ae6d22090e56c39bc60dce9086ab0729930c6900bc4162b4c"

	SCALLOP_VERSION       = "0x07871c4b3c847a0f674510d4978d5cf6f960452795e8ff6f189fd2088a3f6ac7"
	SCALLOP_MARKET_OBJECT = "0xa757975255146dc9686aa823b7838b507f315d704f428cbadad2f4ea061939d9"
	SCALLOP_MINT_PACKAGE  = "0x83bbe0b3985c5e3857803e2678899b03f3c4a31be75006ab03faf268c014ce41"

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
	METADATA     = "0x680cd26af32b2bde8d3361e804c53ec1d1cfe24c7f039eb7f549e8dfde389a60"

	HAEDAL_PACKAGE = "0x3f45767c1aa95b25422f675800f02d8a813ec793a00b60667d071a77ba7178a2"
	HAEDAL_STAKING = "0x47b224762220393057ebf4f70501b6e657c3e56684737568439a04f80849b2ca"

	ALPHAFI_PACKAGE = "0x059f94b85c07eb74d2847f8255d8cc0a67c9a8dcc039eabf9f8b9e23a0de2700"
	ALPHAFI_STAKING = "0x1adb343ab351458e151bc392fbf1558b3332467f23bda45ae67cd355a57fd5f5"

	LP_HASUI_VAULT = "0xde97452e63505df696440f86f0b805263d8659b77b8c316739106009d514c270"
	LP_HASUI_POOL  = "0x871d8a227114f375170f149f7e9d45be822dd003eba225e83c05ac80828596bc"

	LP_AFSUI_VAULT = "0xff4cc0af0ad9d50d4a3264dfaafd534437d8b66c8ebe9f92b4c39d898d6870a3"
	LP_AFSUI_POOL  = "0xa528b26eae41bcfca488a9feaa3dca614b2a1d9b9b5c78c256918ced051d4c50"

	LP_VSUI_VAULT = "0x5732b81e659bd2db47a5b55755743dde15be99490a39717abc80d62ec812bcb6"
	LP_VSUI_POOL  = "0x6c545e78638c8c1db7a48b282bb8ca79da107993fcb185f75cedc1f5adb2f535"

	STSBUCK_VAULT           = "0xe83e455a9e99884c086c8c79c13367e7a865de1f953e75bcf3e529cdf03c6224"
	STSBUCK_PACKAGE         = "0x2a721777dc1fcf7cda19492ad7c2272ee284214652bde3e9740e2f49c3bff457"
	STSBUCK_DEPOSIT_PACKAGE = "0x75fe358d87679b30befc498a8dae1d28ca9eed159ab6f2129a654a8255e5610e"

	MSTABLE_REGISTRY = "0x5ff2396592a20f7bf6ff291963948d6fc2abec279e11f50ee74d193c4cf0bba8"
	MSTABLE_VAULT = "0x3062285974a5e517c88cf3395923aac788dce74f3640029a01e25d76c4e76f5d"
)

// 定义一个 map 来存储 coinType 和 treasury 的映射关系
var sCoinMap = map[string]string{
	"0xaafc4f740de0dd0dde642a31148fb94517087052f19afb0f7bed1dc41a50c77b::scallop_sui::SCALLOP_SUI":                     "0x5c1678c8261ac9eec024d4d630006a9f55c80dc0b1aa38a003fcb1d425818c6b",
	"0xea346ce428f91ab007210443efcea5f5cdbbb3aae7e9affc0ca93f9203c31f0c::scallop_cetus::SCALLOP_CETUS":                 "0xa283c63488773c916cb3d6c64109536160d5eb496caddc721eb39aad2977d735",
	"0x5ca17430c1d046fae9edeaa8fd76c7b4193a00d764a0ecfa9418d733ad27bc1e::scallop_sca::SCALLOP_SCA":                     "0xe04bfc95e00252bd654ee13c08edef9ac5e4b6ae4074e8390db39e9a0109c529",
	"0xad4d71551d31092230db1fd482008ea42867dbf27b286e9c70a79d2a6191d58d::scallop_wormhole_usdc::SCALLOP_WORMHOLE_USDC": "0x50c5cfcbcca3aaacab0984e4d7ad9a6ad034265bebb440f0d1cd688ec20b2548",
	"0xe6e5a012ec20a49a3d1d57bd2b67140b96cd4d3400b9d79e541f7bdbab661f95::scallop_wormhole_usdt::SCALLOP_WORMHOLE_USDT": "0x1f02e2fed702b477732d4ad6044aaed04f2e8e586a169153694861a901379df0",
	"0x67540ceb850d418679e69f1fb6b2093d6df78a2a699ffc733f7646096d552e9b::scallop_wormhole_eth::SCALLOP_WORMHOLE_ETH":   "0x4b7f5da0e306c9d52490a0c1d4091e653d6b89778b9b4f23c877e534e4d9cd21",
	"0x00671b1fa2a124f5be8bdae8b91ee711462c5d9e31bda232e70fd9607b523c88::scallop_af_sui::SCALLOP_AF_SUI":               "0x55f4dfe9e40bc4cc11c70fcb1f3daefa2bdc330567c58d4f0792fbd9f9175a62",
	"0x9a2376943f7d22f88087c259c5889925f332ca4347e669dc37d54c2bf651af3c::scallop_ha_sui::SCALLOP_HA_SUI":               "0x404ccc1404d74a90eb6f9c9d4b6cda6d417fb03189f80d9070a35e5dab1df0f5",
	"0xe1a1cc6bcf0001a015eab84bcc6713393ce20535f55b8b6f35c142e057a25fbe::scallop_v_sui::SCALLOP_V_SUI":                 "0xc06688ee1af25abc286ffb1d18ce273d1d5907cd1064c25f4e8ca61ea989c1d1",
	"0x1392650f2eca9e3f6ffae3ff89e42a3590d7102b80e2b430f674730bc30d3259::scallop_wormhole_sol::SCALLOP_WORMHOLE_SOL":   "0x760fd66f5be869af4382fa32b812b3c67f0eca1bb1ed7a5578b21d56e1848819",
	"0x2cf76a9cf5d3337961d1154283234f94da2dcff18544dfe5cbdef65f319591b5::scallop_wormhole_btc::SCALLOP_WORMHOLE_BTC":   "0xe2883934ea42c99bc998bbe0f01dd6d27aa0e27a56455707b1b34e6a41c20baa",
	"0x854950aa624b1df59fe64e630b2ba7c550642e9342267a33061d59fb31582da5::scallop_usdc::SCALLOP_USDC":                   "0xbe6b63021f3d82e0e7e977cdd718ed7c019cf2eba374b7b546220402452f938e",
	"0xb14f82d8506d139eacef109688d1b71e7236bcce9b2c0ad526abcd6aa5be7de0::scallop_sb_eth::SCALLOP_SB_ETH":               "0xfd0f02def6358a1f266acfa1493d4707ee8387460d434fb667d63d755ff907ed",
	"0x6711551c1e7652a270d9fbf0eee25d99594c157cde3cb5fbb49035eb59b1b001::scallop_fdusd::SCALLOP_FDUSD":                 "0xdad9bc6293e694f67a5274ea51b596e0bdabfafc585ae6d7e82888e65f1a03e0",
	"0xeb7a05a3224837c5e5503575aed0be73c091d1ce5e43aa3c3e716e0ae614608f::scallop_deep::SCALLOP_DEEP":                   "0xc63838fabe37b25ad897392d89876d920f5e0c6a406bf3abcb84753d2829bc88",
	"0xe56d5167f427cbe597da9e8150ef5c337839aaf46891d62468dcf80bdd8e10d1::scallop_fud::SCALLOP_FUD":                     "0xf25212f11d182decff7a86165699a73e3d5787aced203ca539f43cfbc10db867",
	"0xb1d7df34829d1513b73ba17cb7ad90c88d1e104bb65ab8f62f13e0cc103783d3::scallop_sb_usdt::SCALLOP_SB_USDT":             "0x58bdf6a9752e3a60144d0b70e8608d630dfd971513e2b2bfa7282f5eaa7d04d8",
}

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

	scallopTreasury, err := GetTreasuryByCoinType(coinType)
	if err != nil {
		return nil, err
	}

	scaTreasuryCallArg, err := GetObjectArg(client, scallopTreasury, false, SCALLOP_PACKAGE, "s_coin_converter", "mint_s_coin")
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

	scallopTreasury, err := GetTreasuryByCoinType(coinType)
	if err != nil {
		return nil, err
	}

	scaTreasuryCallArg, err := GetObjectArg(client, scallopTreasury, false, SCALLOP_PACKAGE, "s_coin_converter", "mint_s_coin")
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

func SplitCoinFromMerged(ptb *sui_types.ProgrammableTransactionBuilder, mergeCoinArgument sui_types.Argument, netSyIn uint64) (splitCoin sui_types.Argument, err error) {
	splitCoinArgument, err := ptb.Pure(netSyIn)
	if err != nil {
		return sui_types.Argument{}, fmt.Errorf("failed to create split coin argument: %w", err)
	}

	splitResult := ptb.Command(sui_types.Command{
		SplitCoins: &struct {
			Argument  sui_types.Argument
			Arguments []sui_types.Argument
		}{
			Argument:  mergeCoinArgument,
			Arguments: []sui_types.Argument{splitCoinArgument},
		},
	})

	if splitResult.Result == nil {
		return sui_types.Argument{}, fmt.Errorf("split coins command failed")
	}

	splitCoin = sui_types.Argument{
		NestedResult: &struct {
			Result1 uint16
			Result2 uint16
		}{
			Result1: *splitResult.Result,
			Result2: 0, // 第一个分割结果
		},
	}

	return splitCoin, nil
}

func MintToSCoin(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, sCoinArgument *sui_types.Argument) (underlyingCoinArgument *sui_types.Argument, err error) {
	if constant.IsScallopCoin(nemoConfig.CoinType) || nemoConfig.ProviderProtocol == constant.SCALLOP {
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
	} else if constant.IsHaSui(nemoConfig.CoinType) {
		return MintHaedalCoin(ptb, client, nemoConfig, sCoinArgument)
	} else if constant.IsStSui(nemoConfig.CoinType) {
		return MintStCoin(ptb, client, nemoConfig, sCoinArgument)
	} else if constant.IsStsBuck(nemoConfig.CoinType) {
		return MintByBuck(ptb, client, nemoConfig, sCoinArgument)
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

	metadataArgument, err := GetObjectArgument(ptb, client, METADATA, false, VOLO_PACKAGE, moduleName, functionName)
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

	address, _ := sui_types.NewAddressFromHex("0x0000000000000000000000000000000000000000000000000000000000000000")
	addressArgument, err := ptb.Pure(*address)

	var arguments []sui_types.Argument

	arguments = append(arguments, systemStateArgument, haedalStakingArgument, *coinArgument, addressArgument)
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

func MintStCoin(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, coinArgument *sui_types.Argument) (underlyingCoinArgument *sui_types.Argument, err error) {
	scallopMintSPackage, err := sui_types.NewObjectIdFromHex(ALPHAFI_PACKAGE)
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

	systemStateArgument, err := GetObjectArgument(ptb, client, SYSTEM_STATE, false, ALPHAFI_PACKAGE, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	alphafiStakingArgument, err := GetObjectArgument(ptb, client, ALPHAFI_STAKING, false, ALPHAFI_PACKAGE, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	var arguments []sui_types.Argument

	arguments = append(arguments, alphafiStakingArgument, systemStateArgument, *coinArgument)
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

func StsBuckWithdraw(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, coinBalanceArgument *sui_types.Argument) (tickerArgument *sui_types.Argument, err error) {
	stsbuckWithdrawPackage, err := sui_types.NewObjectIdFromHex(STSBUCK_PACKAGE)
	if err != nil {
		return nil, err
	}

	moduleName := "vault"
	functionName := "withdraw"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	underlyingCoinTypeStructTag, err := GetStructTag(nemoConfig.UnderlyingCoinType)
	if err != nil {
		return nil, err
	}
	type1Tag := move_types.TypeTag{
		Struct: underlyingCoinTypeStructTag,
	}
	coinTypeStructTag, err := GetStructTag(nemoConfig.CoinType)
	if err != nil {
		return nil, err
	}
	type2Tag := move_types.TypeTag{
		Struct: coinTypeStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, type1Tag, type2Tag)

	stsBuckVaultArgument, err := GetObjectArgument(ptb, client, STSBUCK_VAULT, false, STSBUCK_PACKAGE, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	clockArgument, err := GetObjectArgument(ptb, client, constant.CLOCK, false, STSBUCK_PACKAGE, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	var arguments []sui_types.Argument

	arguments = append(arguments, stsBuckVaultArgument, *coinBalanceArgument, clockArgument)
	command := ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *stsbuckWithdrawPackage,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)
	return &command, nil
}

func StsBuckWithdrawTicket(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, withdrawTicketArgument *sui_types.Argument) (balanceArgument *sui_types.Argument, err error) {
	stsbuckWithdrawPackage, err := sui_types.NewObjectIdFromHex(STSBUCK_PACKAGE)
	if err != nil {
		return nil, err
	}

	moduleName := "vault"
	functionName := "redeem_withdraw_ticket"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	underlyingCoinTypeStructTag, err := GetStructTag(nemoConfig.UnderlyingCoinType)
	if err != nil {
		return nil, err
	}
	type1Tag := move_types.TypeTag{
		Struct: underlyingCoinTypeStructTag,
	}
	coinTypeStructTag, err := GetStructTag(nemoConfig.CoinType)
	if err != nil {
		return nil, err
	}
	type2Tag := move_types.TypeTag{
		Struct: coinTypeStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, type1Tag, type2Tag)

	stsBuckVaultArgument, err := GetObjectArgument(ptb, client, STSBUCK_VAULT, false, STSBUCK_PACKAGE, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	var arguments []sui_types.Argument

	arguments = append(arguments, stsBuckVaultArgument, *withdrawTicketArgument)
	command := ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *stsbuckWithdrawPackage,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)
	return &command, nil
}

func StsBuckDeposit(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, coinBalanceArgument *sui_types.Argument) (argument *sui_types.Argument, err error) {
	stsbuckWithdrawPackage, err := sui_types.NewObjectIdFromHex(STSBUCK_DEPOSIT_PACKAGE)
	if err != nil {
		return nil, err
	}

	moduleName := "sbuck_saving_vault"
	functionName := "deposit"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	typeArguments := make([]move_types.TypeTag, 0)

	stsBuckVaultArgument, err := GetObjectArgument(ptb, client, STSBUCK_VAULT, false, STSBUCK_PACKAGE, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	clockArgument, err := GetObjectArgument(ptb, client, constant.CLOCK, false, STSBUCK_PACKAGE, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	var arguments []sui_types.Argument

	arguments = append(arguments, stsBuckVaultArgument, *coinBalanceArgument, clockArgument)
	command := ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *stsbuckWithdrawPackage,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)
	return &command, nil
}

func GetTreasuryByCoinType(coinType string) (string, error) {
	if treasury, exists := sCoinMap[coinType]; exists {
		return treasury, nil
	}
	return "", fmt.Errorf("coinType not found: %s, not support redeem to underlying coin", coinType)
}

func BurnToBuck(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, coinArgument *sui_types.Argument) (*sui_types.Argument, error) {
	coinBalanceArgument, err := CoinIntoBalance(ptb, coinArgument, nemoConfig.CoinType)
	if err != nil {
		return nil, err
	}

	withdrawTicketArgument, err := StsBuckWithdraw(ptb, client, nemoConfig, coinBalanceArgument)
	if err != nil {
		return nil, err
	}

	balanceArgument, err := StsBuckWithdrawTicket(ptb, client, nemoConfig, withdrawTicketArgument)
	if err != nil {
		return nil, err
	}

	return CoinFromBalance(ptb, balanceArgument, nemoConfig.UnderlyingCoinType)
}

func MintByBuck(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, coinArgument *sui_types.Argument) (*sui_types.Argument, error) {
	coinBalanceArgument, err := CoinIntoBalance(ptb, coinArgument, nemoConfig.UnderlyingCoinType)
	if err != nil {
		return nil, err
	}

	balanceArgument, err := StsBuckDeposit(ptb, client, nemoConfig, coinBalanceArgument)
	if err != nil {
		return nil, err
	}

	return CoinFromBalance(ptb, balanceArgument, nemoConfig.CoinType)
}
