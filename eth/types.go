package eth

import (
	"fmt"
	"math/big"
	"strings"

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

// Quoted big ints
type QuotedBigInt big.Int

// Serialize the big.Int to JSON
func (i QuotedBigInt) MarshalJSON() ([]byte, error) {
	nativeInt := big.Int(i)
	return []byte("\"" + nativeInt.String() + "\""), nil
}

// Deserialize the big.Int from JSON
func (i *QuotedBigInt) UnmarshalJSON(data []byte) error {
	strippedString := strings.Trim(string(data), "\"")
	nativeInt, success := big.NewInt(0).SetString(strippedString, 0)
	if !success {
		return fmt.Errorf("%s is not a valid big integer", strippedString)
	}

	// Set value and return
	*i = QuotedBigInt(*nativeInt)
	return nil
}

// Converts the QuotedBigInt to the native type
func (i *QuotedBigInt) ToInt() *big.Int {
	return (*big.Int)(i)
}
