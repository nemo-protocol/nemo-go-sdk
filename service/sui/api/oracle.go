package api

import (
	"errors"
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/move_types"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/nemo-protocol/nemo-go-sdk/service/sui/common/constant"
	"github.com/nemo-protocol/nemo-go-sdk/service/sui/common/models"
)

const (
	MMT_ORACLE_PACKAGE_ID = "0x45fe3ef1ed2d9b444b8041a84e426242ad129483bf56000cdd514b8065967f4d"
	MMT_REGISTRY_ID = "0x6f8c395de3f250e08c01a25500c185d74cb182002d76750189c7e20a514befa8"
	MMT_ORACLE_STATE = "0x1f9310238ee9298fb703c3419030b35b22bb1cc37113e3bb5007c99aec79e5b8"
	PRICE_ADAPTER_PACKAGE_ID = "0x454ae856685130db7e5f86851add03d8252cebddbaff59db196548f3bd93d32a"
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

	shareObjectMap := map[string]bool{
		nemoConfig.PriceOracle: false,
		nemoConfig.OracleTicket: false,
		SCALLOP_VERSION: false,
		SCALLOP_MARKET_OBJECT: false,
		nemoConfig.SyState: false,
		constant.CLOCK: false,
	}

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, nemoConfig.OraclePackage, moduleName, functionName, nemoConfig.CacheContractPackageInfo[nemoConfig.OraclePackage])
	if err != nil{
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs,
		sui_types.CallArg{Object: objectArgMap[nemoConfig.PriceOracle]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.OracleTicket]},
		sui_types.CallArg{Object: objectArgMap[SCALLOP_VERSION]},
		sui_types.CallArg{Object: objectArgMap[SCALLOP_MARKET_OBJECT]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.SyState]},
		sui_types.CallArg{Object: objectArgMap[constant.CLOCK]},
	)

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

	shareObjectMap := map[string]bool{
		nemoConfig.PriceOracle: false,
		nemoConfig.OracleTicket: false,
		NATIVE_POOL: false,
		METADATA: false,
		nemoConfig.SyState: false,
	}

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, nemoConfig.OraclePackage, moduleName, functionName, nemoConfig.CacheContractPackageInfo[nemoConfig.OraclePackage])
	if err != nil{
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs,
		sui_types.CallArg{Object: objectArgMap[nemoConfig.PriceOracle]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.OracleTicket]},
		sui_types.CallArg{Object: objectArgMap[NATIVE_POOL]},
		sui_types.CallArg{Object: objectArgMap[METADATA]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.SyState]},
	)

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

	shareObjectMap := map[string]bool{
		nemoConfig.PriceOracle: false,
		nemoConfig.OracleTicket: false,
		nemoConfig.LstInfo: false,
		nemoConfig.SyState: false,
	}

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, nemoConfig.OraclePackage, moduleName, functionName, nemoConfig.CacheContractPackageInfo[nemoConfig.OraclePackage])
	if err != nil{
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs,
		sui_types.CallArg{Object: objectArgMap[nemoConfig.PriceOracle]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.OracleTicket]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.LstInfo]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.SyState]},
	)

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

	shareObjectMap := map[string]bool{
		nemoConfig.PriceOracle: false,
		nemoConfig.OracleTicket: false,
		STAKED_SUI_VAULT: false,
		SAFE: false,
		nemoConfig.SyState: false,
	}

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, nemoConfig.OraclePackage, moduleName, functionName, nemoConfig.CacheContractPackageInfo[nemoConfig.OraclePackage])
	if err != nil{
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs,
		sui_types.CallArg{Object: objectArgMap[nemoConfig.PriceOracle]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.OracleTicket]},
		sui_types.CallArg{Object: objectArgMap[STAKED_SUI_VAULT]},
		sui_types.CallArg{Object: objectArgMap[SAFE]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.SyState]},
	)

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

	shareObjectMap := map[string]bool{
		nemoConfig.PriceOracle: false,
		nemoConfig.OracleTicket: false,
		HAEDAL_STAKING: false,
		nemoConfig.SyState: false,
	}

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, nemoConfig.OraclePackage, moduleName, functionName, nemoConfig.CacheContractPackageInfo[nemoConfig.OraclePackage])
	if err != nil{
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs,
		sui_types.CallArg{Object: objectArgMap[nemoConfig.PriceOracle]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.OracleTicket]},
		sui_types.CallArg{Object: objectArgMap[HAEDAL_STAKING]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.SyState]},
	)

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

	shareObjectMap := map[string]bool{
		nemoConfig.PriceOracle: false,
		nemoConfig.OracleTicket: false,
		HAEDAL_STAKING: false,
		lpVault: false,
		lpPool: false,
		nemoConfig.SyState: false,
		STAKED_SUI_VAULT: false,
		SAFE: false,
		NATIVE_POOL: false,
		METADATA: false,
	}

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, nemoConfig.OraclePackage, moduleName, functionName, nemoConfig.CacheContractPackageInfo[nemoConfig.OraclePackage])
	if err != nil{
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)

	//default hasui
	if constant.IsLpTokenAfSui(nemoConfig.CoinType){
		callArgs = append(callArgs,
			sui_types.CallArg{Object: objectArgMap[nemoConfig.PriceOracle]},
			sui_types.CallArg{Object: objectArgMap[nemoConfig.OracleTicket]},
			sui_types.CallArg{Object: objectArgMap[STAKED_SUI_VAULT]},
			sui_types.CallArg{Object: objectArgMap[SAFE]},
			sui_types.CallArg{Object: objectArgMap[lpVault]},
			sui_types.CallArg{Object: objectArgMap[lpPool]},
			sui_types.CallArg{Object: objectArgMap[nemoConfig.SyState]},
		)
	} else if constant.IsLpTokenVSui(nemoConfig.CoinType){
		callArgs = append(callArgs,
			sui_types.CallArg{Object: objectArgMap[nemoConfig.PriceOracle]},
			sui_types.CallArg{Object: objectArgMap[nemoConfig.OracleTicket]},
			sui_types.CallArg{Object: objectArgMap[NATIVE_POOL]},
			sui_types.CallArg{Object: objectArgMap[METADATA]},
			sui_types.CallArg{Object: objectArgMap[lpVault]},
			sui_types.CallArg{Object: objectArgMap[lpPool]},
			sui_types.CallArg{Object: objectArgMap[nemoConfig.SyState]},
		)
	} else if constant.IsLpTokenHaSui(nemoConfig.CoinType){
		callArgs = append(callArgs,
			sui_types.CallArg{Object: objectArgMap[nemoConfig.PriceOracle]},
			sui_types.CallArg{Object: objectArgMap[nemoConfig.OracleTicket]},
			sui_types.CallArg{Object: objectArgMap[HAEDAL_STAKING]},
			sui_types.CallArg{Object: objectArgMap[lpVault]},
			sui_types.CallArg{Object: objectArgMap[lpPool]},
			sui_types.CallArg{Object: objectArgMap[nemoConfig.SyState]},
		)
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

	shareObjectMap := map[string]bool{
		nemoConfig.PriceOracle: false,
		nemoConfig.OracleTicket: false,
		STSBUCK_VAULT: false,
		constant.CLOCK: false,
	}

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, nemoConfig.OraclePackage, moduleName, functionName, nemoConfig.CacheContractPackageInfo[nemoConfig.OraclePackage])
	if err != nil{
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs,
		sui_types.CallArg{Object: objectArgMap[nemoConfig.PriceOracle]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.OracleTicket]},
		sui_types.CallArg{Object: objectArgMap[STSBUCK_VAULT]},
		sui_types.CallArg{Object: objectArgMap[constant.CLOCK]},
	)

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

	shareObjectMap := map[string]bool{
		nemoConfig.PriceOracle: false,
		nemoConfig.OracleTicket: false,
		MSTABLE_REGISTRY: false,
		MSTABLE_VAULT: false,
		nemoConfig.SyState: false,
	}

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, nemoConfig.OraclePackage, moduleName, functionName, nemoConfig.CacheContractPackageInfo[nemoConfig.OraclePackage])
	if err != nil{
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs,
		sui_types.CallArg{Object: objectArgMap[nemoConfig.PriceOracle]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.OracleTicket]},
		sui_types.CallArg{Object: objectArgMap[MSTABLE_REGISTRY]},
		sui_types.CallArg{Object: objectArgMap[MSTABLE_VAULT]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.SyState]},
	)

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

	shareObjectMap := map[string]bool{
		nemoConfig.PriceOracle: false,
		nemoConfig.OracleTicket: false,
		HAWAL_STAKING: false,
		nemoConfig.SyState: false,
	}

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, nemoConfig.OraclePackage, moduleName, functionName, nemoConfig.CacheContractPackageInfo[nemoConfig.OraclePackage])
	if err != nil{
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs,
		sui_types.CallArg{Object: objectArgMap[nemoConfig.PriceOracle]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.OracleTicket]},
		sui_types.CallArg{Object: objectArgMap[HAWAL_STAKING]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.SyState]},
	)

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

	shareObjectMap := map[string]bool{
		nemoConfig.PriceOracle: false,
		nemoConfig.OracleTicket: false,
		nemoConfig.WinterStaking: false,
		WALRUS_STAKING: false,
		nemoConfig.SyState: false,
	}

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, nemoConfig.OraclePackage, moduleName, functionName, nemoConfig.CacheContractPackageInfo[nemoConfig.OraclePackage])
	if err != nil{
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs,
		sui_types.CallArg{Object: objectArgMap[nemoConfig.PriceOracle]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.OracleTicket]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.WinterStaking]},
		sui_types.CallArg{Object: objectArgMap[WALRUS_STAKING]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.SyState]},
	)

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

