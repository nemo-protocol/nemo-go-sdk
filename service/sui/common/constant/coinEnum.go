package constant

import (
	"github.com/nemo-protocol/nemo-go-sdk/utils"
	"strings"
)

var (
	GASCOINTYPE = "0x2::sui::SUI"
	SCALLOP     = "Scallop"
	Nemo     = "Nemo"

	SCALLOPSSUI   = "0xaafc4f740de0dd0dde642a31148fb94517087052f19afb0f7bed1dc41a50c77b::scallop_sui::SCALLOP_SUI"
	SCALLOPSCA    = "0x5ca17430c1d046fae9edeaa8fd76c7b4193a00d764a0ecfa9418d733ad27bc1e::scallop_sca::SCALLOP_SCA"
	SCALLOPWUSDC  = "0xad4d71551d31092230db1fd482008ea42867dbf27b286e9c70a79d2a6191d58d::scallop_wormhole_usdc::SCALLOP_WORMHOLE_USDC"
	SCALLOPWUSDT  = "0xe6e5a012ec20a49a3d1d57bd2b67140b96cd4d3400b9d79e541f7bdbab661f95::scallop_wormhole_usdt::SCALLOP_WORMHOLE_USDT"
	SCALLOPDEEP   = "0xeb7a05a3224837c5e5503575aed0be73c091d1ce5e43aa3c3e716e0ae614608f::scallop_deep::SCALLOP_DEEP"
	SCALLOPAFSUI  = "0x00671b1fa2a124f5be8bdae8b91ee711462c5d9e31bda232e70fd9607b523c88::scallop_af_sui::SCALLOP_AF_SUI"
	SCALLOPUSDC   = "0x854950aa624b1df59fe64e630b2ba7c550642e9342267a33061d59fb31582da5::scallop_usdc::SCALLOP_USDC"
	SCALLOPSBUSDT = "0xb1d7df34829d1513b73ba17cb7ad90c88d1e104bb65ab8f62f13e0cc103783d3::scallop_sb_usdt::SCALLOP_SB_USDT"
	SCALLOPSBETH  = "0xb14f82d8506d139eacef109688d1b71e7236bcce9b2c0ad526abcd6aa5be7de0::scallop_sb_eth::SCALLOP_SB_ETH"
	VSUI          = "0x549e8b69270defbfafd4f94e17ec44cdbdd99820b33bda2278dea3b9a32d3f55::cert::CERT"
	SPRINGSUI     = "0x83556891f4a0f233ce7b05cfe7f957d4020492a34f5405b2cb9377d060bef4bf::spring_sui::SPRING_SUI"
	AFSUI         = "0xf325ce1300e8dac124071d3152c5c5ee6174914f8bc2161e88329cf579246efc::afsui::AFSUI"
	HASUI         = "0xbde4ba4c2e274a60ce15c1cfff9e5c42e41654ac8b6d906a57efa4bd3c29f47d::hasui::HASUI"
	STSUI         = "0xd1b72982e40348d069bb1ff701e634c117bb5f741f44dff91e472d3b01461e55::stsui::STSUI"
	STSBUCK       = "0xd01d27939064d79e4ae1179cd11cfeeff23943f32b1a842ea1a1e15a0045d77d::st_sbuck::ST_SBUCK"
	SUPERSUI      = "0x790f258062909e3a0ffc78b3c53ac2f62d7084c3bab95644bdeb05add7250001::super_sui::SUPER_SUI"
	HAWAL         = "0x8b4d553839b219c3fd47608a0cc3d5fcc572cb25d41b7df3833208586a8d2470::hawal::HAWAL"
	WWAL          = "0xb1b0650a8862e30e3f604fd6c5838bc25464b8d3d827fbd58af7cb9685b832bf::wwal::WWAL"

	LPTOKENHASUI = "0x828b452d2aa239d48e4120c24f4a59f451b8cd8ac76706129f4ac3bd78ac8809::lp_token::LP_TOKEN"
	LPTOKENVSUI  = "0xb490d6fa9ead588a9d72da07a02914da42f6b5b1339b8118a90011a42b67a44f::lp_token::LP_TOKEN"
	LPTOKENAFSUI = "0x0c8a5fcbe32b9fc88fe1d758d33dd32586143998f68656f43f3a6ced95ea4dc3::lp_token::LP_TOKEN"

	SUI   = "0x0000000000000000000000000000000000000000000000000000000000000002::sui::SUI"
	SCA   = "0x7016aae72cfc67f2fadf55769c0a7dd54291a583b63051a5ed71081cce836ac6::sca::SCA"
	DEEP  = "0xdeeb7a4662eec9f2f3def03fb937a663dddaa2e215b8078a284d026b7946c270::deep::DEEP"
	USDC  = "0xdba34672e30cb065b1f93e3ab55318768fd6fef66c15942c9f7cb846e2f900e7::usdc::USDC"
	USDT  = "0x375f70cf2ae4c00bf37117d0c85a2c71545e6ee05c4a5c7d282cd66a4504b068::usdt::USDT"
	WUSDT = "0xc060006111016b8a020ad5b33834984a437aaa7d3c74c18e09a95d48aceab08c::coin::COIN"
	WUSDC = "0x5d4b302506645c37ff133b98c4b50a5ae14841659738d6d733d59d0d217a93bf::coin::COIN"
	BUCK  = "0xce7ff77a83ea0cb6fd39bd8748e2ec89a3f41e8efdc3f4eb123e0ca37b184db2::buck::BUCK"

	SPRING_PROVIDER = "SpringSui"
	WINTER_PROVIDER = "Winter"
	LIQUID_AGENTS = "Liquid Agents"

)

