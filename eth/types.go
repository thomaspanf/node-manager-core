package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	batch "github.com/rocket-pool/batch-query"
)

// Information about a transaction's simulation
type SimulationResult struct {
	// True if the transaction was simulated, false if it was not
	IsSimulated bool `json:"isSimulated"`

	// The raw amount of gas, in gwei, the transaction took during simulation
	EstimatedGasLimit uint64 `json:"estimatedGasLimit"`

	// A safe gas limit to use for the transaction, in gwei, with a reasonable safety buffer to account for variances from simulation
	SafeGasLimit uint64 `json:"safeGasLimit"`

	// Any error / revert that occurred during simulation, indicating the transaction may fail if submitted
	SimulationError string `json:"simulationError"`
}

// Information of a candidate transaction
type TransactionInfo struct {
	// The transaction's data
	Data []byte `json:"data"`

	// The address to send the transaction to
	To common.Address `json:"to"`

	// The ETH value, in wei, to send along with the transaction
	Value *big.Int `json:"value"`

	// Info about the transaction's simulation
	SimulationResult SimulationResult `json:"simulationResult"`
}

// Information for submitting a candidate transaction to the network
type TransactionSubmission struct {
	// The transaction info
	TxInfo *TransactionInfo `json:"txInfo"`

	// The gas limit to use when submitting this transaction
	GasLimit uint64 `json:"gasLimit"`
}

// Represents structs that can have their values queried during a multicall
type IQueryable interface {
	// Adds the struct's values to the provided multicall query before it runs
	AddToQuery(mc *batch.MultiCaller)
}
