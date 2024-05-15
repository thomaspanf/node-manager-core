package types

// This is a wrapper for the EC / BN status report
type ClientStatus struct {
	IsWorking    bool    `json:"isWorking"`
	IsSynced     bool    `json:"isSynced"`
	SyncProgress float64 `json:"syncProgress"`
	ChainId      uint    `json:"networkId"`
	Error        string  `json:"error"`
}

// This is a wrapper for the manager's overall status report
type ClientManagerStatus struct {
	PrimaryClientStatus  ClientStatus `json:"primaryEcStatus"`
	FallbackEnabled      bool         `json:"fallbackEnabled"`
	FallbackClientStatus ClientStatus `json:"fallbackEcStatus"`
}
