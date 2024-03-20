package services

import (
	"context"
	"errors"
	"fmt"
	"net"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fatih/color"
	"github.com/rocket-pool/node-manager-core/api/types"
	"github.com/rocket-pool/node-manager-core/beacon"
	"github.com/rocket-pool/node-manager-core/beacon/client"
	"github.com/rocket-pool/node-manager-core/utils/log"
)

// This is a proxy for multiple Beacon clients, providing natural fallback support if one of them fails.
type BeaconClientManager struct {
	primaryBc       beacon.IBeaconClient
	fallbackBc      beacon.IBeaconClient
	logger          log.ColorLogger
	primaryReady    bool
	fallbackReady   bool
	ignoreSyncCheck bool
}

// This is a signature for a wrapped Beacon client function that only returns an error
type bcFunction0 func(beacon.IBeaconClient) error

// This is a signature for a wrapped Beacon client function that returns 1 var and an error
type bcFunction1 func(beacon.IBeaconClient) (interface{}, error)

// This is a signature for a wrapped Beacon client function that returns 2 vars and an error
type bcFunction2 func(beacon.IBeaconClient) (interface{}, interface{}, error)

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
		logger:        log.NewColorLogger(color.FgHiBlue),
		primaryReady:  true,
		fallbackReady: fallbackBc != nil,
	}, nil
}

func (m *BeaconClientManager) IsPrimaryReady() bool {
	return m.primaryReady
}

func (m *BeaconClientManager) IsFallbackReady() bool {
	return m.fallbackReady
}

/// ======================
/// IBeaconClient Functions
/// ======================

// Get the client's sync status
func (m *BeaconClientManager) GetSyncStatus(ctx context.Context) (beacon.SyncStatus, error) {
	result, err := m.runFunction1(func(client beacon.IBeaconClient) (interface{}, error) {
		return client.GetSyncStatus(ctx)
	})
	if err != nil {
		return beacon.SyncStatus{}, err
	}
	return result.(beacon.SyncStatus), nil
}

// Get the Beacon configuration
func (m *BeaconClientManager) GetEth2Config(ctx context.Context) (beacon.Eth2Config, error) {
	result, err := m.runFunction1(func(client beacon.IBeaconClient) (interface{}, error) {
		return client.GetEth2Config(ctx)
	})
	if err != nil {
		return beacon.Eth2Config{}, err
	}
	return result.(beacon.Eth2Config), nil
}

// Get the Beacon configuration
func (m *BeaconClientManager) GetEth2DepositContract(ctx context.Context) (beacon.Eth2DepositContract, error) {
	result, err := m.runFunction1(func(client beacon.IBeaconClient) (interface{}, error) {
		return client.GetEth2DepositContract(ctx)
	})
	if err != nil {
		return beacon.Eth2DepositContract{}, err
	}
	return result.(beacon.Eth2DepositContract), nil
}

// Get the attestations in a Beacon chain block
func (m *BeaconClientManager) GetAttestations(ctx context.Context, blockId string) ([]beacon.AttestationInfo, bool, error) {
	result1, result2, err := m.runFunction2(func(client beacon.IBeaconClient) (interface{}, interface{}, error) {
		return client.GetAttestations(ctx, blockId)
	})
	if err != nil {
		return nil, false, err
	}
	return result1.([]beacon.AttestationInfo), result2.(bool), nil
}

// Get a Beacon chain block
func (m *BeaconClientManager) GetBeaconBlock(ctx context.Context, blockId string) (beacon.BeaconBlock, bool, error) {
	result1, result2, err := m.runFunction2(func(client beacon.IBeaconClient) (interface{}, interface{}, error) {
		return client.GetBeaconBlock(ctx, blockId)
	})
	if err != nil {
		return beacon.BeaconBlock{}, false, err
	}
	return result1.(beacon.BeaconBlock), result2.(bool), nil
}

// Get the header of a Beacon chain block
func (m *BeaconClientManager) GetBeaconBlockHeader(ctx context.Context, blockId string) (beacon.BeaconBlockHeader, bool, error) {
	result1, result2, err := m.runFunction2(func(client beacon.IBeaconClient) (interface{}, interface{}, error) {
		return client.GetBeaconBlockHeader(ctx, blockId)
	})
	if err != nil {
		return beacon.BeaconBlockHeader{}, false, err
	}
	return result1.(beacon.BeaconBlockHeader), result2.(bool), nil
}

// Get the Beacon chain's head information
func (m *BeaconClientManager) GetBeaconHead(ctx context.Context) (beacon.BeaconHead, error) {
	result, err := m.runFunction1(func(client beacon.IBeaconClient) (interface{}, error) {
		return client.GetBeaconHead(ctx)
	})
	if err != nil {
		return beacon.BeaconHead{}, err
	}
	return result.(beacon.BeaconHead), nil
}

// Get a validator's status by its index
func (m *BeaconClientManager) GetValidatorStatusByIndex(ctx context.Context, index string, opts *beacon.ValidatorStatusOptions) (beacon.ValidatorStatus, error) {
	result, err := m.runFunction1(func(client beacon.IBeaconClient) (interface{}, error) {
		return client.GetValidatorStatusByIndex(ctx, index, opts)
	})
	if err != nil {
		return beacon.ValidatorStatus{}, err
	}
	return result.(beacon.ValidatorStatus), nil
}

