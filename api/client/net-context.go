package client

import (
	"context"
	"log/slog"
	"net"
	"net/http"
)

// The context passed into a requester
type NetworkRequesterContext struct {
	// The base address of the server
	address string

	// An HTTP client for sending requests
	client *http.Client

	// Logger to print debug messages to
	logger *slog.Logger
}

// Creates a new API client requester context for network-based servers
func NewNetworkRequesterContext(address string, log *slog.Logger) *NetworkRequesterContext {
	requesterContext := &NetworkRequesterContext{
		address: address,
		logger:  log,
		client: &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return net.Dial("tcp", address)
				},
			},
		},
	}

	return requesterContext
}

// Get the base of the address used for submitting server requests
func (r *NetworkRequesterContext) GetAddressBase() string {
	return r.address
}

// Get the logger for the context
func (r *NetworkRequesterContext) GetLogger() *slog.Logger {
	return r.logger
}

// Set the logger for the context
func (r *NetworkRequesterContext) SetLogger(logger *slog.Logger) {
	r.logger = logger
}

// Send an HTTP request to the server
func (r *NetworkRequesterContext) SendRequest(request *http.Request) (*http.Response, error) {
	return r.client.Do(request)
}
