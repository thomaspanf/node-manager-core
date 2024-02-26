package beacon

import (
	"encoding/hex"
	"fmt"

	"github.com/goccy/go-json"
	nmc_utils "github.com/rocket-pool/node-manager-core/utils"
	"gopkg.in/yaml.v3"
)

const (
	ValidatorPubkeyLength int = 48
)

// A validator's pubkey key
type ValidatorPubkey [ValidatorPubkeyLength]byte

// Gets the string representation of the pubkey without a 0x prefix.
func (v ValidatorPubkey) Hex() string {
	return hex.EncodeToString(v[:])
}

// Gets the string representation of the pubkey with a 0x prefix.
func (v ValidatorPubkey) HexWithPrefix() string {
	return nmc_utils.EncodeHexWithPrefix(v[:])
}

// Gets the string representation of the pubkey without a 0x prefix.
func (v ValidatorPubkey) String() string {
	return v.Hex()
}

// Converts a hex-encoded validator pubkey (with an optional 0x prefix) to a validator pubkey.
func HexToValidatorPubkey(value string) (ValidatorPubkey, error) {
	// Decode the value
	bytes, err := nmc_utils.DecodeHex(value)
	if err != nil {
		return ValidatorPubkey{}, fmt.Errorf("error decoding validator pubkey: %w", err)
	}

	// Sanity check the length
	if len(bytes) != ValidatorPubkeyLength {
		return ValidatorPubkey{}, fmt.Errorf("invalid validator pubkey hex string %s: invalid length %d", value, len(value))
	}
	return ValidatorPubkey(bytes), nil
}

// Serializes the pubkey to JSON.
func (v ValidatorPubkey) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Hex())
}

// Deserializes the pubkey from JSON.
func (v *ValidatorPubkey) UnmarshalJSON(data []byte) error {
	// Unmarshal the JSON
	var dataStr string
	if err := json.Unmarshal(data, &dataStr); err != nil {
		return fmt.Errorf("error decoding validator pubkey: %w", err)
	}

	// Decode the string
	pubkey, err := HexToValidatorPubkey(dataStr)
	if err != nil {
		return fmt.Errorf("value '%s' cannot be decoded into a validator pubkey: %w", dataStr, err)
	}

	*v = pubkey
	return nil
}

// Serializes the pubkey to YAML.
func (v ValidatorPubkey) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(v.Hex())
}

// Deserializes the pubkey from YAML.
func (v *ValidatorPubkey) UnmarshalYAML(data []byte) error {
	// Unmarshal the YAML
	var dataStr string
	if err := yaml.Unmarshal(data, &dataStr); err != nil {
		return fmt.Errorf("error decoding validator pubkey: %w", err)
	}

	// Decode the string
	pubkey, err := HexToValidatorPubkey(dataStr)
	if err != nil {
		return fmt.Errorf("value '%s' cannot be decoded into a validator pubkey: %w", dataStr, err)
	}

	*v = pubkey
	return nil
}
