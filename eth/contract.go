package eth

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	batch "github.com/rocket-pool/batch-query"
)

// Contract is a wrapper for go-ethereum bound contracts
type Contract struct {
	// A human-readable name of the contract
	Name string

	// The contract's address
	Address common.Address

	// The contract's ABI
	ABI *abi.ABI

	// The underlying bound contract
	ContractImpl *bind.BoundContract
}

// This is a helper for adding calls to multicall
func (c *Contract) AddCall(mc *batch.MultiCaller, output any, method string, args ...any) {
	mc.AddCall(c.Address, c.ABI, output, method, args...)
}
