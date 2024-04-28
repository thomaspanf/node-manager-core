package validator

import (
	"fmt"
	"strings"
	"sync"

	"github.com/rocket-pool/node-manager-core/beacon"
	"github.com/rocket-pool/node-manager-core/node/validator/keystore"
	types "github.com/wealdtech/go-eth2-types/v2"
)

type ValidatorManager struct {
	keystoreManagers map[string]keystore.IKeystoreManager
	lock             *sync.Mutex
}

func NewValidatorManager(validatorPath string) *ValidatorManager {
	mgr := &ValidatorManager{
		keystoreManagers: map[string]keystore.IKeystoreManager{
			"lighthouse": keystore.NewLighthouseKeystoreManager(validatorPath),
			"lodestar":   keystore.NewLodestarKeystoreManager(validatorPath),
			"nimbus":     keystore.NewNimbusKeystoreManager(validatorPath),
			"prysm":      keystore.NewPrysmKeystoreManager(validatorPath),
			"teku":       keystore.NewTekuKeystoreManager(validatorPath),
		},
		lock: &sync.Mutex{},
	}
	return mgr
}

// Stores a validator key into all of the manager's client keystores
func (m *ValidatorManager) StoreKey(key *types.BLSPrivateKey, derivationPath string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	for name, mgr := range m.keystoreManagers {
		err := mgr.StoreValidatorKey(key, derivationPath)
		if err != nil {
			pubkey := beacon.ValidatorPubkey(key.PublicKey().Marshal())
			return fmt.Errorf("error saving validator key %s (path %s) to the %s keystore: %w", pubkey.HexWithPrefix(), derivationPath, name, err)
		}
	}
	return nil
}

// Loads a validator key from the manager's client keystores
func (m *ValidatorManager) LoadKey(pubkey beacon.ValidatorPubkey) (*types.BLSPrivateKey, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	errors := []string{}
	// Try loading the key from all of the keystores, caching errors but not breaking on them
	for name := range m.keystoreManagers {
		key, err := m.keystoreManagers[name].LoadValidatorKey(pubkey)
		if err != nil {
			errors = append(errors, err.Error())
		}
		if key != nil {
			return key, nil
		}
	}

	if len(errors) > 0 {
		// If there were errors, return them
		return nil, fmt.Errorf("encountered the following errors while trying to load the key for validator %s:\n%s", pubkey.Hex(), strings.Join(errors, "\n"))
	} else {
		// If there were no errors, the key just didn't exist
		return nil, fmt.Errorf("couldn't find the key for validator %s in any of the validator manager's keystores", pubkey.Hex())
	}
}
