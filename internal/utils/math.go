package utils

import (
	"math/big"
)

// WeiToEther converts Wei (*big.Int) into Ether (*big.Float) securely avoiding precision loss in floats.
func WeiToEther(wei *big.Int) *big.Float {
	if wei == nil {
		return big.NewFloat(0)
	}
	f := new(big.Float).SetInt(wei)
	// 1 Ether = 10^18 Wei
	divisor := new(big.Float).SetFloat64(1e18)
	return new(big.Float).Quo(f, divisor)
}
