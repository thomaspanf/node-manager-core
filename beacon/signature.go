package beacon

import (
	"encoding/hex"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/rocket-pool/node-manager-core/utils"
	"gopkg.in/yaml.v3"
)

const (
	ValidatorSignatureLength int = 96
)

// A signature produced by a validator's private key.
type ValidatorSignature [ValidatorSignatureLength]byte

// Gets the string representation of the signature without a 0x prefix.
func (v ValidatorSignature) Hex() string {
	return hex.EncodeToString(v[:])
}

// Gets the string representation of the signature with a 0x prefix.
func (v ValidatorSignature) HexWithPrefix() string {
	return utils.EncodeHexWithPrefix(v[:])
}

// Gets the string representation of the signature without a 0x prefix.
func (v ValidatorSignature) String() string {
	return v.Hex()
}

// Converts a hex-encoded validator signature (with an optional 0x prefix) to a validator signature.
func HexToValidatorSignature(value string) (ValidatorSignature, error) {
	// Decode the value
	bytes, err := utils.DecodeHex(value)
	if err != nil {
		return ValidatorSignature{}, fmt.Errorf("error decoding validator signature: %w", err)
	}

	// Sanity check the length
	if len(bytes) != ValidatorSignatureLength {
		return ValidatorSignature{}, fmt.Errorf("invalid validator signature hex string %s: invalid length %d", value, len(value))
	}
	return ValidatorSignature(bytes), nil
}

// Serializes the signature to JSON.
func (v ValidatorSignature) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Hex())
}

// Deserializes the signature from JSON.
func (v *ValidatorSignature) UnmarshalJSON(data []byte) error {
	// Unmarshal the JSON
	var dataStr string
	if err := json.Unmarshal(data, &dataStr); err != nil {
		return fmt.Errorf("error decoding validator signature: %w", err)
	}

	// Decode the string
	signature, err := HexToValidatorSignature(dataStr)
	if err != nil {
		return fmt.Errorf("value '%s' cannot be decoded into a validator signature: %w", dataStr, err)
	}

	*v = signature
	return nil
}

// Serializes the signature to YAML.
func (v ValidatorSignature) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(v.Hex())
}

// Deserializes the signature from YAML.
func (v *ValidatorSignature) UnmarshalYAML(data []byte) error {
	// Unmarshal the YAML
	var dataStr string
	if err := yaml.Unmarshal(data, &dataStr); err != nil {
		return fmt.Errorf("error decoding validator signature: %w", err)
	}

	// Decode the string
	signature, err := HexToValidatorSignature(dataStr)
	if err != nil {
		return fmt.Errorf("value '%s' cannot be decoded into a validator signature: %w", dataStr, err)
	}

	*v = signature
	return nil
}
