package models

import (
	"encoding/json"
	"fmt"
	"github.com/nemo-protocol/nemo-go-sdk/utils"
	"strconv"
)

type NemoConfig struct {
	CoinPrice                string            `json:"coinPrice"`
	CoinType                 string            `json:"coinType"`
	SyCoinType               string            `json:"syCoinType"`
	UnderlyingCoinType       string            `json:"underlyingCoinType"`
	UnderlyingCoinPrice      string            `json:"UnderlyingCoinPrice"`
	UnderlyingApy            string            `json:"underlyingApy"`
	Decimal                  uint64            `json:"decimal"`
	ConversionRate           string            `json:"conversionRate"`
	PyState                  string            `json:"pyState"`
	Version                  string            `json:"version"`
	YieldFactoryConfig       string            `json:"yieldFactoryConfig"`
	MarketFactoryConfig      string            `json:"marketFactoryConfig"`
	MarketState              string            `json:"marketState"`
	SyState                  string            `json:"syState"`
	PyStore                  string            `json:"pyStore"`
	PriceOracle              string            `json:"priceOracle"`
	HaedalStakeing           string            `json:"haedalStakeing"`
	NativePool               string            `json:"nativePool"`
	Metadata                 string            `json:"metadata"`
	ProviderMarket           string            `json:"providerMarket"`
	ProviderVersion          string            `json:"providerVersion"`
	LstInfo                  string            `json:"lstInfo"`
	NemoContract             string            `json:"nemoContract"`
	NemoContractList         []string          `json:"nemoContractList"`
	ProviderProtocol         string            `json:"providerProtocol"`
	OraclePackage            string            `json:"oraclePackageId"`
	OracleTicket             string            `json:"oracleTicket"`
	OracleVoucherPackage     string            `json:"oracleVoucherPackageId"`
	SwapFeeForLpHolder       string            `json:"swapFeeForLpHolder"`
	YieldTokenType           string            `json:"yieldTokenType"`
	WinterStaking            string            `json:"winterStaking"`
	CacheContractPackageInfo map[string]string `json:"cacheContractPackageInfo"`
	Incentives               []Incentives      `json:"incentives"`
	VaultId                  string            `json:"vaultId"`
	MmtOracle                string            `json:"mmtOracle"`
	PoolId                   string            `json:"poolId"`
	StableType               string            `json:"stableType"`
}

type NemoConfigInfo struct {
	CoinPrice            string       `json:"coinPrice"`
	CoinType             string       `json:"coinType"`
	SyCoinType           string       `json:"syCoinType"`
	UnderlyingCoinType   string       `json:"underlyingCoinType"`
	UnderlyingApy        string       `json:"underlyingApy"`
	Decimal              string       `json:"decimal"`
	ConversionRate       string       `json:"conversionRate"`
	PyState              string       `json:"pyStateId"`
	Version              string       `json:"version"`
	YieldFactoryConfig   string       `json:"yieldFactoryConfigId"`
	MarketFactoryConfig  string       `json:"marketFactoryConfigId"`
	MarketState          string       `json:"marketStateId"`
	SyState              string       `json:"syStateId"`
	PyStore              string       `json:"pyStoreId"`
	PriceOracle          string       `json:"priceOracleConfigId"`
	HaedalStakeing       string       `json:"haedalStakeingId"`
	NativePool           string       `json:"nativePool"`
	Metadata             string       `json:"metadataId"`
	ProviderMarket       string       `json:"providerMarket"`
	ProviderVersion      string       `json:"providerVersion"`
	NemoContract         string       `json:"nemoContractId"`
	NemoContractList     []string     `json:"nemoContractIdList"`
	ProviderProtocol     string       `json:"underlyingProtocol"`
	OraclePackage        string       `json:"oraclePackageId"`
	OracleTicket         string       `json:"oracleTicket"`
	OracleVoucherPackage string       `json:"oracleVoucherPackageId"`
	SwapFeeForLpHolder   string       `json:"swapFeeForLpHolder"`
	YieldTokenType       string       `json:"yieldTokenType"`
	LstInfo              string       `json:"lstInfo"`
	WinterStaking        string       `json:"winterStaking"`
	Incentives           []Incentives `json:"incentives"`
	VaultId              string       `json:"vaultId"`
	MmtOracle            string       `json:"mmtOracle"`
	PoolId               string            `json:"poolId"`
}

type NemoInfoResponse struct {
	NemoConfigInfo NemoConfigInfo `json:"data"`
}

type NemoPage struct {
	Id       string `json:"id"`
	CoinType string `json:"coinType"`
}

type NemoPageResponse struct {
	NemoPage []NemoPage `json:"data"`
}

