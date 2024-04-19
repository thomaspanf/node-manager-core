package validator

import (
	"fmt"
	"strconv"

	"github.com/rocket-pool/node-manager-core/beacon"
	"github.com/rocket-pool/node-manager-core/beacon/ssz_types"
	eth2types "github.com/wealdtech/go-eth2-types/v2"
)

// Get a voluntary exit message signature for a given validator key and index
func GetSignedExitMessage(validatorKey *eth2types.BLSPrivateKey, validatorIndex string, epoch uint64, signatureDomain []byte) (beacon.ValidatorSignature, error) {
	// Parse the validator index
	indexNum, err := strconv.ParseUint(validatorIndex, 10, 64)
	if err != nil {
		return beacon.ValidatorSignature{}, fmt.Errorf("error parsing validator index (%s): %w", validatorIndex, err)
	}
	// Build voluntary exit message
	exitMessage := ssz_types.VoluntaryExit{
		Epoch:          epoch,
		ValidatorIndex: indexNum,
	}
	// Get object root
	or, err := exitMessage.HashTreeRoot()
	if err != nil {
		return beacon.ValidatorSignature{}, err
	}
	// Get signing root
	sr := ssz_types.SigningRoot{
		ObjectRoot: or[:],
		Domain:     signatureDomain,
	}

	srHash, err := sr.HashTreeRoot()
	if err != nil {
		return beacon.ValidatorSignature{}, err
	}
	// Sign message
	signature := validatorKey.Sign(srHash[:]).Marshal()
	// Return
	return beacon.ValidatorSignature(signature), nil

}
