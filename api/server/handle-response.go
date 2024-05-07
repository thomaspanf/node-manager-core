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
func HandleInvalidMethod(logger *slog.Logger, w http.ResponseWriter) error {
	return writeResponse(w, logger, http.StatusMethodNotAllowed, "", nil, []byte{})
}

// Handles an error related to parsing the input parameters of a request
func HandleInputError(logger *slog.Logger, w http.ResponseWriter, err error) error {
	msg := err.Error()
	return writeResponse(w, logger, http.StatusBadRequest, "", err, formatError(msg))
}

// The request couldn't complete because the node requires an address but one wasn't present
func HandleAddressNotPresent(logger *slog.Logger, w http.ResponseWriter, err error) error {
	msg := fmt.Sprintf(addressNotPresentMessage, err.Error())
	return writeResponse(w, logger, http.StatusUnprocessableEntity, "Address not present", err, formatError(msg))
}

// The request couldn't complete because the node requires a wallet but one isn't present or useable
func HandleWalletNotReady(logger *slog.Logger, w http.ResponseWriter, err error) error {
	msg := fmt.Sprintf(walletNotReadyMessage, err.Error())
	return writeResponse(w, logger, http.StatusUnprocessableEntity, "Wallet not ready", err, formatError(msg))
}

// The request couldn't complete because it's trying to create a resource that already exists, or use a resource that conflicts with what's requested
func HandleResourceConflict(logger *slog.Logger, w http.ResponseWriter, err error) error {
	msg := fmt.Sprintf(resourceConflictMessage, err.Error())
	return writeResponse(w, logger, http.StatusConflict, "Resource conflict", err, formatError(msg))
}

// The request couldn't complete because it's trying to access a resource that didn't exist or couldn't be found
func HandleResourceNotFound(logger *slog.Logger, w http.ResponseWriter, err error) error {
	msg := fmt.Sprintf(resourceNotFoundMessage, err.Error())
	return writeResponse(w, logger, http.StatusNotFound, "Resource not found", err, formatError(msg))
}

// The request couldn't complete because the clients aren't synced yet
func HandleClientNotSynced(logger *slog.Logger, w http.ResponseWriter, err error) error {
	msg := clientsNotSyncedMessage
	return writeResponse(w, logger, http.StatusUnprocessableEntity, "Clients not synced", err, formatError(msg))
}

// The request couldn't complete because the chain state is preventing the request (it will revert if submitted)
func HandleInvalidChainState(logger *slog.Logger, w http.ResponseWriter, err error) error {
	msg := fmt.Sprintf(invalidChainStateMessage, err.Error())
	return writeResponse(w, logger, http.StatusUnprocessableEntity, "Invalid chain state", err, formatError(msg))
}

// The request couldn't complete because of a server error
func HandleServerError(logger *slog.Logger, w http.ResponseWriter, err error) error {
	msg := err.Error()
	return writeResponse(w, logger, http.StatusInternalServerError, "", err, formatError(msg))
}

// The request completed successfully
func HandleSuccess(logger *slog.Logger, w http.ResponseWriter, response any) error {
	// Serialize the response
	bytes, err := json.Marshal(response)
	if err != nil {
		return HandleServerError(logger, w, fmt.Errorf("error serializing response: %w", err))
	}

	// Write it
	logger.Debug("Response body", slog.String(log.BodyKey, string(bytes)))
	return writeResponse(w, logger, http.StatusOK, "", nil, bytes)
}

// Handles an API response for a request that could not be completed
func HandleFailedResponse(logger *slog.Logger, w http.ResponseWriter, status types.ResponseStatus, err error) error {
	switch status {
	case types.ResponseStatus_InvalidArguments:
		return HandleInputError(logger, w, err)
	case types.ResponseStatus_AddressNotPresent:
		return HandleAddressNotPresent(logger, w, err)
	case types.ResponseStatus_WalletNotReady:
		return HandleWalletNotReady(logger, w, err)
	case types.ResponseStatus_ResourceConflict:
		return HandleResourceConflict(logger, w, err)
	case types.ResponseStatus_ResourceNotFound:
		return HandleResourceNotFound(logger, w, err)
	case types.ResponseStatus_ClientsNotSynced:
		return HandleClientNotSynced(logger, w, err)
	case types.ResponseStatus_InvalidChainState:
		return HandleInvalidChainState(logger, w, err)
	case types.ResponseStatus_Error:
		return HandleServerError(logger, w, err)
	default:
		return HandleServerError(logger, w, fmt.Errorf("unknown response status: %d", status))
	}
}

// Handles an API response
func HandleResponse(logger *slog.Logger, w http.ResponseWriter, status types.ResponseStatus, response any, err error) error {
	switch status {
	case types.ResponseStatus_Success:
		return HandleSuccess(logger, w, response)
	default:
		return HandleFailedResponse(logger, w, status, err)
	}
}

// Writes a response to an HTTP request back to the client and logs it
func writeResponse(w http.ResponseWriter, logger *slog.Logger, statusCode int, cause string, err error, message []byte) error {
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
	logMsg := "Responded with:"
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
	_, writeErr := w.Write(message)
	return writeErr
}

// JSONifies an error for responding to requests
func formatError(message string) []byte {
	msg := types.ApiResponse[any]{
		Error: message,
	}

	bytes, _ := json.Marshal(msg)
	return bytes
}