func InitConfig() []NemoConfig {
	/**
	scallopSui := &NemoConfig{
		CoinType:            "0xaafc4f740de0dd0dde642a31148fb94517087052f19afb0f7bed1dc41a50c77b::scallop_sui::SCALLOP_SUI",
		SyCoinType:          "0x53a8c1ffcdac36d993ce3c454d001eca57224541d1953d827ef96ac6d7f8142e::sSUI::SSUI",
		UnderlyingCoinType:  "0x2::sui::SUI",
		Decimal:             9,
		ConversionRate:      "1.0580698767",
		PyState:             "0xc6840365f500bee8732a3a256344a11343936b864c144b7e9de5bb8c54224fbe",
		Version:             "0x467ce1e4351ff0e0dcdd7ab98ded8a0573c6acfb3ce10cb33bcfb2db06caca88",
		YieldFactoryConfig:  "0x56392aa1f849e901f4b9d0d313cdd02ab74543bbecb97708653051ab680cf281",
		MarketFactoryConfig: "0x4a8d13937be10f97e450d1b8eb5846b749f9d3f470243b6cfa660e3d75b1fc49",
		MarketState:         "0x7472959314b24ebfbd4da49cc36abb3da29f722746019c692407aaf6b47e9a08",
		SyState:             "0xccd3898005a269c1e9074fe28bca2ff46784e8ee7c13b576862d9758266c3a4d",
		PyStore:             "0x0f589f1f1937b39cc51cd04b66dffe69ff6358693a5014dac75d6621730dbd9b",
		PriceOracle:         "0xb9cc723bf7494325be2f3333a3fb72f46d53abe3603e3f326fc761287850db0e",
		NemoContract:        "0x2b71664477755b90f9fb71c9c944d5d0d3832fec969260e3f18efc7d855f57c4",
		NemoContractList: []string{
			"0x2b71664477755b90f9fb71c9c944d5d0d3832fec969260e3f18efc7d855f57c4"},
		ProviderProtocol:     "Scallop",
		OraclePackage:        "0xee1ff66985a76b2c0170935fb29144b4007827ed2c4f3d6a1189578afb92bcdd",
		OracleTicket:         "0x0fa9dc987f71878b91d389c145aab67f462744b695054578ca4ae4d6ced01099",
		OracleVoucherPackage: "0x8783841625738f73a6b0085f5dad270b4b0bd2e5cdb278dc95201e45bd1a332b",
	}
	*/
	url := "https://api.nemoprotocol.com/api/v1/market/coinInfo?isShowExpiry=0"
	pageByte, err := utils.SendGetRpc(url)
	if err != nil {
		return nil
	}

	response := NemoPageResponse{}
	_ = json.Unmarshal(pageByte, &response)

	infoUrl := "https://api.nemoprotocol.com/api/v1/market/config/detail?id=%v"
	infoList := make([]NemoConfig, 0)
	for _, v := range response.NemoPage {
		infoByte, err := utils.SendGetRpc(fmt.Sprintf(infoUrl, v.Id))
		if err != nil {
			return nil
		}
		info := NemoInfoResponse{}
		err = json.Unmarshal(infoByte, &info)
		if err != nil {
			continue
		}
		innerInfo := FormatStruct(info.NemoConfigInfo)
		infoList = append(infoList, innerInfo)
	}
	return infoList
}

func FormatStruct(resInfo NemoConfigInfo) NemoConfig {
	innerInfo := NemoConfig{}
	innerInfo.CoinType = resInfo.CoinType
	innerInfo.SyCoinType = resInfo.SyCoinType
	innerInfo.UnderlyingCoinType = resInfo.UnderlyingCoinType
	decimal, _ := strconv.ParseInt(resInfo.Decimal, 10, 64)
	innerInfo.Decimal = uint64(decimal)
	innerInfo.ConversionRate = resInfo.ConversionRate
	innerInfo.PyState = resInfo.PyState
	innerInfo.Version = resInfo.Version
	innerInfo.YieldFactoryConfig = resInfo.YieldFactoryConfig
	innerInfo.MarketFactoryConfig = resInfo.MarketFactoryConfig
	innerInfo.MarketState = resInfo.MarketState
	innerInfo.SyState = resInfo.SyState
	innerInfo.PyStore = resInfo.PyStore
	innerInfo.PriceOracle = resInfo.PriceOracle
	innerInfo.HaedalStakeing = resInfo.HaedalStakeing
	innerInfo.NativePool = resInfo.NativePool
	innerInfo.Metadata = resInfo.Metadata
	innerInfo.ProviderMarket = resInfo.ProviderMarket
	innerInfo.ProviderVersion = resInfo.ProviderVersion
	innerInfo.NemoContract = resInfo.NemoContract
	innerInfo.NemoContractList = resInfo.NemoContractList
	innerInfo.ProviderProtocol = resInfo.ProviderProtocol
	innerInfo.OraclePackage = resInfo.OraclePackage
	innerInfo.OracleTicket = resInfo.OracleTicket
	innerInfo.OracleVoucherPackage = resInfo.OracleVoucherPackage
	innerInfo.CoinPrice = resInfo.CoinPrice
	innerInfo.UnderlyingApy = resInfo.UnderlyingApy
	innerInfo.SwapFeeForLpHolder = resInfo.SwapFeeForLpHolder
	innerInfo.YieldTokenType = resInfo.YieldTokenType
	innerInfo.LstInfo = resInfo.LstInfo
	innerInfo.WinterStaking = resInfo.WinterStaking
	innerInfo.Incentives = resInfo.Incentives
	innerInfo.MmtOracle = resInfo.MmtOracle
	innerInfo.VaultId = resInfo.VaultId
	innerInfo.PoolId = resInfo.PoolId
	return innerInfo
}

type NemoVaultConfig struct {
	VaultContract    string `json:"vaultContract"`
	VaultId          string `json:"vaultId"`
	PoolId           string `json:"poolId"`
	LeftCoinType     string `json:"leftCoinType"`
	LeftCoinPrice    string `json:"leftCoinPrice"`
	LeftCoinDecimal  string `json:"leftCoinDecimal"`
	RightCoinType    string `json:"rightCoinType"`
	RightCoinPrice   string `json:"rightCoinPrice"`
	RightCoinDecimal string `json:"rightCoinDecimal"`
	VaultType        string `json:"vaultType"`
	StableType       string `json:"stableType"`
}

type VaultPageResponse struct {
	VaultPage []NemoVaultConfig `json:"data"`
}

func InitVaultConfig() []NemoVaultConfig {
	url := "https://api.nemoprotocol.com/api/v1/market/user/vaultInfo?"
	pageByte, err := utils.SendGetRpc(url)
	if err != nil {
		return nil
	}

	response := VaultPageResponse{}
	_ = json.Unmarshal(pageByte, &response)
	return response.VaultPage
}
