package client

import (
	"context"
	"log/slog"
	"net"
	"net/http"
)

// The context passed into a requester
type RequesterContext struct {
	// The path to the socket to send requests to
	socketPath string

	// An HTTP client for sending requests
	client *http.Client

	// Logger to print debug messages to
	logger *slog.Logger

	// The base route for the client to send requests to (<http://<base>/<route>/<method>)
	base string
}

// Creates a new API client requester context
func NewRequesterContext(baseRoute string, socketPath string, log *slog.Logger) *RequesterContext {
	requesterContext := &RequesterContext{
		socketPath: socketPath,
		base:       baseRoute,
		logger:     log,
		client: &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return net.Dial("unix", socketPath)
				},
			},
		},
	}

	return requesterContext
}

// Set the logger for the context
func (r *RequesterContext) SetLogger(logger *slog.Logger) {
	r.logger = logger
}
