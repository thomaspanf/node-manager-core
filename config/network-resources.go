package config

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// A collection of network-specific resources and getters for them
type NetworkResources struct {
	// The Network being used
	Network Network

	// The actual name of the underlying Ethereum network, passed into the clients
	EthNetworkName string

	// The chain ID for the current network
	ChainID uint

	// The genesis fork version for the network according to the Beacon config for the network
	GenesisForkVersion []byte

	// The address of the multicall contract
	MulticallAddress common.Address

	// The BalanceChecker contract address
	BalanceBatcherAddress common.Address

	// The URL for transaction monitoring on the network's chain explorer
	TxWatchUrl string

	// The FlashBots Protect RPC endpoint
	FlashbotsProtectUrl string
}

// Creates a new resource collection for the given network
func NewResources(network Network) *NetworkResources {
	// Mainnet
	mainnetResources := &NetworkResources{
		Network:               Network_Mainnet,
		EthNetworkName:        string(Network_Mainnet),
		ChainID:               1,
		GenesisForkVersion:    common.FromHex("0x00000000"), // https://github.com/eth-clients/eth2-networks/tree/master/shared/mainnet#genesis-information
		MulticallAddress:      common.HexToAddress("0x5BA1e12693Dc8F9c48aAD8770482f4739bEeD696"),
		BalanceBatcherAddress: common.HexToAddress("0xb1f8e55c7f64d203c1400b9d8555d050f94adf39"),
		TxWatchUrl:            "https://etherscan.io/tx",
		FlashbotsProtectUrl:   "https://rpc.flashbots.net/",
	}

	// Holesky
	holeskyResources := &NetworkResources{
		Network:               Network_Holesky,
		EthNetworkName:        string(Network_Holesky),
		ChainID:               17000,
		GenesisForkVersion:    common.FromHex("0x01017000"), // https://github.com/eth-clients/holesky
		MulticallAddress:      common.HexToAddress("0x0540b786f03c9491f3a2ab4b0e3ae4ecd4f63ce7"),
		BalanceBatcherAddress: common.HexToAddress("0xfAa2e7C84eD801dd9D27Ac1ed957274530796140"),
		TxWatchUrl:            "https://holesky.etherscan.io/tx",
		FlashbotsProtectUrl:   "",
	}

	switch network {
	case Network_Mainnet:
		return mainnetResources
	case Network_Holesky:
		return holeskyResources
	}

	panic(fmt.Sprintf("network %s is not supported", network))
}