func GetPriceVoucherFromNemo(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig) (*sui_types.Argument,error) {
	_,err := PreNemoProcess(ptb, client, nemoConfig)
	if err != nil {
		return nil, err
	}

	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.OraclePackage)
	if err != nil {
		return nil, err
	}

	moduleName := "vault"
	functionName := "get_pair_price_voucher_usd_from_mmt_vault"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	syTypeStructTag, err := GetStructTag(nemoConfig.SyCoinType)
	if err != nil {
		return nil, err
	}
	syTypeTag := move_types.TypeTag{
		Struct: syTypeStructTag,
	}

	leftCoinTypeStructTag, err := GetStructTag(nemoConfig.LeftCoinType)
	if err != nil {
		return nil, err
	}
	leftCoinTypeTag := move_types.TypeTag{
		Struct: leftCoinTypeStructTag,
	}

	rightCoinTypeStructTag, err := GetStructTag(nemoConfig.RightCoinType)
	if err != nil {
		return nil, err
	}
	rightCoinTypeTag := move_types.TypeTag{
		Struct: rightCoinTypeStructTag,
	}

	vCoinTypeStructTag, err := GetStructTag(nemoConfig.CoinType)
	if err != nil {
		return nil, err
	}
	vaultCoinTypeTag := move_types.TypeTag{
		Struct: vCoinTypeStructTag,
	}

	stableTypeStructTag, err := GetStructTag(nemoConfig.StableType)
	if err != nil {
		return nil, err
	}
	stableTypeTag := move_types.TypeTag{
		Struct: stableTypeStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, syTypeTag, leftCoinTypeTag, rightCoinTypeTag, vaultCoinTypeTag, stableTypeTag)
	//marshal, err := json.Marshal(typeArguments)

	shareObjectMap := map[string]bool{
		nemoConfig.PriceOracle: false,
		nemoConfig.OracleTicket: false,
		nemoConfig.MmtOracle: false,
		nemoConfig.VaultId: false,
		nemoConfig.PoolId: false,
		nemoConfig.SyState: false,
		constant.CLOCK: false,
	}

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, nemoConfig.OraclePackage, moduleName, functionName, nemoConfig.CacheContractPackageInfo[nemoConfig.OraclePackage])
	if err != nil{
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs,
		sui_types.CallArg{Object: objectArgMap[nemoConfig.PriceOracle]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.OracleTicket]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.MmtOracle]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.VaultId]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.PoolId]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.SyState]},
		sui_types.CallArg{Object: objectArgMap[constant.CLOCK]},
	)

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

