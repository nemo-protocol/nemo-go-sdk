package constant

var (
	GASCOINTYPE = "0x2::sui::SUI"
)

func IsGasCoinType(coinType string) bool{
	return coinType == GASCOINTYPE
}