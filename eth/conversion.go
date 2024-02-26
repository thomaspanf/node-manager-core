package eth

import (
	"math/big"
	"strconv"
)

// Conversion factors
const (
	// Amount of wei in 1 ETH
	WeiPerEth float64 = 1e18

	// Amount of wei in 1 gwei
	WeiPerGwei float64 = 1e9

	// Amount of gwei in 1 ETH
	GweiPerEth float64 = WeiPerEth / WeiPerGwei
)

var (
	weiPerEthFloat  *big.Float = big.NewFloat(WeiPerEth)
	WeiPerGweiFloat *big.Float = big.NewFloat(WeiPerGwei)
)

// Convert a wei amount (a native uint256 value on the execution layer) to a floating-point ETH amount
func WeiToEth(wei *big.Int) float64 {
	var weiFloat big.Float
	var eth big.Float
	weiFloat.SetInt(wei)
	eth.Quo(&weiFloat, weiPerEthFloat)
	eth64, _ := eth.Float64()
	return eth64
}

// Convert a floating-point ETH amount to a wei amount (a native uint256 value on the execution layer)
func EthToWei(eth float64) *big.Int {
	var ethFloat big.Float
	var weiFloat big.Float
	var wei big.Int
	ethFloat.SetString(strconv.FormatFloat(eth, 'f', -1, 64))
	weiFloat.Mul(&ethFloat, weiPerEthFloat)
	weiFloat.Int(&wei)
	return &wei
}

// Convert a wei amount (a native uint256 value on the execution layer) to a floating-point gwei amount
func WeiToGwei(wei *big.Int) float64 {
	var weiFloat big.Float
	var gwei big.Float
	weiFloat.SetInt(wei)
	gwei.Quo(&weiFloat, WeiPerGweiFloat)
	gwei64, _ := gwei.Float64()
	return gwei64
}

// Convert a floating-point gwei amount to a wei amount (a native uint256 value on the execution layer)
func GweiToWei(gwei float64) *big.Int {
	var gweiFloat big.Float
	var weiFloat big.Float
	var wei big.Int
	gweiFloat.SetString(strconv.FormatFloat(gwei, 'f', -1, 64))
	weiFloat.Mul(&gweiFloat, WeiPerGweiFloat)
	weiFloat.Int(&wei)
	return &wei
}

// Convert a floating-point ETH amount to a floating-point gwei amount
func EthToGwei(eth float64) float64 {
	return eth * GweiPerEth
}

// Convert a floating-point gwei amount to a floating-point ETH amount
func GweiToEth(gwei float64) float64 {
	return gwei / GweiPerEth
}
