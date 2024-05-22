package client

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/goccy/go-json"
	"github.com/rocket-pool/node-manager-core/utils"
)

// Request types
type VoluntaryExitMessage struct {
	Epoch          Uinteger `json:"epoch"`
	ValidatorIndex string   `json:"validator_index"`
}
type VoluntaryExitRequest struct {
	Message   VoluntaryExitMessage `json:"message"`
	Signature ByteArray            `json:"signature"`
}
type BLSToExecutionChangeMessage struct {
	ValidatorIndex     string    `json:"validator_index"`
	FromBLSPubkey      ByteArray `json:"from_bls_pubkey"`
	ToExecutionAddress ByteArray `json:"to_execution_address"`
}
type BLSToExecutionChangeRequest struct {
	Message   BLSToExecutionChangeMessage `json:"message"`
	Signature ByteArray                   `json:"signature"`
}

// Response types
type SyncStatusResponse struct {
	Data struct {
		IsSyncing    bool     `json:"is_syncing"`
		HeadSlot     Uinteger `json:"head_slot"`
		SyncDistance Uinteger `json:"sync_distance"`
	} `json:"data"`
}
type Eth2ConfigResponse struct {
	Data struct {
		SecondsPerSlot               Uinteger  `json:"SECONDS_PER_SLOT"`
		SlotsPerEpoch                Uinteger  `json:"SLOTS_PER_EPOCH"`
		EpochsPerSyncCommitteePeriod Uinteger  `json:"EPOCHS_PER_SYNC_COMMITTEE_PERIOD"`
		CapellaForkVersion           ByteArray `json:"CAPELLA_FORK_VERSION"`
	} `json:"data"`
}
type Eth2DepositContractResponse struct {
	Data struct {
		ChainID Uinteger       `json:"chain_id"`
		Address common.Address `json:"address"`
	} `json:"data"`
}
type GenesisResponse struct {
	Data struct {
		GenesisTime           Uinteger  `json:"genesis_time"`
		GenesisForkVersion    ByteArray `json:"genesis_fork_version"`
		GenesisValidatorsRoot ByteArray `json:"genesis_validators_root"`
	} `json:"data"`
}
type FinalityCheckpointsResponse struct {
	Data struct {
		PreviousJustified struct {
			Epoch Uinteger `json:"epoch"`
		} `json:"previous_justified"`
		CurrentJustified struct {
			Epoch Uinteger `json:"epoch"`
		} `json:"current_justified"`
		Finalized struct {
			Epoch Uinteger `json:"epoch"`
		} `json:"finalized"`
	} `json:"data"`
}
type ForkResponse struct {
	Data struct {
		PreviousVersion ByteArray `json:"previous_version"`
		CurrentVersion  ByteArray `json:"current_version"`
		Epoch           Uinteger  `json:"epoch"`
	} `json:"data"`
}
type AttestationsResponse struct {
	Data []Attestation `json:"data"`
}
type BeaconBlockResponse struct {
	Data struct {
		Message struct {
			Slot          Uinteger `json:"slot"`
			ProposerIndex string   `json:"proposer_index"`
			Body          struct {
				Eth1Data struct {
					DepositRoot  ByteArray `json:"deposit_root"`
					DepositCount Uinteger  `json:"deposit_count"`
					BlockHash    ByteArray `json:"block_hash"`
				} `json:"eth1_data"`
				Attestations     []Attestation `json:"attestations"`
				ExecutionPayload *struct {
					FeeRecipient ByteArray `json:"fee_recipient"`
					BlockNumber  Uinteger  `json:"block_number"`
				} `json:"execution_payload"`
			} `json:"body"`
		} `json:"message"`
	} `json:"data"`
}
type BeaconBlockHeaderResponse struct {
	Finalized bool `json:"finalized"`
	Data      struct {
		Root      string `json:"root"`
		Canonical bool   `json:"canonical"`
		Header    struct {
			Message struct {
				Slot          Uinteger `json:"slot"`
				ProposerIndex string   `json:"proposer_index"`
			} `json:"message"`
		} `json:"header"`
	} `json:"data"`
}
type ValidatorsResponse struct {
	Data []Validator `json:"data"`
}
type Validator struct {
	Index     string   `json:"index"`
	Balance   Uinteger `json:"balance"`
	Status    string   `json:"status"`
	Validator struct {
		Pubkey                     ByteArray `json:"pubkey"`
		WithdrawalCredentials      ByteArray `json:"withdrawal_credentials"`
		EffectiveBalance           Uinteger  `json:"effective_balance"`
		Slashed                    bool      `json:"slashed"`
		ActivationEligibilityEpoch Uinteger  `json:"activation_eligibility_epoch"`
		ActivationEpoch            Uinteger  `json:"activation_epoch"`
		ExitEpoch                  Uinteger  `json:"exit_epoch"`
		WithdrawableEpoch          Uinteger  `json:"withdrawable_epoch"`
	} `json:"validator"`
}
type SyncDutiesResponse struct {
	Data []SyncDuty `json:"data"`
}
type SyncDuty struct {
	Pubkey               ByteArray  `json:"pubkey"`
	ValidatorIndex       string     `json:"validator_index"`
	SyncCommitteeIndices []Uinteger `json:"validator_sync_committee_indices"`
}
type ProposerDutiesResponse struct {
	Data []ProposerDuty `json:"data"`
}
type ProposerDuty struct {
	ValidatorIndex string `json:"validator_index"`
}

type CommitteesResponse struct {
	Data []Committee `json:"data"`
}

type Attestation struct {
	AggregationBits string `json:"aggregation_bits"`
	Data            struct {
		Slot  Uinteger `json:"slot"`
		Index Uinteger `json:"index"`
	} `json:"data"`
}

// Unsigned integer type
type Uinteger uint64

func (i Uinteger) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.Itoa(int(i)))
}
func (i *Uinteger) UnmarshalJSON(data []byte) error {

	// Unmarshal string
	var dataStr string
	if err := json.Unmarshal(data, &dataStr); err != nil {
		return err
	}

	// Parse integer value
	value, err := strconv.ParseUint(dataStr, 10, 64)
	if err != nil {
		return err
	}

	// Set value and return
	*i = Uinteger(value)
	return nil

}

// Byte array type
type ByteArray []byte

func (b ByteArray) MarshalJSON() ([]byte, error) {
	return json.Marshal(utils.EncodeHexWithPrefix(b))
}
func (b *ByteArray) UnmarshalJSON(data []byte) error {

	// Unmarshal string
	var dataStr string
	if err := json.Unmarshal(data, &dataStr); err != nil {
		return err
	}

	// Decode hex
	value, err := utils.DecodeHex(dataStr)
	if err != nil {
		return err
	}

	// Set value and return
	*b = value
	return nil

}