func IsGasCoinType(coinType string) bool {
	return coinType == GASCOINTYPE
}

func IsScallopCoin(coinType string) bool {
	sCoinList := []string{
		SCALLOPSSUI, SCALLOPSCA, SCALLOPWUSDC, SCALLOPWUSDT, SCALLOPDEEP, SCALLOPAFSUI, SCALLOPUSDC, SCALLOPSBUSDT, SCALLOPSBETH,
	}
	return utils.Contains(sCoinList, coinType) || strings.Contains(coinType, "::scallop")
}

func IsVSui(coinType string) bool {
	return VSUI == coinType
}

func IsSpringSui(coinType string) bool {
	return SPRINGSUI == coinType
}

func IsAfSui(coinType string) bool {
	return AFSUI == coinType
}

func IsHaSui(coinType string) bool {
	return HASUI == coinType
}

func IsSui(coinType string) bool {
	return coinType == GASCOINTYPE || coinType == SUI
}

func IsBuck(coinType string) bool {
	return coinType == BUCK
}

func IsStSui(coinType string) bool {
	return coinType == STSUI
}

func IsLpTokenHaSui(coinType string) bool{
	return coinType == LPTOKENHASUI
}

func IsLpTokenAfSui(coinType string) bool{
	return coinType == LPTOKENAFSUI
}

func IsLpTokenVSui(coinType string) bool{
	return coinType == LPTOKENVSUI
}

func IsStsBuck(coinType string) bool{
	return coinType == STSBUCK
}

func IsSuperSui(coinType string) bool{
	return coinType == SUPERSUI
}

func IsHaWal(coinType string) bool{
	return coinType == HAWAL
}

func IsWWal(coinType string) bool{
	return coinType == WWAL
}

func IsSpringCoin(provider string) bool{
	return provider == SPRING_PROVIDER || provider == LIQUID_AGENTS
}

func IsWinterCoin(provider string) bool{
	return provider == WINTER_PROVIDER
}

func IsDeep(coinType string) bool {
	return coinType == DEEP
}

func IsSca(coinType string) bool {
	return coinType == SCA
}

func IsUsdc(coinType string) bool {
	return coinType == USDC
}

func IsUsdt(coinType string) bool {
	return coinType == USDT
}

func IsWUsdt(coinType string) bool {
	return coinType == WUSDT
}

func IsWUsdc(coinType string) bool {
	return coinType == WUSDC
}
