package utils

import "math/big"

func ReadUint128ToBigInt(byteSlice []byte) *big.Int {
	if len(byteSlice) < 16 {
		return big.NewInt(0)
	}

	result := new(big.Int)

	reversed := make([]byte, len(byteSlice))
	for i := 0; i < len(byteSlice); i++ {
		reversed[len(byteSlice)-1-i] = byteSlice[i]
	}

	return result.SetBytes(reversed)
}
