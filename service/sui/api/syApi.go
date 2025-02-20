package api

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/move_types"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"nemo-go-sdk/service/sui/common/constant"
	"nemo-go-sdk/service/sui/common/models"
)

func Deposit(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, coinArgument *sui_types.Argument) (*sui_types.Argument,error) {
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.NemoContract)
	if err != nil {
		return nil, err
	}

	moduleName := "sy"
	functionName := "deposit"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)
	syStructTag, err := GetStructTag(nemoConfig.SyCoinType)
	if err != nil {
		return nil, err
	}
	syTypeTag := move_types.TypeTag{
		Struct: syStructTag,
	}

	structTag, err := GetStructTag(nemoConfig.CoinType)
	if err != nil {
		return nil, err
	}
	typeTag := move_types.TypeTag{
		Struct: structTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, typeTag, syTypeTag)

	var arguments []sui_types.Argument
	versionArgument,err := GetObjectArgument(ptb, client, nemoConfig.Version, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	syStateArgument,err := GetObjectArgument(ptb, client, nemoConfig.SyState, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}
	arguments = append(arguments, versionArgument, *coinArgument, syStateArgument)

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

func SeedLiquidity() {}

func SwapExactPtForSy(ptb *sui_types.ProgrammableTransactionBuilder, blockClient *sui.ISuiAPI, client *client.Client, nemoConfig *models.NemoConfig, amountIn, minSyOut uint64, ownerAddress string, oracleArgument *sui_types.Argument) (*sui_types.Argument,error){
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.NemoContract)
	if err != nil {
		return nil, err
	}

	moduleName := "market"
	functionName := "swap_exact_pt_for_sy"
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

	versionArgument,err := GetObjectArgument(ptb, client, nemoConfig.Version, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	pyStateInfo, err := GetObjectFieldByObjectId(client, nemoConfig.PyState)
	if err != nil{
		return nil, err
	}
	maturity := pyStateInfo["expiry"].(string)

	expectPyPositionTypeList := make([]string, 0)
	for _, pkg := range nemoConfig.NemoContractList{
		expectPyPositionTypeList = append(expectPyPositionTypeList, fmt.Sprintf("%v::py_position::PyPosition", pkg))
	}

	pyPosition,err := GetOwnerObjectByType(blockClient, client, expectPyPositionTypeList, nemoConfig.SyCoinType, maturity, ownerAddress)
	if err != nil {
		return nil, err
	}
	if pyPosition == ""{
		return nil, errors.New("without pyPosition！")
	}

	var pyPositionArgument *sui_types.Argument
	argument, err := GetObjectArgument(ptb, client, pyPosition, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil{
		return nil, err
	}
	pyPositionArgument = &argument

	pyStateArgument,err := GetObjectArgument(ptb, client, nemoConfig.PyState, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	marketStateArgument,err := GetObjectArgument(ptb, client, nemoConfig.MarketState, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	marketGlobalConfigArgument,err := GetObjectArgument(ptb, client, nemoConfig.MarketFactoryConfig, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	clockArgument,err := GetObjectArgument(ptb, client, constant.CLOCK, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	ptAmountIn := CreatePureU64CallArg(amountIn)
	ptAmountArgument,err := ptb.Input(ptAmountIn)
	if err != nil {
		return nil, err
	}
	syOut := CreatePureU64CallArg(minSyOut)
	syOutArgument,err := ptb.Input(syOut)
	if err != nil {
		return nil, err
	}

	var arguments []sui_types.Argument
	resultArg := &sui_types.Argument{
		NestedResult: &struct {
			Result1 uint16
			Result2 uint16
		}{Result1: *oracleArgument.Result, Result2: 0},
	}

	arguments = append(arguments, versionArgument, ptAmountArgument, syOutArgument, *pyPositionArgument, pyStateArgument, *resultArg, marketGlobalConfigArgument, marketStateArgument, clockArgument)
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

func SyRedeem(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, argument *sui_types.Argument) (*sui_types.Argument,error){
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.NemoContract)
	if err != nil {
		return nil, err
	}

	moduleName := "sy"
	functionName := "redeem"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	arg0Tag, err := GetStructTag(nemoConfig.CoinType)
	if err != nil {
		return nil, err
	}
	arg0TypeTag := move_types.TypeTag{
		Struct: arg0Tag,
	}
	arg1Tag, err := GetStructTag(nemoConfig.SyCoinType)
	if err != nil {
		return nil, err
	}
	arg1TypeTag := move_types.TypeTag{
		Struct: arg1Tag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, arg0TypeTag, arg1TypeTag)

	versionArgument,err := GetObjectArgument(ptb, client, nemoConfig.Version, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	syStateArgument,err := GetObjectArgument(ptb, client, nemoConfig.SyState, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	var arguments []sui_types.Argument
	arguments = append(arguments, versionArgument, *argument, syStateArgument)
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

func CreatePureU64CallArg(value uint64) sui_types.CallArg {
	// 创建一个8字节的缓冲区
	buf := make([]byte, 8)

	// 使用 binary.LittleEndian 将 uint64 写入字节数组
	binary.LittleEndian.PutUint64(buf, value)

	// 返回构造的 CallArg
	return sui_types.CallArg{
		Pure: &buf,
		Object: nil,  // Pure 类型不需要 Object
	}
}

func SwapExactYtForSy(ptb *sui_types.ProgrammableTransactionBuilder, blockClient *sui.ISuiAPI, client *client.Client, nemoConfig *models.NemoConfig, amountIn, minSyOut uint64, ownerAddress string, oracleArgument *sui_types.Argument) (*sui_types.Argument,error){
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.NemoContract)
	if err != nil {
		return nil, err
	}

	moduleName := "router"
	functionName := "swap_exact_yt_for_sy"
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

	versionArgument,err := GetObjectArgument(ptb, client, nemoConfig.Version, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	pyStateInfo, err := GetObjectFieldByObjectId(client, nemoConfig.PyState)
	if err != nil{
		return nil, err
	}
	maturity := pyStateInfo["expiry"].(string)

	expectPyPositionTypeList := make([]string, 0)
	for _, pkg := range nemoConfig.NemoContractList{
		expectPyPositionTypeList = append(expectPyPositionTypeList, fmt.Sprintf("%v::py_position::PyPosition", pkg))
	}

	pyPosition,err := GetOwnerObjectByType(blockClient, client, expectPyPositionTypeList, nemoConfig.SyCoinType, maturity, ownerAddress)
	if err != nil {
		return nil, err
	}
	if pyPosition == ""{
		return nil, errors.New("without pyPosition！")
	}

	var pyPositionArgument *sui_types.Argument
	argument, err := GetObjectArgument(ptb, client, pyPosition, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil{
		return nil, err
	}
	pyPositionArgument = &argument

	pyStateArgument,err := GetObjectArgument(ptb, client, nemoConfig.PyState, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	yieldFactoryArgument,err := GetObjectArgument(ptb, client, nemoConfig.YieldFactoryConfig, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	marketStateArgument,err := GetObjectArgument(ptb, client, nemoConfig.MarketState, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	marketGlobalConfigArgument,err := GetObjectArgument(ptb, client, nemoConfig.MarketFactoryConfig, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	clockArgument,err := GetObjectArgument(ptb, client, constant.CLOCK, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	ytAmountIn := CreatePureU64CallArg(amountIn)
	ytAmountArgument,err := ptb.Input(ytAmountIn)
	if err != nil {
		return nil, err
	}
	syOut := CreatePureU64CallArg(minSyOut)
	syOutArgument,err := ptb.Input(syOut)
	if err != nil {
		return nil, err
	}

	var arguments []sui_types.Argument
	resultArg := &sui_types.Argument{
		NestedResult: &struct {
			Result1 uint16
			Result2 uint16
		}{Result1: *oracleArgument.Result, Result2: 0},
	}

	arguments = append(arguments, versionArgument, ytAmountArgument, syOutArgument, *pyPositionArgument, pyStateArgument, *resultArg, yieldFactoryArgument, marketGlobalConfigArgument, marketStateArgument, clockArgument)
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

func SwapExactSyForYt(ptb *sui_types.ProgrammableTransactionBuilder, blockClient *sui.ISuiAPI, client *client.Client, nemoConfig *models.NemoConfig,ownerAddress string, approxYtOut, netSyTokenization, minYtOut uint64, oracleArgument, coinArgument *sui_types.Argument) (*sui_types.Argument,error){
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.NemoContract)
	if err != nil {
		return nil, err
	}

	moduleName := "router"
	functionName := "swap_exact_sy_for_yt"
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

	versionArgument,err := GetObjectArgument(ptb, client, nemoConfig.Version, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	pyStateInfo, err := GetObjectFieldByObjectId(client, nemoConfig.PyState)
	if err != nil{
		return nil, err
	}
	maturity := pyStateInfo["expiry"].(string)

	expectPyPositionTypeList := make([]string, 0)
	for _, pkg := range nemoConfig.NemoContractList{
		expectPyPositionTypeList = append(expectPyPositionTypeList, fmt.Sprintf("%v::py_position::PyPosition", pkg))
	}

	pyPosition,err := GetOwnerObjectByType(blockClient, client, expectPyPositionTypeList, nemoConfig.SyCoinType, maturity, ownerAddress)
	if err != nil {
		return nil, err
	}
	var pyPositionArgument *sui_types.Argument
	if pyPosition == ""{
		pyPositionArgument, err = InitPyPosition(ptb, client, nemoConfig)
		if err != nil{
			return nil, err
		}
	}else {
		argument, err := GetObjectArgument(ptb, client, pyPosition, false, nemoConfig.NemoContract, moduleName, functionName)
		if err != nil{
			return nil, err
		}
		pyPositionArgument = &argument
	}

	pyStateArgument,err := GetObjectArgument(ptb, client, nemoConfig.PyState, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	yieldFactoryArgument,err := GetObjectArgument(ptb, client, nemoConfig.YieldFactoryConfig, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	marketGlobalConfigArgument,err := GetObjectArgument(ptb, client, nemoConfig.MarketFactoryConfig, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	marketStateArgument,err := GetObjectArgument(ptb, client, nemoConfig.MarketState, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	clockArgument,err := GetObjectArgument(ptb, client, constant.CLOCK, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	approxYtOutCallArg := CreatePureU64CallArg(approxYtOut)
	approxYtOutArgument,err := ptb.Input(approxYtOutCallArg)
	if err != nil {
		return nil, err
	}

	netSyTokenizationCallArg := CreatePureU64CallArg(netSyTokenization)
	netSyTokenizationArgument,err := ptb.Input(netSyTokenizationCallArg)
	if err != nil {
		return nil, err
	}

	minYtOutCallArg := CreatePureU64CallArg(minYtOut)
	minYtOutArgument,err := ptb.Input(minYtOutCallArg)
	if err != nil {
		return nil, err
	}

	var arguments []sui_types.Argument
	arguments = append(arguments, versionArgument, minYtOutArgument, approxYtOutArgument, netSyTokenizationArgument, *coinArgument, *oracleArgument, *pyPositionArgument, pyStateArgument, yieldFactoryArgument, marketGlobalConfigArgument, marketStateArgument, clockArgument)
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

func SwapExactSyForPt(ptb *sui_types.ProgrammableTransactionBuilder, blockClient *sui.ISuiAPI, client *client.Client, nemoConfig *models.NemoConfig, ownerAddress string, approxPtOut, minYtOut uint64, oracleArgument, coinArgument *sui_types.Argument) (*sui_types.Argument,error){
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.NemoContract)
	if err != nil {
		return nil, err
	}

	moduleName := "router"
	functionName := "swap_exact_sy_for_pt"
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

	versionArgument,err := GetObjectArgument(ptb, client, nemoConfig.Version, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	pyStateInfo, err := GetObjectFieldByObjectId(client, nemoConfig.PyState)
	if err != nil{
		return nil, err
	}
	maturity := pyStateInfo["expiry"].(string)

	expectPyPositionTypeList := make([]string, 0)
	for _, pkg := range nemoConfig.NemoContractList{
		expectPyPositionTypeList = append(expectPyPositionTypeList, fmt.Sprintf("%v::py_position::PyPosition", pkg))
	}

	pyPosition,err := GetOwnerObjectByType(blockClient, client, expectPyPositionTypeList, nemoConfig.SyCoinType, maturity, ownerAddress)
	if err != nil {
		return nil, err
	}
	var pyPositionArgument *sui_types.Argument
	if pyPosition == ""{
		pyPositionArgument, err = InitPyPosition(ptb, client, nemoConfig)
		if err != nil{
			return nil, err
		}
	}else {
		argument, err := GetObjectArgument(ptb, client, pyPosition, false, nemoConfig.NemoContract, moduleName, functionName)
		if err != nil{
			return nil, err
		}
		pyPositionArgument = &argument
	}

	pyStateArgument,err := GetObjectArgument(ptb, client, nemoConfig.PyState, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	marketGlobalConfigArgument,err := GetObjectArgument(ptb, client, nemoConfig.MarketFactoryConfig, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	marketStateArgument,err := GetObjectArgument(ptb, client, nemoConfig.MarketState, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	clockArgument,err := GetObjectArgument(ptb, client, constant.CLOCK, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	approxPtOutCallArg := CreatePureU64CallArg(approxPtOut)
	approxPtOutArgument,err := ptb.Input(approxPtOutCallArg)
	if err != nil {
		return nil, err
	}

	minPtOutCallArg := CreatePureU64CallArg(minYtOut)
	minPtOutArgument,err := ptb.Input(minPtOutCallArg)
	if err != nil {
		return nil, err
	}

	var arguments []sui_types.Argument
	arguments = append(arguments, versionArgument, minPtOutArgument, approxPtOutArgument, *coinArgument, *oracleArgument, *pyPositionArgument, pyStateArgument, marketGlobalConfigArgument, marketStateArgument, clockArgument)
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

func GetApproxYtOutForNetSyInInternal(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoPackage, syType, pyState, marketGlobalConfig, marketState string, netSyIn, minYtOut uint64, oracleArgument *sui_types.Argument) (*sui_types.Argument,error){
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoPackage)
	if err != nil {
		return nil, err
	}

	moduleName := "offchain"
	functionName := "get_approx_yt_out_for_net_sy_in_internal"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)
	syStructTag, err := GetStructTag(syType)
	if err != nil {
		return nil, err
	}
	syTypeTag := move_types.TypeTag{
		Struct: syStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, syTypeTag)

	pyStateArgument,err := GetObjectArgument(ptb, client, pyState, false, nemoPackage, moduleName, functionName)

	marketGlobalConfigArgument,err := GetObjectArgument(ptb, client, marketGlobalConfig, false, nemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	marketStateArgument,err := GetObjectArgument(ptb, client, marketState, false, nemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	clockArgument,err := GetObjectArgument(ptb, client, constant.CLOCK, false, nemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	netSyInArg := CreatePureU64CallArg(netSyIn)
	netSyInArgument,err := ptb.Pure(netSyInArg)
	if err != nil {
		return nil, err
	}
	minYtOutArg := CreatePureU64CallArg(minYtOut)
	minYtOutArgument,err := ptb.Pure(minYtOutArg)
	if err != nil {
		return nil, err
	}

	var arguments []sui_types.Argument
	arguments = append(arguments, netSyInArgument, minYtOutArgument, *oracleArgument, pyStateArgument, pyStateArgument, marketStateArgument, marketGlobalConfigArgument, clockArgument)
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

func MintPy(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, coinArgument, priceOracleArgument, pyPositionArgument *sui_types.Argument) (*sui_types.Argument,error) {
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.NemoContract)
	if err != nil {
		return nil, err
	}

	moduleName := "yield_factory"
	functionName := "mint_py"
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

	versionArgument,err := GetObjectArgument(ptb, client, nemoConfig.Version, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	pyStateArgument,err := GetObjectArgument(ptb, client, nemoConfig.PyState, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	yieldFactoryConfigArgument,err := GetObjectArgument(ptb, client, nemoConfig.YieldFactoryConfig, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	clockArgument,err := GetObjectArgument(ptb, client, constant.CLOCK, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	var arguments []sui_types.Argument
	arguments = append(arguments, versionArgument, *coinArgument, *priceOracleArgument, *pyPositionArgument, pyStateArgument, yieldFactoryConfigArgument, clockArgument)
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

func RedeemPy(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig, amountIn uint64, priceOracleArgument, pyPositionArgument *sui_types.Argument) (*sui_types.Argument,error) {
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.NemoContract)
	if err != nil {
		return nil, err
	}

	moduleName := "yield_factory"
	functionName := "redeem_py"
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

	versionArgument,err := GetObjectArgument(ptb, client, nemoConfig.Version, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	pyStateArgument,err := GetObjectArgument(ptb, client, nemoConfig.PyState, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	yieldFactoryConfigArgument,err := GetObjectArgument(ptb, client, nemoConfig.YieldFactoryConfig, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	clockArgument,err := GetObjectArgument(ptb, client, constant.CLOCK, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	ptInArgument,err := ptb.Pure(amountIn)
	if err != nil {
		return nil, err
	}

	ytInArgument,err := ptb.Pure(amountIn)
	if err != nil {
		return nil, err
	}

	var arguments []sui_types.Argument
	arguments = append(arguments, versionArgument, ptInArgument, ytInArgument, *priceOracleArgument, *pyPositionArgument, pyStateArgument, yieldFactoryConfigArgument, clockArgument)
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