// Get a validator's status by its pubkey
func (m *BeaconClientManager) GetValidatorStatus(ctx context.Context, pubkey beacon.ValidatorPubkey, opts *beacon.ValidatorStatusOptions) (beacon.ValidatorStatus, error) {
	result, err := m.runFunction1(func(client beacon.IBeaconClient) (interface{}, error) {
		return client.GetValidatorStatus(ctx, pubkey, opts)
	})
	if err != nil {
		return beacon.ValidatorStatus{}, err
	}
	return result.(beacon.ValidatorStatus), nil
}

// Get the statuses of multiple validators by their pubkeys
func (m *BeaconClientManager) GetValidatorStatuses(ctx context.Context, pubkeys []beacon.ValidatorPubkey, opts *beacon.ValidatorStatusOptions) (map[beacon.ValidatorPubkey]beacon.ValidatorStatus, error) {
	result, err := m.runFunction1(func(client beacon.IBeaconClient) (interface{}, error) {
		return client.GetValidatorStatuses(ctx, pubkeys, opts)
	})
	if err != nil {
		return nil, err
	}
	return result.(map[beacon.ValidatorPubkey]beacon.ValidatorStatus), nil
}

// Get a validator's index
func (m *BeaconClientManager) GetValidatorIndex(ctx context.Context, pubkey beacon.ValidatorPubkey) (string, error) {
	result, err := m.runFunction1(func(client beacon.IBeaconClient) (interface{}, error) {
		return client.GetValidatorIndex(ctx, pubkey)
	})
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

// Get a validator's sync duties
func (m *BeaconClientManager) GetValidatorSyncDuties(ctx context.Context, indices []string, epoch uint64) (map[string]bool, error) {
	result, err := m.runFunction1(func(client beacon.IBeaconClient) (interface{}, error) {
		return client.GetValidatorSyncDuties(ctx, indices, epoch)
	})
	if err != nil {
		return nil, err
	}
	return result.(map[string]bool), nil
}

// Get a validator's proposer duties
func (m *BeaconClientManager) GetValidatorProposerDuties(ctx context.Context, indices []string, epoch uint64) (map[string]uint64, error) {
	result, err := m.runFunction1(func(client beacon.IBeaconClient) (interface{}, error) {
		return client.GetValidatorProposerDuties(ctx, indices, epoch)
	})
	if err != nil {
		return nil, err
	}
	return result.(map[string]uint64), nil
}

// Get the Beacon chain's domain data
func (m *BeaconClientManager) GetDomainData(ctx context.Context, domainType []byte, epoch uint64, useGenesisFork bool) ([]byte, error) {
	result, err := m.runFunction1(func(client beacon.IBeaconClient) (interface{}, error) {
		return client.GetDomainData(ctx, domainType, epoch, useGenesisFork)
	})
	if err != nil {
		return nil, err
	}
	return result.([]byte), nil
}

// Voluntarily exit a validator
func (m *BeaconClientManager) ExitValidator(ctx context.Context, validatorIndex string, epoch uint64, signature beacon.ValidatorSignature) error {
	err := m.runFunction0(func(client beacon.IBeaconClient) error {
		return client.ExitValidator(ctx, validatorIndex, epoch, signature)
	})
	return err
}

// Close the connection to the Beacon client
func (m *BeaconClientManager) Close(ctx context.Context) error {
	err := m.runFunction0(func(client beacon.IBeaconClient) error {
		return client.Close(ctx)
	})
	return err
}

// Get the EL data for a CL block
func (m *BeaconClientManager) GetEth1DataForEth2Block(ctx context.Context, blockId string) (beacon.Eth1Data, bool, error) {
	result1, result2, err := m.runFunction2(func(client beacon.IBeaconClient) (interface{}, interface{}, error) {
		return client.GetEth1DataForEth2Block(ctx, blockId)
	})
	if err != nil {
		return beacon.Eth1Data{}, false, err
	}
	return result1.(beacon.Eth1Data), result2.(bool), nil
}

// Get the attestation committees for an epoch
func (m *BeaconClientManager) GetCommitteesForEpoch(ctx context.Context, epoch *uint64) (beacon.Committees, error) {
	result, err := m.runFunction1(func(client beacon.IBeaconClient) (interface{}, error) {
		return client.GetCommitteesForEpoch(ctx, epoch)
	})
	if err != nil {
		return nil, err
	}
	return result.(beacon.Committees), nil
}

// Change the withdrawal credentials for a validator
func (m *BeaconClientManager) ChangeWithdrawalCredentials(ctx context.Context, validatorIndex string, fromBlsPubkey beacon.ValidatorPubkey, toExecutionAddress common.Address, signature beacon.ValidatorSignature) error {
	err := m.runFunction0(func(client beacon.IBeaconClient) error {
		return client.ChangeWithdrawalCredentials(ctx, validatorIndex, fromBlsPubkey, toExecutionAddress, signature)
	})
	if err != nil {
		return err
	}
	return nil
}

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

/// ==================
/// Internal Functions
/// ==================

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

// Attempts to run a function progressively through each client until one succeeds or they all fail.
func (m *BeaconClientManager) runFunction0(function bcFunction0) error {

	// Check if we can use the primary
	if m.primaryReady {
		// Try to run the function on the primary
		err := function(m.primaryBc)
		if err != nil {
			if m.isDisconnected(err) {
				// If it's disconnected, log it and try the fallback
				m.primaryReady = false
				if m.fallbackBc != nil {
					m.logger.Printlnf("WARNING: Primary Beacon client disconnected (%s), using fallback...", err.Error())
					return m.runFunction0(function)
				} else {
					m.logger.Printlnf("WARNING: Primary Beacon client disconnected (%s) and no fallback is configured.", err.Error())
					return fmt.Errorf("all Beacon clients failed")
				}
			}
			// If it's a different error, just return it
			return err
		}
		// If there's no error, return the result
		return nil
	}

	if m.fallbackReady {
		// Try to run the function on the fallback
		err := function(m.fallbackBc)
		if err != nil {
			if m.isDisconnected(err) {
				// If it's disconnected, log it and try the fallback
				m.logger.Printlnf("WARNING: Fallback Beacon client disconnected (%s)", err.Error())
				m.fallbackReady = false
				return fmt.Errorf("all Beacon clients failed")
			}

			// If it's a different error, just return it
			return err
		}
		// If there's no error, return the result
		return nil
	}

	return fmt.Errorf("no Beacon clients were ready")
}

// Attempts to run a function progressively through each client until one succeeds or they all fail.
func (m *BeaconClientManager) runFunction1(function bcFunction1) (interface{}, error) {
	// Check if we can use the primary
	if m.primaryReady {
		// Try to run the function on the primary
		result, err := function(m.primaryBc)
		if err != nil {
			if m.isDisconnected(err) {
				// If it's disconnected, log it and try the fallback
				m.primaryReady = false
				if m.fallbackBc != nil {
					m.logger.Printlnf("WARNING: Primary Beacon client disconnected (%s), using fallback...", err.Error())
					return m.runFunction1(function)
				} else {
					m.logger.Printlnf("WARNING: Primary Beacon client disconnected (%s) and no fallback is configured.", err.Error())
					return nil, fmt.Errorf("all Beacon clients failed")
				}
			}
			// If it's a different error, just return it
			return nil, err
		}
		// If there's no error, return the result
		return result, nil
	}

	if m.fallbackReady {
		// Try to run the function on the fallback
		result, err := function(m.fallbackBc)
		if err != nil {
			if m.isDisconnected(err) {
				// If it's disconnected, log it and try the fallback
				m.logger.Printlnf("WARNING: Fallback Beacon client disconnected (%s)", err.Error())
				m.fallbackReady = false
				return nil, fmt.Errorf("all Beacon clients failed")
			}
			// If it's a different error, just return it
			return nil, err
		}
		// If there's no error, return the result
		return result, nil
	}

	return nil, fmt.Errorf("no Beacon clients were ready")

}

// Attempts to run a function progressively through each client until one succeeds or they all fail.
func (m *BeaconClientManager) runFunction2(function bcFunction2) (interface{}, interface{}, error) {
	// Check if we can use the primary
	if m.primaryReady {
		// Try to run the function on the primary
		result1, result2, err := function(m.primaryBc)
		if err != nil {
			if m.isDisconnected(err) {
				// If it's disconnected, log it and try the fallback
				m.primaryReady = false
				if m.fallbackBc != nil {
					m.logger.Printlnf("WARNING: Primary Beacon client disconnected (%s), using fallback...", err.Error())
					return m.runFunction2(function)
				} else {
					m.logger.Printlnf("WARNING: Primary Beacon client disconnected (%s) and no fallback is configured.", err.Error())
					return nil, nil, fmt.Errorf("all Beacon clients failed")
				}
			}
			// If it's a different error, just return it
			return nil, nil, err
		}
		// If there's no error, return the result
		return result1, result2, nil
	}

	if m.fallbackReady {
		// Try to run the function on the fallback
		result1, result2, err := function(m.fallbackBc)
		if err != nil {
			if m.isDisconnected(err) {
				// If it's disconnected, log it and try the fallback
				m.logger.Printlnf("WARNING: Fallback Beacon client request failed (%s)", err.Error())
				m.fallbackReady = false
				return nil, nil, fmt.Errorf("all Beacon clients failed")
			}
			// If it's a different error, just return it
			return nil, nil, err
		}
		// If there's no error, return the result
		return result1, result2, nil
	}

	return nil, nil, fmt.Errorf("no Beacon clients were ready")

}

// Returns true if the error was a connection failure and a backup client is available
func (m *BeaconClientManager) isDisconnected(err error) bool {
	var sysErr syscall.Errno
	if errors.As(err, &sysErr) {
		return true
	}
	var netErr net.Error
	return errors.As(err, &netErr)
}
