package services

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/node-manager-core/api/types"
	"github.com/rocket-pool/node-manager-core/beacon"
	"github.com/rocket-pool/node-manager-core/beacon/client"
)

// This is a proxy for multiple Beacon clients, providing natural fallback support if one of them fails.
type BeaconClientManager struct {
	primaryBc       beacon.IBeaconClient
	fallbackBc      beacon.IBeaconClient
	primaryReady    bool
	fallbackReady   bool
	ignoreSyncCheck bool
}

// Creates a new BeaconClientManager instance
func NewBeaconClientManager(primaryProvider string, fallbackProvider string, clientTimeout time.Duration) (*BeaconClientManager, error) {
	var primaryBc beacon.IBeaconClient
	var fallbackBc beacon.IBeaconClient
	primaryBc = client.NewStandardHttpClient(primaryProvider, clientTimeout)
	if fallbackProvider != "" {
		fallbackBc = client.NewStandardHttpClient(fallbackProvider, clientTimeout)
	}

	return &BeaconClientManager{
		primaryBc:     primaryBc,
		fallbackBc:    fallbackBc,
		primaryReady:  true,
		fallbackReady: fallbackBc != nil,
	}, nil
}

/// ========================
/// IClientManager Functions
/// ========================

func (m *BeaconClientManager) GetPrimaryClient() beacon.IBeaconClient {
	return m.primaryBc
}

func (m *BeaconClientManager) GetFallbackClient() beacon.IBeaconClient {
	return m.fallbackBc
}

func (m *BeaconClientManager) IsPrimaryReady() bool {
	return m.primaryReady
}

func (m *BeaconClientManager) IsFallbackReady() bool {
	return m.fallbackReady
}

func (m *BeaconClientManager) IsFallbackEnabled() bool {
	return m.fallbackBc != nil
}

func (m *BeaconClientManager) GetClientTypeName() string {
	return "Beacon Node"
}

func (m *BeaconClientManager) setPrimaryReady(ready bool) {
	m.primaryReady = ready
}

func (m *BeaconClientManager) setFallbackReady(ready bool) {
	m.fallbackReady = ready
}

/// =======================
/// IBeaconClient Functions
/// =======================

// Get the client's sync status
func (m *BeaconClientManager) GetSyncStatus(ctx context.Context) (beacon.SyncStatus, error) {
	return runFunction1(m, ctx, func(client beacon.IBeaconClient) (beacon.SyncStatus, error) {
		return client.GetSyncStatus(ctx)
	})
}

// Get the Beacon configuration
func (m *BeaconClientManager) GetEth2Config(ctx context.Context) (beacon.Eth2Config, error) {
	return runFunction1(m, ctx, func(client beacon.IBeaconClient) (beacon.Eth2Config, error) {
		return client.GetEth2Config(ctx)
	})
}

// Get the Beacon configuration
func (m *BeaconClientManager) GetEth2DepositContract(ctx context.Context) (beacon.Eth2DepositContract, error) {
	return runFunction1(m, ctx, func(client beacon.IBeaconClient) (beacon.Eth2DepositContract, error) {
		return client.GetEth2DepositContract(ctx)
	})
}

// Get the attestations in a Beacon chain block
func (m *BeaconClientManager) GetAttestations(ctx context.Context, blockId string) ([]beacon.AttestationInfo, bool, error) {
	return runFunction2(m, ctx, func(client beacon.IBeaconClient) ([]beacon.AttestationInfo, bool, error) {
		return client.GetAttestations(ctx, blockId)
	})
}

// Get a Beacon chain block
func (m *BeaconClientManager) GetBeaconBlock(ctx context.Context, blockId string) (beacon.BeaconBlock, bool, error) {
	return runFunction2(m, ctx, func(client beacon.IBeaconClient) (beacon.BeaconBlock, bool, error) {
		return client.GetBeaconBlock(ctx, blockId)
	})
}

// Get the header of a Beacon chain block
func (m *BeaconClientManager) GetBeaconBlockHeader(ctx context.Context, blockId string) (beacon.BeaconBlockHeader, bool, error) {
	return runFunction2(m, ctx, func(client beacon.IBeaconClient) (beacon.BeaconBlockHeader, bool, error) {
		return client.GetBeaconBlockHeader(ctx, blockId)
	})
}

// Get the Beacon chain's head information
func (m *BeaconClientManager) GetBeaconHead(ctx context.Context) (beacon.BeaconHead, error) {
	return runFunction1(m, ctx, func(client beacon.IBeaconClient) (beacon.BeaconHead, error) {
		return client.GetBeaconHead(ctx)
	})
}

// Get a validator's status by its index
func (m *BeaconClientManager) GetValidatorStatusByIndex(ctx context.Context, index string, opts *beacon.ValidatorStatusOptions) (beacon.ValidatorStatus, error) {
	return runFunction1(m, ctx, func(client beacon.IBeaconClient) (beacon.ValidatorStatus, error) {
		return client.GetValidatorStatusByIndex(ctx, index, opts)
	})
}

// Get a validator's status by its pubkey
func (m *BeaconClientManager) GetValidatorStatus(ctx context.Context, pubkey beacon.ValidatorPubkey, opts *beacon.ValidatorStatusOptions) (beacon.ValidatorStatus, error) {
	return runFunction1(m, ctx, func(client beacon.IBeaconClient) (beacon.ValidatorStatus, error) {
		return client.GetValidatorStatus(ctx, pubkey, opts)
	})
}