func GetPriceVoucherFromNemoMmt(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig) (*sui_types.Argument,error) {
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.OraclePackage)
	if err != nil {
		return nil, err
	}

	moduleName := "spring"
	functionName := "get_price_voucher_in_sui_from_mmt_vault"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	leftCoinTypeStructTag, err := GetStructTag(nemoConfig.LeftCoinType)
	if err != nil {
		return nil, err
	}
	leftCoinTypeTag := move_types.TypeTag{
		Struct: leftCoinTypeStructTag,
	}

	syTypeStructTag, err := GetStructTag(nemoConfig.SyCoinType)
	if err != nil {
		return nil, err
	}
	syTypeTag := move_types.TypeTag{
		Struct: syTypeStructTag,
	}

	vCoinTypeStructTag, err := GetStructTag(nemoConfig.CoinType)
	if err != nil {
		return nil, err
	}
	vaultCoinTypeTag := move_types.TypeTag{
		Struct: vCoinTypeStructTag,
	}

	stableTypeStructTag, err := GetStructTag(nemoConfig.StableType)
	if err != nil {
		return nil, err
	}
	stableTypeTag := move_types.TypeTag{
		Struct: stableTypeStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, syTypeTag, leftCoinTypeTag, vaultCoinTypeTag, stableTypeTag)
	//marshal, err := json.Marshal(typeArguments
	shareObjectMap := map[string]bool{
		nemoConfig.PriceOracle: false,
		nemoConfig.OracleTicket: false,
		nemoConfig.LstInfo: false,
		nemoConfig.VaultId: false,
		nemoConfig.PoolId: false,
		nemoConfig.SyState: false,
	}

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, nemoConfig.OraclePackage, moduleName, functionName, nemoConfig.CacheContractPackageInfo[nemoConfig.OraclePackage])
	if err != nil{
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs,
		sui_types.CallArg{Object: objectArgMap[nemoConfig.PriceOracle]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.OracleTicket]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.LstInfo]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.VaultId]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.PoolId]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.SyState]},
	)

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

