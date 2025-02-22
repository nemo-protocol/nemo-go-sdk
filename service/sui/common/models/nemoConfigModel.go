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
		SyCoinType:          "0xb31beb8a12c3814cc402ae5eb9ab9405b5633942644b8aaa7c069f00d7cd9ec2::sSUI::SSUI",
		UnderlyingCoinType:  "0x2::sui::SUI",
		Decimal:             9,
		ConversionRate:      "1.0580698767",
		PyState:             "0xdf295eaed18181c9d8d98547de3953aac7a08939de63d5bb0a2f0eaafba48acc",
		Version:             "0x4000b5c20e70358a42ae45421c96d2f110817d75b80df30dad5b5d4f1fdad6af",
		YieldFactoryConfig:  "0xf274dd733e159d3b7a671c4304bcc731def8c206cd62f66518feb9a941b864fe",
		MarketFactoryConfig: "0x87dd91ea90c44e9410c0b899969ccfeb04b4bccb257c75b1865da4b615bba602",
		MarketState:         "0xe9c3ef2ac1cf92b967a69663815498b20af6da1d46edb5d328f4a039730e3747",
		SyState:             "0x6f67f42a3c0f59d5d7a8cd11f9bffec9ae2716282a2da287df5244c79463ecba",
		PyStore:             "0xc0e5fcc424cb9b03c54c5d5cec27766efd194a6c8f50f6720ea8e9c7e710047a",
		PriceOracle:         "0xde87ed5462249dd6928269a37f0fda300971cad8586f1b2dfaf69a2ad070d63a",
		NemoContract:        "0x0dd9407d07f05c7367e29fc485f3ea8e14701a63d8b0a11c6d2cd40377536586",
		NemoContractList: []string{
			"0xe4daf34d4aa8eb89fe9676c8dbec6a0675c53cb8f98fac9b17357f185e992a31",
			"0xa09a5e35bd3be6993a978c674061a12b6796a7d5e696aaf068d1d402c59ac654",
			"0x0dd9407d07f05c7367e29fc485f3ea8e14701a63d8b0a11c6d2cd40377536586"},
		ProviderProtocol:     "Scallop",
		OraclePackage:        "0xc167c8782d73bf1c5777c2b005b730e53b21cc277136c7c3b94cca625dc6cf3f",
		OracleTicket:         "0xcd50df05daa0b6521a326ad534269b34903f4895c370f31eb37be2d1a68e317a",
		OracleVoucherPackage: "0x85f6136e8af827d5cdfbf07927698f5ee035ffa34c018148ee0737fb8e43b7aa",
	}
	return scallopSui
}
