package nmc_utils

import (
	"encoding/hex"
	"strings"
)

const (
	hexPrefix string = "0x"
)

// Encode bytes to a hex string with a 0x prefix
func EncodeHexWithPrefix(value []byte) string {
	hexString := hex.EncodeToString(value)
	return hexPrefix + hexString
}

// Convert a hex-encoded string to a byte array, removing the 0x prefix if present
func DecodeHex(value string) ([]byte, error) {
	value = strings.TrimPrefix(value, hexPrefix)
	return hex.DecodeString(value)
}

// Add a 0x prefix to a hex string if not present
func AddPrefix(value string) string {
	if !strings.HasPrefix(value, hexPrefix) {
		return hexPrefix + value
	}
	return value
}

// Remove a 0x prefix from a hex string if present
func RemovePrefix(value string) string {
	return strings.TrimPrefix(value, hexPrefix)
}