// Get the statuses of multiple validators by their pubkeys
func (m *BeaconClientManager) GetValidatorStatuses(ctx context.Context, pubkeys []beacon.ValidatorPubkey, opts *beacon.ValidatorStatusOptions) (map[beacon.ValidatorPubkey]beacon.ValidatorStatus, error) {
	return runFunction1(m, ctx, func(client beacon.IBeaconClient) (map[beacon.ValidatorPubkey]beacon.ValidatorStatus, error) {
		return client.GetValidatorStatuses(ctx, pubkeys, opts)
	})
}

// Get a validator's index
func (m *BeaconClientManager) GetValidatorIndex(ctx context.Context, pubkey beacon.ValidatorPubkey) (string, error) {
	return runFunction1(m, ctx, func(client beacon.IBeaconClient) (string, error) {
		return client.GetValidatorIndex(ctx, pubkey)
	})
}

// Get a validator's sync duties
func (m *BeaconClientManager) GetValidatorSyncDuties(ctx context.Context, indices []string, epoch uint64) (map[string]bool, error) {
	return runFunction1(m, ctx, func(client beacon.IBeaconClient) (map[string]bool, error) {
		return client.GetValidatorSyncDuties(ctx, indices, epoch)
	})
}

// Get a validator's proposer duties
func (m *BeaconClientManager) GetValidatorProposerDuties(ctx context.Context, indices []string, epoch uint64) (map[string]uint64, error) {
	return runFunction1(m, ctx, func(client beacon.IBeaconClient) (map[string]uint64, error) {
		return client.GetValidatorProposerDuties(ctx, indices, epoch)
	})
}

// Get the Beacon chain's domain data
func (m *BeaconClientManager) GetDomainData(ctx context.Context, domainType []byte, epoch uint64, useGenesisFork bool) ([]byte, error) {
	return runFunction1(m, ctx, func(client beacon.IBeaconClient) ([]byte, error) {
		return client.GetDomainData(ctx, domainType, epoch, useGenesisFork)
	})
}

// Voluntarily exit a validator
func (m *BeaconClientManager) ExitValidator(ctx context.Context, validatorIndex string, epoch uint64, signature beacon.ValidatorSignature) error {
	return runFunction0(m, ctx, func(client beacon.IBeaconClient) error {
		return client.ExitValidator(ctx, validatorIndex, epoch, signature)
	})
}

// Close the connection to the Beacon client
func (m *BeaconClientManager) Close(ctx context.Context) error {
	return runFunction0(m, ctx, func(client beacon.IBeaconClient) error {
		return client.Close(ctx)
	})
}

// Get the EL data for a CL block
func (m *BeaconClientManager) GetEth1DataForEth2Block(ctx context.Context, blockId string) (beacon.Eth1Data, bool, error) {
	return runFunction2(m, ctx, func(client beacon.IBeaconClient) (beacon.Eth1Data, bool, error) {
		return client.GetEth1DataForEth2Block(ctx, blockId)
	})
}

// Get the attestation committees for an epoch
func (m *BeaconClientManager) GetCommitteesForEpoch(ctx context.Context, epoch *uint64) (beacon.Committees, error) {
	return runFunction1(m, ctx, func(client beacon.IBeaconClient) (beacon.Committees, error) {
		return client.GetCommitteesForEpoch(ctx, epoch)
	})
}

// Change the withdrawal credentials for a validator
func (m *BeaconClientManager) ChangeWithdrawalCredentials(ctx context.Context, validatorIndex string, fromBlsPubkey beacon.ValidatorPubkey, toExecutionAddress common.Address, signature beacon.ValidatorSignature) error {
	return runFunction0(m, ctx, func(client beacon.IBeaconClient) error {
		return client.ChangeWithdrawalCredentials(ctx, validatorIndex, fromBlsPubkey, toExecutionAddress, signature)
	})
}

/// =================
/// Manager Functions
/// =================

// Get the status of the primary and fallback clients
func (m *BeaconClientManager) CheckStatus(ctx context.Context) *types.ClientManagerStatus {
	status := &types.ClientManagerStatus{
		FallbackEnabled: m.fallbackBc != nil,
	}

	// Ignore the sync check and just use the predefined settings if requested
	if m.ignoreSyncCheck {
		status.PrimaryClientStatus.IsWorking = m.primaryReady
		status.PrimaryClientStatus.IsSynced = m.primaryReady
		if status.FallbackEnabled {
			status.FallbackClientStatus.IsWorking = m.fallbackReady
			status.FallbackClientStatus.IsSynced = m.fallbackReady
		}
		return status
	}

	// Get the primary BC status
	status.PrimaryClientStatus = checkBcStatus(ctx, m.primaryBc)

	// Get the fallback BC status if applicable
	if status.FallbackEnabled {
		status.FallbackClientStatus = checkBcStatus(ctx, m.fallbackBc)
	}

	// Flag the ready clients
	m.primaryReady = (status.PrimaryClientStatus.IsWorking && status.PrimaryClientStatus.IsSynced)
	m.fallbackReady = (status.FallbackEnabled && status.FallbackClientStatus.IsWorking && status.FallbackClientStatus.IsSynced)

	return status
}

// Check the client status
func checkBcStatus(ctx context.Context, client beacon.IBeaconClient) types.ClientStatus {
	status := types.ClientStatus{}

	// Get the fallback's sync progress
	syncStatus, err := client.GetSyncStatus(ctx)
	if err != nil {
		status.Error = fmt.Sprintf("Sync progress check failed with [%s]", err.Error())
		status.IsSynced = false
		status.IsWorking = false
		return status
	}

	// Return the sync status
	if !syncStatus.Syncing {
		status.IsWorking = true
		status.IsSynced = true
		status.SyncProgress = 1
	} else {
		status.IsWorking = true
		status.IsSynced = false
		status.SyncProgress = syncStatus.Progress
	}
	return status
}
