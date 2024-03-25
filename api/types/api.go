package types

import (
	"github.com/rocket-pool/node-manager-core/eth"
)

type ApiResponse[Data any] struct {
	Data *Data `json:"data"`
}

type SuccessData struct {
}

type DataBatch[DataType any] struct {
	Batch []DataType `json:"batch"`
}

type TxInfoData struct {
	TxInfo *eth.TransactionInfo `json:"txInfo"`
}

type BatchTxInfoData struct {
	TxInfos []*eth.TransactionInfo `json:"txInfos"`
}

// ResponseStatus is used to signify the status of an API request's result.
type ResponseStatus int

const (
	// Unknown (default value)
	ResponseStatus_Unknown ResponseStatus = iota

	// The request succeeded
	ResponseStatus_Success

	// The request failed because of an internal error within the daemon
	ResponseStatus_Error

	// The request failed because the clients weren't synced yet, but synced clients were required for the request
	ResponseStatus_ClientsNotSynced

	// The request failed because the chain's state won't allow it to proceed. This is usually used for methods that
	// build transactions, but the preconditions for it aren't correct (and executing it will revert)
	ResponseStatus_InvalidChainState
)
