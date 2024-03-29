package config

type IConfig interface {
	IConfigSection

	GetApiLogFilePath() string

	GetDaemonLogFilePath() string

	GetNodeAddressFilePath() string

	GetWalletFilePath() string

	GetPasswordFilePath() string

	GetNetworkResources() *NetworkResources

	GetExecutionClientUrls() (string, string)

	GetBeaconNodeUrls() (string, string)
}
