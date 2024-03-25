package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/rocket-pool/node-manager-core/api/types"
	"github.com/rocket-pool/node-manager-core/utils/log"
)

const (
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
	writeResponse(w, formatError(errorMsg, nil))
	log.Printlnf("[%d BAD_REQUEST] <= %s", http.StatusBadRequest, errorMsg)
}

// The request couldn't complete because the clients aren't synced yet
func HandleClientNotSynced(log *log.ColorLogger, w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	writeResponse(w, formatError(clientsNotSyncedMessage, nil))
	log.Printlnf("[%d UNPROCESSABLE ENTITY (Clients not synced)]", http.StatusUnprocessableEntity)
}

// The request couldn't complete because the chain state is preventing the request (it will revert if submitted)
func HandleInvalidChainState(log *log.ColorLogger, w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	errorMsg := err.Error()
	writeResponse(w, formatError(fmt.Sprintf(invalidChainStateMessage, errorMsg), nil))
	log.Printlnf("[%d UNPROCESSABLE ENTITY (Invalid chain state: %s)]", http.StatusUnprocessableEntity, errorMsg)
}

// The request couldn't complete because of a server error
func HandleServerError(log *log.ColorLogger, w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	stack := debug.Stack()
	errorMsg := err.Error()
	writeResponse(w, formatError(errorMsg, stack))
	log.Printlnf("[%d INTERNAL SERVER ERROR (%s)]\n%s", http.StatusInternalServerError, errorMsg, stack)
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
	case types.ResponseStatus_ClientsNotSynced:
		HandleClientNotSynced(log, w)
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
func formatError(message string, stackTrace []byte) []byte {
	type errorMessage struct {
		Error      string `json:"error"`
		StackTrace string `json:"stackTrace,omitempty"`
	}

	msg := errorMessage{
		Error: message,
	}
	if stackTrace != nil {
		msg.StackTrace = string(stackTrace)
	}

	bytes, _ := json.Marshal(msg)
	return bytes
}
