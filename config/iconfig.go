package config

// NMC servers typically provide some kind of persistent configuration; it must implement this interface.
type IConfig interface {
	IConfigSection

	// The path to use for the API log file
	GetApiLogFilePath() string

	// The path to use for the tasks log file
	GetTasksLogFilePath() string

	// The path to use for the node address file
	GetNodeAddressFilePath() string

	// The path to use for the wallet keystore file
	GetWalletFilePath() string

	// The path to use for the wallet keystore's password file
	GetPasswordFilePath() string

	// The resources for the selected network
	GetNetworkResources() *NetworkResources

	// The URLs for the Execution clients to use
	GetExecutionClientUrls() (string, string)

	// The URLs for the Beacon nodes to use
	GetBeaconNodeUrls() (string, string)
}
