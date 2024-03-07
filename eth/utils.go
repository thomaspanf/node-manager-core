package eth

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"regexp"

	batch "github.com/rocket-pool/batch-query"
)

const (
	// Regex to check for reversion messages from Nethermind
	nethermindRevertRegexString string = "Reverted 0x(?P<message>[0-9a-fA-F]+).*"
)

var (
	// Regex to check for reversion messages from Nethermind
	nethermindRevertRegex *regexp.Regexp = regexp.MustCompile(nethermindRevertRegexString)
)

// Create a transaction submission directly from serialized info (and the error provided by the transaction info constructor),
// using the SafeGasLimit as the GasLimit for the submission automatically.
func CreateTxSubmissionFromInfo(txInfo *TransactionInfo, err error) (*TransactionSubmission, error) {
	if err != nil {
		return nil, err
	}
	return &TransactionSubmission{
		TxInfo:   txInfo,
		GasLimit: txInfo.SimulationResult.SafeGasLimit,
	}, nil
}

// Simple convenience method to add a contract call to a multicaller
func AddCallToMulicaller(mc *batch.MultiCaller, contract *Contract, output any, method string, args ...any) {
	mc.AddCall(contract.Address, contract.ABI, output, method, args...)
}

// Adds a collection of IQueryable calls to a multicall
func AddQueryablesToMulticall(mc *batch.MultiCaller, queryables ...IQueryable) {
	for _, queryable := range queryables {
		queryable.AddToQuery(mc)
	}
}

// Adds all of the object's fields that implement IQueryable to the provided multicaller
func QueryAllFields(object any, mc *batch.MultiCaller) error {
	objectValue := reflect.ValueOf(object)
	objectType := reflect.TypeOf(object)
	if objectType.Kind() == reflect.Pointer {
		// If this is a pointer, switch to what it's pointing at
		objectValue = objectValue.Elem()
		objectType = objectType.Elem()
	}

	// Run through each field
	for i := 0; i < objectType.NumField(); i++ {
		field := objectValue.Field(i)
		typeField := objectType.Field(i)
		if typeField.IsExported() {
			fieldAsQueryable, isQueryable := field.Interface().(IQueryable)
			if isQueryable {
				// If it's IQueryable, run it
				fieldAsQueryable.AddToQuery(mc)
			} else if typeField.Type.Kind() == reflect.Pointer &&
				typeField.Type.Elem().Kind() == reflect.Struct {
				// If it's a pointer to a struct, recurse
				err := QueryAllFields(field.Interface(), mc)
				if err != nil {
					return err
				}
			} else if typeField.Type.Kind() == reflect.Struct {
				// If it's a struct, recurse
				err := QueryAllFields(field.Interface(), mc)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Normalize revert messages so they're all in ASCII format
func normalizeRevertMessage(err error) error {
	if err == nil {
		return err
	}

	// Get the message in hex format, if it exists
	matches := nethermindRevertRegex.FindStringSubmatch(err.Error())
	if matches == nil {
		return err
	}
	messageIndex := nethermindRevertRegex.SubexpIndex("message")
	if messageIndex == -1 {
		return err
	}
	message := matches[messageIndex]

	// Convert the hex message to ASCII
	bytes, err2 := hex.DecodeString(message)
	if err2 != nil {
		return err // Return the original error if decoding failed somehow
	}

	return fmt.Errorf("reverted: %s", string(bytes))
}
