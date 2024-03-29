package config

type IConfig interface {
	IConfigSection

	GetApiLogFilePath() string

	GetTasksLogFilePath() string

	GetNodeAddressFilePath() string

	GetWalletFilePath() string

	GetPasswordFilePath() string

	GetNetworkResources() *NetworkResources

	GetExecutionClientUrls() (string, string)

	GetBeaconNodeUrls() (string, string)
}