func SetKOraclePrice(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, priceReceipt *sui_types.Argument, coinType string, priceOracleObjId string) (*sui_types.Argument,error){
	nemoPackageId, err := sui_types.NewObjectIdFromHex(PRICE_ADAPTER_PACKAGE_ID)
	if err != nil {
		return nil, err
	}

	moduleName := "price_source"
	functionName := "set_k_oracle_price"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	coinTypeTypeStructTag, err := GetStructTag(coinType)
	if err != nil {
		return nil, err
	}
	coinTypeTag := move_types.TypeTag{
		Struct: coinTypeTypeStructTag,
	}

	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, coinTypeTag)

	shareObjectMap := map[string]bool{
		nemoConfig.MmtOracle: false,
		MMT_REGISTRY_ID: false,
		MMT_ORACLE_STATE: false,
		priceOracleObjId: false,
		constant.CLOCK: false,
	}

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, PRICE_ADAPTER_PACKAGE_ID, moduleName, functionName, nemoConfig.CacheContractPackageInfo[PRICE_ADAPTER_PACKAGE_ID])
	if err != nil{
		return nil, err
	}
	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs,
		sui_types.CallArg{Object: objectArgMap[nemoConfig.MmtOracle]},
		sui_types.CallArg{Object: objectArgMap[MMT_REGISTRY_ID]},
		sui_types.CallArg{Object: objectArgMap[MMT_ORACLE_STATE]},
		sui_types.CallArg{Object: objectArgMap[priceOracleObjId]},
		sui_types.CallArg{Object: objectArgMap[constant.CLOCK]},
	)

	boolArguemnt,err := ptb.Pure(true)
	if err != nil{
		return nil, err
	}

	var arguments []sui_types.Argument
	arguments = append(arguments, *priceReceipt)
	for k, v := range callArgs {
		if k == 2{
			arguments = append(arguments, boolArguemnt)
		}
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

func SetMmtOraclePriceSuiPair(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, priceReceipt *sui_types.Argument, coinType string, priceOracleObjId string) (*sui_types.Argument,error){
	nemoPackageId, err := sui_types.NewObjectIdFromHex(PRICE_ADAPTER_PACKAGE_ID)
	if err != nil {
		return nil, err
	}

	moduleName := "price_source"
	functionName := "set_mmt_oracle_price_sui_pair"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	coinTypeTypeStructTag, err := GetStructTag(coinType)
	if err != nil {
		return nil, err
	}
	coinTypeTag := move_types.TypeTag{
		Struct: coinTypeTypeStructTag,
	}

	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, coinTypeTag)

	shareObjectMap := map[string]bool{
		nemoConfig.MmtOracle: false,
		MMT_REGISTRY_ID: false,
		MMT_ORACLE_STATE: false,
		priceOracleObjId: false,
		constant.CLOCK: false,
		nemoConfig.PoolId: false,
	}

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, PRICE_ADAPTER_PACKAGE_ID, moduleName, functionName, nemoConfig.CacheContractPackageInfo[PRICE_ADAPTER_PACKAGE_ID])
	if err != nil{
		return nil, err
	}
	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs,
		sui_types.CallArg{Object: objectArgMap[nemoConfig.PoolId]},
		sui_types.CallArg{Object: objectArgMap[nemoConfig.MmtOracle]},
		sui_types.CallArg{Object: objectArgMap[MMT_ORACLE_STATE]},
		sui_types.CallArg{Object: objectArgMap[MMT_REGISTRY_ID]},
		sui_types.CallArg{Object: objectArgMap[priceOracleObjId]},
		sui_types.CallArg{Object: objectArgMap[constant.CLOCK]},
	)

	var arguments []sui_types.Argument
	arguments = append(arguments, *priceReceipt)
	for _, v := range callArgs {
		argument, err := ptb.Input(v)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, argument)
	}
	boolArguemnt,err := ptb.Pure(true)
	if err != nil{
		return nil, err
	}
	arguments = append(arguments, boolArguemnt)

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

