package config

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// A collection of network-specific resources and getters for them
type NetworkResources struct {
	// The network being used
	network Network

	// The address of the multicall contract
	MulticallAddress common.Address

	// The chain ID for the current network
	ChainID uint

	// The URL for transaction monitoring on the network's chain explorer
	TxWatchUrl string
}

// Creates a new resource collection for the given network
func NewResources(network Network) *NetworkResources {
	// Mainnet
	mainnetResources := &NetworkResources{
		network:          network,
		MulticallAddress: common.HexToAddress("0x5BA1e12693Dc8F9c48aAD8770482f4739bEeD696"),
		ChainID:          1,
		TxWatchUrl:       "https://etherscan.io/tx",
	}

	// Holesky
	holeskyResources := &NetworkResources{
		network:          network,
		MulticallAddress: common.HexToAddress("0x0540b786f03c9491f3a2ab4b0e3ae4ecd4f63ce7"),
		ChainID:          17000,
		TxWatchUrl:       "https://holesky.etherscan.io/tx",
	}

	switch network {
	case Network_Mainnet:
		return mainnetResources
	case Network_Holesky:
		return holeskyResources
	}

	panic(fmt.Sprintf("network %s is not supported", network))
}
