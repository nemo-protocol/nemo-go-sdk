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

func GetPriceVoucherFromXOracle(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, cacheContractPackageInfo ...string) (*sui_types.Argument,error) {
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

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, nemoConfig.OraclePackage, moduleName, functionName, cacheContractPackageInfo...)
	if err != nil{
		return nil, err
	}

	fmt.Printf("\n==objectArgMap:%+v==\n",objectArgMap)
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

func GetPriceVoucherFromVolo(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, cacheContractPackageInfo ...string) (*sui_types.Argument,error) {
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

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, nemoConfig.OraclePackage, moduleName, functionName, cacheContractPackageInfo...)
	if err != nil{
		return nil, err
	}

	fmt.Printf("\n==objectArgMap:%+v==\n",objectArgMap)
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

func GetPriceVoucherFromSpring(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, lstInfo string, moduleName string, cacheContractPackageInfo ...string) (*sui_types.Argument,error) {
	fmt.Printf("nemoConfig:%+v==\n",nemoConfig)
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

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, nemoConfig.OraclePackage, moduleName, functionName, cacheContractPackageInfo...)
	if err != nil{
		return nil, err
	}

	fmt.Printf("\n==objectArgMap:%+v==\n",objectArgMap)
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

func GetPriceVoucherFromAftermath(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, cacheContractPackageInfo ...string) (*sui_types.Argument,error) {
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

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, nemoConfig.OraclePackage, moduleName, functionName, cacheContractPackageInfo...)
	if err != nil{
		return nil, err
	}

	fmt.Printf("\n==objectArgMap:%+v==\n",objectArgMap)
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

func GetPriceVoucherFromHasui(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, cacheContractPackageInfo ...string) (*sui_types.Argument,error) {
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

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, nemoConfig.OraclePackage, moduleName, functionName, cacheContractPackageInfo...)
	if err != nil{
		return nil, err
	}

	fmt.Printf("\n==objectArgMap:%+v==\n",objectArgMap)
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

func GetPriceVoucherFromLpToken(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, lpVault, lpPool, moduleName string, cacheContractPackageInfo ...string) (*sui_types.Argument,error) {
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

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, nemoConfig.OraclePackage, moduleName, functionName, cacheContractPackageInfo...)
	if err != nil{
		return nil, err
	}

	fmt.Printf("\n==objectArgMap:%+v==\n",objectArgMap)
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

func GetPriceVoucherFromBuck(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, cacheContractPackageInfo ...string) (*sui_types.Argument,error) {
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

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, nemoConfig.OraclePackage, moduleName, functionName, cacheContractPackageInfo...)
	if err != nil{
		return nil, err
	}

	fmt.Printf("\n==objectArgMap:%+v==\n",objectArgMap)
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

func GetPriceVoucherFromMsTable(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, cacheContractPackageInfo ...string) (*sui_types.Argument,error) {
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

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, nemoConfig.OraclePackage, moduleName, functionName, cacheContractPackageInfo...)
	if err != nil{
		return nil, err
	}

	fmt.Printf("\n==objectArgMap:%+v==\n",objectArgMap)
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

func GetPriceVoucherFromHaWal(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, cacheContractPackageInfo ...string) (*sui_types.Argument,error) {
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

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, nemoConfig.OraclePackage, moduleName, functionName, cacheContractPackageInfo...)
	if err != nil{
		return nil, err
	}

	fmt.Printf("\n==objectArgMap:%+v==\n",objectArgMap)
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

func GetPriceVoucherFromWWal(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, cacheContractPackageInfo ...string) (*sui_types.Argument,error) {
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

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, nemoConfig.OraclePackage, moduleName, functionName, cacheContractPackageInfo...)
	if err != nil{
		return nil, err
	}

	fmt.Printf("\n==objectArgMap:%+v==\n",objectArgMap)
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

func GetPriceVoucher(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, cacheContractPackageInfo ...string) (*sui_types.Argument,error){
	if constant.IsScallopCoin(nemoConfig.CoinType) || nemoConfig.ProviderProtocol == constant.SCALLOP{
		return GetPriceVoucherFromXOracle(ptb, client, nemoConfig, cacheContractPackageInfo...)
	}else if constant.IsVSui(nemoConfig.CoinType){
		return GetPriceVoucherFromVolo(ptb, client, nemoConfig, cacheContractPackageInfo...)
	}else if constant.IsSpringCoin(nemoConfig.ProviderProtocol){
		return GetPriceVoucherFromSpring(ptb, client, nemoConfig, constant.SPRINGLSTINFO, "spring", cacheContractPackageInfo...)
	}else if constant.IsAfSui(nemoConfig.CoinType) {
		return GetPriceVoucherFromAftermath(ptb, client, nemoConfig, cacheContractPackageInfo...)
	}else if constant.IsHaSui(nemoConfig.CoinType) {
		return GetPriceVoucherFromHasui(ptb, client, nemoConfig, cacheContractPackageInfo...)
	}else if constant.IsStSui(nemoConfig.CoinType){
		return GetPriceVoucherFromSpring(ptb, client, nemoConfig, constant.ALPHAFILSTINFO, "alphafi", cacheContractPackageInfo...)
	}else if constant.IsLpTokenHaSui(nemoConfig.CoinType){
		return GetPriceVoucherFromLpToken(ptb, client, nemoConfig, LP_HASUI_VAULT, LP_HASUI_POOL,"haedal", cacheContractPackageInfo...)
	}else if constant.IsLpTokenAfSui(nemoConfig.CoinType){
		return GetPriceVoucherFromLpToken(ptb, client, nemoConfig, LP_AFSUI_VAULT, LP_AFSUI_POOL,"aftermath", cacheContractPackageInfo...)
	}else if constant.IsLpTokenVSui(nemoConfig.CoinType){
		return GetPriceVoucherFromLpToken(ptb, client, nemoConfig, LP_VSUI_VAULT, LP_VSUI_POOL,"volo", cacheContractPackageInfo...)
	}else if constant.IsStsBuck(nemoConfig.CoinType){
		return GetPriceVoucherFromBuck(ptb, client, nemoConfig, cacheContractPackageInfo...)
	}else if constant.IsSuperSui(nemoConfig.CoinType){
		return GetPriceVoucherFromMsTable(ptb, client, nemoConfig, cacheContractPackageInfo...)
	}else if constant.IsHaWal(nemoConfig.CoinType){
		return GetPriceVoucherFromHaWal(ptb, client, nemoConfig, cacheContractPackageInfo...)
	}else if constant.IsWinterCoin(nemoConfig.ProviderProtocol){
		return GetPriceVoucherFromWWal(ptb, client, nemoConfig, cacheContractPackageInfo...)
	}
	return nil, errors.New("coinType oracle not support！")
}