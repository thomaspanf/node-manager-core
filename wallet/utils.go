package wallet

// Check if the node wallet is ready for transacting
func IsWalletReady(status WalletStatus) bool {
	return status.Address.HasAddress &&
		status.Wallet.IsLoaded &&
		status.Address.NodeAddress == status.Wallet.WalletAddress
}
