package models

type NemoConfig struct {
	CoinType            string   `json:"coinType"`
	SyCoinType          string   `json:"syCoinType"`
	UnderlyingCoinType  string   `json:"underlyingCoinType"`
	Decimal             uint64   `json:"decimal"`
	PyState             string   `json:"pyState"`
	Version             string   `json:"version"`
	YieldFactoryConfig  string   `json:"yieldFactoryConfig"`
	MarketFactoryConfig string   `json:"marketFactoryConfig"`
	MarketState         string   `json:"marketState"`
	SyState             string   `json:"syState"`
	PyStore             string   `json:"pyStore"`
	PriceOracle         string   `json:"priceOracle"`
	HaedalStakeing      string   `json:"haedalStakeing"`
	NativePool          string   `json:"nativePool"`
	Metadata            string   `json:"metadata"`
	ProviderMarket      string   `json:"providerMarket"`
	ProviderVersion     string   `json:"providerVersion"`
	LstInfo             string   `json:"lstInfo"`
	NemoContract        string   `json:"nemoContract"`
	NemoContractList    []string `json:"nemoContractList"`
	ProviderProtocol    string   `json:"providerProtocol"`
}

func InitConfig() *NemoConfig {
	scallopSui := &NemoConfig{
		CoinType:            "0xaafc4f740de0dd0dde642a31148fb94517087052f19afb0f7bed1dc41a50c77b::scallop_sui::SCALLOP_SUI",
		SyCoinType:          "0x36a4c63cd17d48d33e32c3796245b2d1ebe50c2898ee80e682b787fb9b6519d5::sSUI::SSUI",
		UnderlyingCoinType:  "0x2::sui::SUI",
		Decimal:             9,
		PyState:             "0xe8713fc5aefcbdf4f25fea27a48901e878d7b0a6681f44672672e6914437004c",
		Version:             "0x4000b5c20e70358a42ae45421c96d2f110817d75b80df30dad5b5d4f1fdad6af",
		YieldFactoryConfig:  "0x0f3e1b1922a2445a4ed5ec936a348cf6bfe50f829b92da0ba9ed3490ae1f1439",
		MarketFactoryConfig: "0x9bdde7b16ccaa212b80cb3ae8d644aa1c7f65fd12764ce9bc267fe28de72b54d",
		MarketState:         "0xd0b859d37963898cb96f8f635018651c8f1f5b68fa270027c61f589fdcde5ff1",
		SyState:             "0xdc04b5fffe78fae13e967c5943ea6b543637df8955afca1e89a70d0cf5a1a0c2",
		PyStore:             "0x13380f35aab79ee4a1314a05edbae50a6b6efbc03bf905a002854121680db7a4",
		PriceOracle:         "0x8dc043ba780bc9f5b4eab09c4e6d82d7af295e5c5ab32be5c27d9933fb02421b",
		NemoContract:        "0xa035d268323e40ab99ce8e4b12353bd89a63270935b4969d5bba87aa850c2b19",
		NemoContractList: []string{
			"0xbde9dd9441697413cf312a2d4e37721f38814b96d037cb90d5af10b79de1d446",
			"0xa035d268323e40ab99ce8e4b12353bd89a63270935b4969d5bba87aa850c2b19"},
		ProviderProtocol: "Scallop",
	}
	return scallopSui
}
