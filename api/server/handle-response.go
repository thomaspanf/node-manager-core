package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rocket-pool/node-manager-core/api/types"
	"github.com/rocket-pool/node-manager-core/utils/log"
)

const (
	addressNotPresentMessage string = "The node requires an address for this request but one isn't present: %s"
	walletNotReadyMessage    string = "A wallet is required for this request but the node wallet isn't ready: %s"
	resourceConflictMessage  string = "Encountered a resource conflict: %s"
	resourceNotFoundMessage  string = "The requested resource could not be found: %s"
	clientsNotSyncedMessage  string = "The Execution Client and/or Beacon Node aren't finished syncing yet. Please try again once they've finished."
	invalidChainStateMessage string = "The Ethereum chain's state is not correct for the request: %s"
)

// Handle routes called with an invalid method
func HandleInvalidMethod(log *log.ColorLogger, w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	writeResponse(w, []byte{})
	log.Printlnf("[%d METHOD_NOT_ALLOWED]", http.StatusMethodNotAllowed)
}

// Handles an error related to parsing the input parameters of a request
func HandleInputError(log *log.ColorLogger, w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	errorMsg := err.Error()
	writeResponse(w, formatError(errorMsg))
	log.Printlnf("[%d BAD_REQUEST] <= %s", http.StatusBadRequest, errorMsg)
}

// The request couldn't complete because the node requires an address but one wasn't present
func HandleAddressNotPresent(log *log.ColorLogger, w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	errorMsg := err.Error()
	writeResponse(w, formatError(fmt.Sprintf(addressNotPresentMessage, errorMsg)))
	log.Printlnf("[%d UNPROCESSABLE ENTITY (Address not present: %s)]", http.StatusUnprocessableEntity, errorMsg)
}

// The request couldn't complete because the node requires a wallet but one isn't present or useable
func HandleWalletNotReady(log *log.ColorLogger, w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	errorMsg := err.Error()
	writeResponse(w, formatError(fmt.Sprintf(walletNotReadyMessage, errorMsg)))
	log.Printlnf("[%d UNPROCESSABLE ENTITY (Wallet not ready: %s)]", http.StatusUnprocessableEntity, errorMsg)
}

// The request couldn't complete because it's trying to create a resource that already exists, or use a resource that conflicts with what's requested
func HandleResourceConflict(log *log.ColorLogger, w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusConflict)
	errorMsg := err.Error()
	writeResponse(w, formatError(fmt.Sprintf(resourceConflictMessage, errorMsg)))
	log.Printlnf("[%d CONFLICT (Resource conflict: %s)]", http.StatusConflict, errorMsg)
}

// The request couldn't complete because it's trying to access a resource that didn't exist or couldn't be found
func HandleResourceNotFound(log *log.ColorLogger, w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusNotFound)
	errorMsg := err.Error()
	writeResponse(w, formatError(fmt.Sprintf(resourceNotFoundMessage, errorMsg)))
	log.Printlnf("[%d NOT FOUND (Resource not found: %s)]", http.StatusNotFound, errorMsg)
}

// The request couldn't complete because the clients aren't synced yet
func HandleClientNotSynced(log *log.ColorLogger, w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	errorMsg := err.Error()
	writeResponse(w, formatError(fmt.Sprintf(clientsNotSyncedMessage, errorMsg)))
	log.Printlnf("[%d UNPROCESSABLE ENTITY (Clients not synced: %s)]", http.StatusUnprocessableEntity, errorMsg)
}

// The request couldn't complete because the chain state is preventing the request (it will revert if submitted)
func HandleInvalidChainState(log *log.ColorLogger, w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	errorMsg := err.Error()
	writeResponse(w, formatError(fmt.Sprintf(invalidChainStateMessage, errorMsg)))
	log.Printlnf("[%d UNPROCESSABLE ENTITY (Invalid chain state: %s)]", http.StatusUnprocessableEntity, errorMsg)
}

// The request couldn't complete because of a server error
func HandleServerError(log *log.ColorLogger, w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	errorMsg := err.Error()
	writeResponse(w, formatError(errorMsg))
	log.Printlnf("[%d INTERNAL SERVER ERROR (%s)]", http.StatusInternalServerError, errorMsg)
}

// The request completed successfully
func HandleSuccess(log *log.ColorLogger, w http.ResponseWriter, response any, debug bool) {
	// Serialize the response
	bytes, err := json.Marshal(response)
	if err != nil {
		HandleServerError(log, w, fmt.Errorf("error serializing response: %w", err))
		return
	}

	// Write it
	writeResponse(w, bytes)
	if debug {
		log.Printlnf("[%d OK] <= %s", http.StatusOK, string(bytes))
	} else {
		log.Printlnf("[%d OK]", http.StatusOK)
	}
}

// Handles an API response for a request that could not be completed
func HandleFailedResponse(log *log.ColorLogger, w http.ResponseWriter, status types.ResponseStatus, err error) {
	switch status {
	case types.ResponseStatus_InvalidArguments:
		HandleInputError(log, w, err)
	case types.ResponseStatus_AddressNotPresent:
		HandleAddressNotPresent(log, w, err)
	case types.ResponseStatus_WalletNotReady:
		HandleWalletNotReady(log, w, err)
	case types.ResponseStatus_ResourceConflict:
		HandleResourceConflict(log, w, err)
	case types.ResponseStatus_ResourceNotFound:
		HandleResourceNotFound(log, w, err)
	case types.ResponseStatus_ClientsNotSynced:
		HandleClientNotSynced(log, w, err)
	case types.ResponseStatus_InvalidChainState:
		HandleInvalidChainState(log, w, err)
	case types.ResponseStatus_Error:
		HandleServerError(log, w, err)
	default:
		HandleServerError(log, w, fmt.Errorf("unknown response status: %d", status))
	}
}

// Handles an API response
func HandleResponse(log *log.ColorLogger, w http.ResponseWriter, status types.ResponseStatus, response any, err error, debug bool) {
	switch status {
	case types.ResponseStatus_Success:
		HandleSuccess(log, w, response, debug)
	default:
		HandleFailedResponse(log, w, status, err)
	}
}

// Writes a response to an HTTP request back to the client
func writeResponse(w http.ResponseWriter, message []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.Write(message)
}

// JSONifies an error for responding to requests
func formatError(message string) []byte {
	msg := types.ApiResponse[any]{
		Error: message,
	}

	bytes, _ := json.Marshal(msg)
	return bytes
}