func UpdatePrice(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, priceReceipt *sui_types.Argument, coinType string) (*sui_types.Argument,error){
	nemoPackageId, err := sui_types.NewObjectIdFromHex(MMT_ORACLE_PACKAGE_ID)
	if err != nil {
		return nil, err
	}

	moduleName := "oracle"
	functionName := "update_price"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	coinTypeTypeStructTag, err := GetStructTag(coinType)
	if err != nil {
		return nil, err
	}
	coinTypeTag := move_types.TypeTag{
		Struct: coinTypeTypeStructTag,
	}

	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, coinTypeTag)

	shareObjectMap := map[string]bool{
		nemoConfig.MmtOracle: false,
	}

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, MMT_ORACLE_PACKAGE_ID, moduleName, functionName, nemoConfig.CacheContractPackageInfo[MMT_ORACLE_PACKAGE_ID])
	if err != nil{
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs,
		sui_types.CallArg{Object: objectArgMap[nemoConfig.MmtOracle]},
	)

	var arguments []sui_types.Argument
	arguments = append(arguments, *priceReceipt)
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

func GetPriceReceipt(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, coinType string) (*sui_types.Argument,error){
	nemoPackageId, err := sui_types.NewObjectIdFromHex(MMT_ORACLE_PACKAGE_ID)
	if err != nil {
		return nil, err
	}

	moduleName := "oracle"
	functionName := "get_price_receipt"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	coinTypeTypeStructTag, err := GetStructTag(coinType)
	if err != nil {
		return nil, err
	}
	coinTypeTag := move_types.TypeTag{
		Struct: coinTypeTypeStructTag,
	}

	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, coinTypeTag)

	shareObjectMap := map[string]bool{
		nemoConfig.MmtOracle: false,
	}

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, MMT_ORACLE_PACKAGE_ID, moduleName, functionName, nemoConfig.CacheContractPackageInfo[MMT_ORACLE_PACKAGE_ID])
	if err != nil{
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs,
		sui_types.CallArg{Object: objectArgMap[nemoConfig.MmtOracle]},
	)

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

