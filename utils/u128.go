package utils

import (
	"fmt"
	"math/big"
)

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

func InterfaceSliceToBytes(data interface{}) ([]byte, error) {
	// 首先尝试直接转换为[]uint8
	if byteSlice, ok := data.([]uint8); ok {
		return byteSlice, nil
	}

	// 如果不是[]uint8，尝试其他类型转换
	slice, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a slice")
	}

	bytes := make([]byte, len(slice))
	for i, v := range slice {
		switch val := v.(type) {
		case uint8:
			bytes[i] = val
		case int64:
			if val < 0 || val > 255 {
				return nil, fmt.Errorf("value out of byte range")
			}
			bytes[i] = byte(val)
		case float64:
			if val < 0 || val > 255 || val != float64(int(val)) {
				return nil, fmt.Errorf("invalid float value")
			}
			bytes[i] = byte(val)
		default:
			return nil, fmt.Errorf("unsupported type: %T", v)
		}
	}
	return bytes, nil
}
