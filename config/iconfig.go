package config

type IConfig interface {
	IConfigSection

	GetNodeAddressFilePath() string

	GetWalletFilePath() string

	GetPasswordFilePath() string

	GetNetworkResources() *NetworkResources

	GetExecutionClientUrls() (string, string)

	GetBeaconNodeUrls() (string, string)
}