func PreNemoProcess(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig) (*sui_types.Argument,error){
	priceReceiptA, err := GetPriceReceipt(ptb, client, nemoConfig, nemoConfig.LeftCoinType)
	if err != nil{
		return nil, err
	}

	if constant.IsSui(nemoConfig.LeftCoinType) {
		_, err = SetMmtOraclePriceSuiPair(ptb, client, nemoConfig, priceReceiptA, nemoConfig.LeftCoinType, nemoConfig.LeftPriceInfoObjectId)
		if err != nil{
			return nil, err
		}
	}else {
		_, err = SetKOraclePrice(ptb, client, nemoConfig, priceReceiptA, nemoConfig.LeftCoinType, nemoConfig.LeftPriceInfoObjectId)
		if err != nil{
			return nil, err
		}
	}

	_, err = UpdatePrice(ptb, client, nemoConfig, priceReceiptA, nemoConfig.LeftCoinType)
	if err != nil{
		return nil, err
	}

	priceReceiptB, err := GetPriceReceipt(ptb, client, nemoConfig, nemoConfig.RightCoinType)
	if err != nil{
		return nil, err
	}

	if constant.IsSui(nemoConfig.RightCoinType) {
		_, err = SetMmtOraclePriceSuiPair(ptb, client, nemoConfig, priceReceiptB, nemoConfig.RightCoinType, nemoConfig.RightPriceInfoObjectId)
		if err != nil{
			return nil, err
		}
	}else {
		_, err = SetKOraclePrice(ptb, client, nemoConfig, priceReceiptB, nemoConfig.RightCoinType, nemoConfig.RightPriceInfoObjectId)
		if err != nil{
			return nil, err
		}
	}

	_, err = UpdatePrice(ptb, client, nemoConfig, priceReceiptB, nemoConfig.RightCoinType)
	if err != nil{
		return nil, err
	}

	return nil, nil
}

func GetPriceVoucher(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, cacheContractPackageInfo ...string) (*sui_types.Argument,error){
	if constant.IsScallopCoin(nemoConfig.CoinType) || nemoConfig.ProviderProtocol == constant.SCALLOP{
		return GetPriceVoucherFromXOracle(ptb, client, nemoConfig)
	}else if constant.IsVSui(nemoConfig.CoinType){
		return GetPriceVoucherFromVolo(ptb, client, nemoConfig)
	}else if constant.IsSpringCoin(nemoConfig.ProviderProtocol){
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
	}else if constant.IsWinterCoin(nemoConfig.ProviderProtocol){
		return GetPriceVoucherFromWWal(ptb, client, nemoConfig)
	}else if nemoConfig.ProviderProtocol == constant.Nemo{
		hasSuiPair := constant.IsSui(nemoConfig.LeftCoinType) || constant.IsSui(nemoConfig.RightCoinType)
		if nemoConfig.VaultId != "" && hasSuiPair{
			return GetPriceVoucherFromNemoMmt(ptb, client, nemoConfig)
		}
		return GetPriceVoucherFromNemo(ptb, client, nemoConfig)
	}
	return nil, errors.New("coinType oracle not support！")
}