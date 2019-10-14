package util

import (
	"math/big"
	"strconv"
)

//进制转换
func Hex2Int(s string) (uint64, error) {
	n, err := strconv.ParseUint(s, 16, 32)
	return n, err
}

func Hex2BigInt(hex string) *big.Int {
	n := new(big.Int)
	n, _ = n.SetString(hex, 16)

	return n
}
