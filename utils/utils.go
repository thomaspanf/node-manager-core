package utils

import (
	"context"
	"time"

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

// Sleeps for the specified time, but can break out if the provided context is cancelled.
// Returns true if the context is cancelled, false if it's not and the full period was slept.
func SleepWithCancel(ctx context.Context, duration time.Duration) bool {
	timer := time.NewTimer(duration)
	select {
	case <-ctx.Done():
		// Cancel occurred
		timer.Stop()
		return true

	case <-timer.C:
		// Duration has passed without a cancel
		return false
	}
}
