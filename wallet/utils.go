package wallet

import "fmt"

// Check if the node wallet is ready for transacting
func IsWalletReady(status WalletStatus) bool {
	return status.Address.HasAddress &&
		status.Wallet.IsLoaded &&
		status.Address.NodeAddress == status.Wallet.WalletAddress
}

// Convert a derivation path type to an actual path value
func GetDerivationPath(pathType DerivationPath) (string, error) {
	// Parse the derivation path
	switch pathType {
	case DerivationPath_Default:
		return DefaultNodeKeyPath, nil
	case DerivationPath_LedgerLive:
		return LedgerLiveNodeKeyPath, nil
	case DerivationPath_Mew:
		return MyEtherWalletNodeKeyPath, nil
	default:
		return "", fmt.Errorf("[%s] is not a valid derivation path type", string(pathType))
	}
}
