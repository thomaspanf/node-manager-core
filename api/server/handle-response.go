package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/rocket-pool/node-manager-core/api/types"
	"github.com/rocket-pool/node-manager-core/log"
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
func HandleInvalidMethod(logger *log.Logger, w http.ResponseWriter) {
	writeResponse(w, logger, http.StatusMethodNotAllowed, "", nil, []byte{})
}

// Handles an error related to parsing the input parameters of a request
func HandleInputError(logger *log.Logger, w http.ResponseWriter, err error) {
	msg := err.Error()
	writeResponse(w, logger, http.StatusBadRequest, "", err, formatError(msg))
}

// The request couldn't complete because the node requires an address but one wasn't present
func HandleAddressNotPresent(logger *log.Logger, w http.ResponseWriter, err error) {
	msg := fmt.Sprintf(addressNotPresentMessage, err.Error())
	writeResponse(w, logger, http.StatusUnprocessableEntity, "Address not present", err, formatError(msg))
}

// The request couldn't complete because the node requires a wallet but one isn't present or useable
func HandleWalletNotReady(logger *log.Logger, w http.ResponseWriter, err error) {
	msg := fmt.Sprintf(walletNotReadyMessage, err.Error())
	writeResponse(w, logger, http.StatusUnprocessableEntity, "Wallet not ready", err, formatError(msg))
}

// The request couldn't complete because it's trying to create a resource that already exists, or use a resource that conflicts with what's requested
func HandleResourceConflict(logger *log.Logger, w http.ResponseWriter, err error) {
	msg := fmt.Sprintf(resourceConflictMessage, err.Error())
	writeResponse(w, logger, http.StatusConflict, "Resource conflict", err, formatError(msg))
}

// The request couldn't complete because it's trying to access a resource that didn't exist or couldn't be found
func HandleResourceNotFound(logger *log.Logger, w http.ResponseWriter, err error) {
	msg := fmt.Sprintf(resourceNotFoundMessage, err.Error())
	writeResponse(w, logger, http.StatusNotFound, "Resource not found", err, formatError(msg))
}

// The request couldn't complete because the clients aren't synced yet
func HandleClientNotSynced(logger *log.Logger, w http.ResponseWriter, err error) {
	msg := clientsNotSyncedMessage
	writeResponse(w, logger, http.StatusUnprocessableEntity, "Clients not synced", err, formatError(msg))
}

// The request couldn't complete because the chain state is preventing the request (it will revert if submitted)
func HandleInvalidChainState(logger *log.Logger, w http.ResponseWriter, err error) {
	msg := fmt.Sprintf(invalidChainStateMessage, err.Error())
	writeResponse(w, logger, http.StatusUnprocessableEntity, "Invalid chain state", err, formatError(msg))
}

// The request couldn't complete because of a server error
func HandleServerError(logger *log.Logger, w http.ResponseWriter, err error) {
	msg := err.Error()
	writeResponse(w, logger, http.StatusInternalServerError, "", err, formatError(msg))
}

// The request completed successfully
func HandleSuccess(logger *log.Logger, w http.ResponseWriter, response any) {
	// Serialize the response
	bytes, err := json.Marshal(response)
	if err != nil {
		HandleServerError(logger, w, fmt.Errorf("error serializing response: %w", err))
		return
	}

	// Write it
	logger.Debug(string(bytes))
	writeResponse(w, logger, http.StatusOK, "", nil, bytes)
}

// Handles an API response for a request that could not be completed
func HandleFailedResponse(log *log.Logger, w http.ResponseWriter, status types.ResponseStatus, err error) {
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
func HandleResponse(log *log.Logger, w http.ResponseWriter, status types.ResponseStatus, response any, err error) {
	switch status {
	case types.ResponseStatus_Success:
		HandleSuccess(log, w, response)
	default:
		HandleFailedResponse(log, w, status, err)
	}
}

// Writes a response to an HTTP request back to the client and logs it
func writeResponse(w http.ResponseWriter, logger *log.Logger, statusCode int, cause string, err error, message []byte) {
	// Prep the log attributes
	codeMsg := fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode))
	attrs := []any{
		slog.String(log.CodeKey, codeMsg),
	}
	if cause != "" {
		attrs = append(attrs, slog.String(log.CauseKey, cause))
	}
	if err != nil {
		attrs = append(attrs, log.Err(err))
	}

	// Log the response
	logMsg := "Response"
	switch statusCode {
	case http.StatusOK:
		logger.Info(logMsg, attrs...)
	case http.StatusInternalServerError:
		logger.Error(logMsg, attrs...)
	default:
		logger.Warn(logMsg, attrs...)
	}

	// Write it to the client
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
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
