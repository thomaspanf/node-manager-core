package beacon

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

// Beacon Node interface
type IBeaconClient interface {
	GetSyncStatus(ctx context.Context) (SyncStatus, error)
	GetEth2Config(ctx context.Context) (Eth2Config, error)
	GetEth2DepositContract(ctx context.Context) (Eth2DepositContract, error)
	GetAttestations(ctx context.Context, blockId string) ([]AttestationInfo, bool, error)
	GetBeaconBlock(ctx context.Context, blockId string) (BeaconBlock, bool, error)
	GetBeaconHead(ctx context.Context) (BeaconHead, error)
	GetValidatorStatusByIndex(ctx context.Context, index string, opts *ValidatorStatusOptions) (ValidatorStatus, error)
	GetValidatorStatus(ctx context.Context, pubkey ValidatorPubkey, opts *ValidatorStatusOptions) (ValidatorStatus, error)
	GetValidatorStatuses(ctx context.Context, pubkeys []ValidatorPubkey, opts *ValidatorStatusOptions) (map[ValidatorPubkey]ValidatorStatus, error)
	GetValidatorIndex(ctx context.Context, pubkey ValidatorPubkey) (string, error)
	GetValidatorSyncDuties(ctx context.Context, indices []string, epoch uint64) (map[string]bool, error)
	GetValidatorProposerDuties(ctx context.Context, indices []string, epoch uint64) (map[string]uint64, error)
	GetDomainData(ctx context.Context, domainType []byte, epoch uint64, useGenesisFork bool) ([]byte, error)
	ExitValidator(ctx context.Context, validatorIndex string, epoch uint64, signature ValidatorSignature) error
	Close(ctx context.Context) error
	GetEth1DataForEth2Block(ctx context.Context, blockId string) (Eth1Data, bool, error)
	GetCommitteesForEpoch(ctx context.Context, epoch *uint64) (Committees, error)
	ChangeWithdrawalCredentials(ctx context.Context, validatorIndex string, fromBlsPubkey ValidatorPubkey, toExecutionAddress common.Address, signature ValidatorSignature) error
}
