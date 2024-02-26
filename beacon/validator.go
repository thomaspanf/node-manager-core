package beacon

import (
	"encoding/hex"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

// Encrypted validator keystore following the EIP-2335 standard
// (https://eips.ethereum.org/EIPS/eip-2335)
type ValidatorKeystore struct {
	Crypto  map[string]interface{} `json:"crypto"`
	Name    string                 `json:"name,omitempty"` // Technically not part of the spec but Prysm needs it
	Version uint                   `json:"version"`
	UUID    uuid.UUID              `json:"uuid"`
	Path    string                 `json:"path"`
	Pubkey  ValidatorPubkey        `json:"pubkey,omitempty"`
}

// Extended deposit data beyond what is required in an actual deposit message to Beacon, emulating what the deposit CLI produces
type ExtendedDepositData struct {
	PublicKey             ByteArray `json:"pubkey"`
	WithdrawalCredentials ByteArray `json:"withdrawal_credentials"`
	Amount                uint64    `json:"amount"`
	Signature             ByteArray `json:"signature"`
	DepositMessageRoot    ByteArray `json:"deposit_message_root"`
	DepositDataRoot       ByteArray `json:"deposit_data_root"`
	ForkVersion           ByteArray `json:"fork_version"`
	NetworkName           string    `json:"network_name"`
}

// Byte array type
type ByteArray []byte

func (b ByteArray) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(b))
}
func (b *ByteArray) UnmarshalJSON(data []byte) error {

	// Unmarshal string
	var dataStr string
	if err := json.Unmarshal(data, &dataStr); err != nil {
		return err
	}

	// Decode hex
	value, err := hex.DecodeString(dataStr)
	if err != nil {
		return err
	}

	// Set value and return
	*b = value
	return nil
}
