package client

import "time"

type StandardHttpClient struct {
	*StandardClient
}

// Create a new client instance
func NewStandardHttpClient(providerAddress string, timeout time.Duration) *StandardHttpClient {
	provider := NewBeaconHttpProvider(providerAddress, timeout)
	return &StandardHttpClient{
		StandardClient: NewStandardClient(provider),
	}
}
