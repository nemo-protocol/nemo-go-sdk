package models

type NemoConfig struct {
	CoinType             string   `json:"coinType"`
	SyCoinType           string   `json:"syCoinType"`
	UnderlyingCoinType   string   `json:"underlyingCoinType"`
	Decimal              uint64   `json:"decimal"`
	ConversionRate       string   `json:"conversionRate"`
	PyState              string   `json:"pyState"`
	Version              string   `json:"version"`
	YieldFactoryConfig   string   `json:"yieldFactoryConfig"`
	MarketFactoryConfig  string   `json:"marketFactoryConfig"`
	MarketState          string   `json:"marketState"`
	SyState              string   `json:"syState"`
	PyStore              string   `json:"pyStore"`
	PriceOracle          string   `json:"priceOracle"`
	HaedalStakeing       string   `json:"haedalStakeing"`
	NativePool           string   `json:"nativePool"`
	Metadata             string   `json:"metadata"`
	ProviderMarket       string   `json:"providerMarket"`
	ProviderVersion      string   `json:"providerVersion"`
	LstInfo              string   `json:"lstInfo"`
	NemoContract         string   `json:"nemoContract"`
	NemoContractList     []string `json:"nemoContractList"`
	ProviderProtocol     string   `json:"providerProtocol"`
	OraclePackage        string   `json:"oraclePackageId"`
	OracleTicket         string   `json:"oracleTicket"`
	OracleVoucherPackage string   `json:"oracleVoucherPackageId"`
}

func InitConfig() *NemoConfig {
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
		OraclePackage:        "0x8d0145043ce10a7d95c27be80169b23f7b3be8993e63ee705982af6ad43e77d0",
		OracleTicket:         "0xa8b7319b326a6e4d1d4baaebbcbd5287fb82484081f0679e68cb5286171d3bd7",
		OracleVoucherPackage: "0x8783841625738f73a6b0085f5dad270b4b0bd2e5cdb278dc95201e45bd1a332b",
	}
	return scallopSui
}
