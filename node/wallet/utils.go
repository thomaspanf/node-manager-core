package wallet

import (
	"fmt"

	"github.com/rocket-pool/node-manager-core/wallet"
)

// Convert a derivation path type to an actual path value
func GetDerivationPath(pathType wallet.DerivationPath) (string, error) {
	// Parse the derivation path
	switch pathType {
	case wallet.DerivationPath_Default:
		return DefaultNodeKeyPath, nil
	case wallet.DerivationPath_LedgerLive:
		return LedgerLiveNodeKeyPath, nil
	case wallet.DerivationPath_Mew:
		return MyEtherWalletNodeKeyPath, nil
	default:
		return "", fmt.Errorf("[%s] is not a valid derivation path type", string(pathType))
	}
}
