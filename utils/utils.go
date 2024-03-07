package utils

import (
	"github.com/rocket-pool/node-manager-core/wallet"
	"github.com/sethvargo/go-password/password"
)

// Check if the node wallet is ready for transacting
func IsWalletReady(status wallet.WalletStatus) bool {
	return status.Address.HasAddress &&
		status.Wallet.IsLoaded &&
		status.Address.NodeAddress == status.Wallet.WalletAddress
}

// Generates a random password
func GenerateRandomPassword() (string, error) {
	// Generate a random 32-character password
	password, err := password.Generate(32, 6, 6, false, false)
	if err != nil {
		return "", err
	}

	return password, nil
}
