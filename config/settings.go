package config

import "time"

const (
	EventLogInterval int    = 1000
	DockerApiVersion string = "1.40"

	// Wallet
	UserAddressFilename    string = "address"
	UserWalletDataFilename string = "wallet"
	UserPasswordFilename   string = "password"

	// HTTP
	ClientTimeout time.Duration = 8 * time.Second
)
