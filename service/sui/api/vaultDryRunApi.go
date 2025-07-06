package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/move_types"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/fardream/go-bcs/bcs"
	"github.com/nemo-protocol/nemo-go-sdk/service/sui/common/constant"
	"github.com/nemo-protocol/nemo-go-sdk/service/sui/common/models"
	"math"
	"strconv"
	"strings"
)

const (
	VAULT_VERSION_ID = "0xf6c198001167e74b6986c8f7619eb168a9c7d6f20ea6a55f6e2c506e7608f710"
	MMT_CLMM_VERSION_ID = "0x2375a0b1ec12010aaea3b2545acfa2ad34cfbba03ce4b59f4c39e1e25eed1b2a"
)

func DryRunVaultWithdraw(client *client.Client, address string, vaultConfig *models.NemoVaultConfig) ([]models.CoinInfo, error){
	ptb := sui_types.NewProgrammableTransactionBuilder()

	vaultPackageId, err := sui_types.NewObjectIdFromHex(vaultConfig.VaultContract)
	if err != nil {
		return nil, err
	}

	leftCoinTypeStructTag, err := GetStructTag(vaultConfig.LeftCoinType)
	if err != nil {
		return nil, err
	}
	leftCoinTypeTag := move_types.TypeTag{
		Struct: leftCoinTypeStructTag,
	}

	rightCoinTypeStructTag, err := GetStructTag(vaultConfig.RightCoinType)
	if err != nil {
		return nil, err
	}
	rightCoinTypeTag := move_types.TypeTag{
		Struct: rightCoinTypeStructTag,
	}

	vCoinTypeStructTag, err := GetStructTag(vaultConfig.VaultType)
	if err != nil {
		return nil, err
	}
	vaultCoinTypeTag := move_types.TypeTag{
		Struct: vCoinTypeStructTag,
	}

	stableTypeStructTag, err := GetStructTag("0xe96f73400fcf04dea660a33ed66b57742e3d936e15b551f7c2f27f2ab5a9dbf1::config::Stable")
	if err != nil {
		return nil, err
	}
	stableTypeTag := move_types.TypeTag{
		Struct: stableTypeStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, leftCoinTypeTag, rightCoinTypeTag, vaultCoinTypeTag, stableTypeTag)

	moduleName := "withdraw"
	functionName := "withdraw"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	sd, err := sui_types.NewAddressFromHex(address)
	if err != nil {
		return nil, err
	}

	coinJson, err := client.GetCoins(context.Background(), *sd, &vaultConfig.VaultType, nil, 100)
	if err != nil {
		return nil, err
	}

	var coins CoinPage
	b, _ := json.Marshal(coinJson)
	err = json.Unmarshal(b, &coins)

	if err != nil {
		return nil, err
	}

	if len(coins.Data) == 0{
		return nil, err
	}

	shareObjectMap := map[string]bool{
		vaultConfig.VaultId: false,
		vaultConfig.PoolId: false,
		constant.CLOCK: false,
		VAULT_VERSION_ID: false,
		MMT_CLMM_VERSION_ID: false,
	}

	objectArgMap, err := MultiGetObjectArg(client, shareObjectMap, vaultConfig.VaultContract, moduleName, functionName)
	if err != nil{
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs,
		sui_types.CallArg{Object: objectArgMap[vaultConfig.VaultId]},
		sui_types.CallArg{Object: objectArgMap[vaultConfig.PoolId]},
		sui_types.CallArg{Object: objectArgMap[constant.CLOCK]},
		sui_types.CallArg{Object: objectArgMap[VAULT_VERSION_ID]},
		sui_types.CallArg{Object: objectArgMap[MMT_CLMM_VERSION_ID]},
	)

	coinArgument, err := MergeAllCoin(ptb, client, coins.Data)
	if err != nil{
		return nil, err
	}

	amount := CreatePureU64CallArg(0)
	amountArgument,err := ptb.Input(amount)
	if err != nil {
		return nil, err
	}

	var arguments []sui_types.Argument
	for k, v := range callArgs {
		if k == 2{
			arguments = append(arguments, *coinArgument, amountArgument, amountArgument)
		}
		argument, err := ptb.Input(v)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, argument)
	}

	ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *vaultPackageId,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)
	pt := ptb.Finish()
	txKind := sui_types.TransactionKind{
		ProgrammableTransaction: &pt,
	}

	txBytes, err := bcs.Marshal(txKind)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize transaction: %w", err)
	}

	senderAddr, err := sui_types.NewAddressFromHex(address)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sender address: %w", err)
	}

	result, err := client.DevInspectTransactionBlock(context.Background(), *senderAddr, txBytes, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect transaction: %w", err)
	}
	if result.Error != nil {
		return nil, errors.New(fmt.Sprintf("%v", *result.Error))
	}

	if len(result.Results) == 0 {
		return nil, fmt.Errorf("no results returned")
	}

	if len(result.Results) == 0 {
		return nil, fmt.Errorf("empty results")
	}
	lastResult := result.Results[1]
	fmt.Printf("result:%+v",len(result.Results))
	if len(lastResult.ReturnValues) == 0 {
		return nil, fmt.Errorf("no return values")
	}

	coinInfoList := make([]models.CoinInfo, 0)
	for _,firstValue := range lastResult.ReturnValues{
		fmt.Printf("\n==firstValue:%v==\n", firstValue)
		var coin models.Coin
		if firstValueArray, ok := firstValue.([]interface{}); ok && len(firstValueArray) > 0 {
			if innerArray, ok := firstValueArray[0].([]interface{}); ok && len(innerArray) > 0 {
				byteSlice := make([]byte, len(innerArray))
				for i, v := range innerArray {
					if num, ok := v.(float64); ok {
						byteSlice[i] = byte(num)
					}
				}
				_, err = bcs.Unmarshal(byteSlice, &coin)
				if err != nil {
					return nil, err
				}
				fmt.Printf("Coin Value: %d\n", coin.Value)

				var decimal int64
				coinType := firstValueArray[1].(string)
				if strings.Contains(coinType, "0x2::sui::SUI"){
					coinType = "0x2::coin::Coin<0x0000000000000000000000000000000000000000000000000000000000000002::sui::SUI>"
				}
				if strings.Contains(coinType, vaultConfig.LeftCoinType){
					decimal,_ = strconv.ParseInt(vaultConfig.LeftCoinDecimal, 10, 64)
				}else{
					decimal,_ = strconv.ParseInt(vaultConfig.RightCoinDecimal, 10, 64)
				}
				coinInfo := models.CoinInfo{}
				coinInfo.Amount = float64(coin.Value) / math.Pow(10, float64(decimal))
				coinInfo.CoinType = coinType
				fmt.Printf("coinInfo:%v",coinInfo)
				coinInfoList = append(coinInfoList, coinInfo)
			}
		}
	}

	return coinInfoList, nil
}
