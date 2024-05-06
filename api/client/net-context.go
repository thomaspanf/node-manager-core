package client

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
)

// The context passed into a requester
type NetworkRequesterContext struct {
	// The base address and route for API calls
	apiUrl *url.URL

	// An HTTP client for sending requests
	client *http.Client

	// Logger to print debug messages to
	logger *slog.Logger

	// Tracer for HTTP requests
	tracer *httptrace.ClientTrace
}

// Creates a new API client requester context for network-based
// traceOpts is optional. If nil, it will not be used.
func NewNetworkRequesterContext(apiUrl *url.URL, log *slog.Logger, tracer *httptrace.ClientTrace) *NetworkRequesterContext {
	requesterContext := &NetworkRequesterContext{
		apiUrl: apiUrl,
		logger: log,
		tracer: tracer,
		client: &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return net.Dial("tcp", apiUrl.Host)
				},
			},
		},
	}

	return requesterContext
}

// Get the base of the address used for submitting server requests
func (r *NetworkRequesterContext) GetAddressBase() string {
	return r.apiUrl.String()
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
	if r.tracer != nil {
		request = request.WithContext(httptrace.WithClientTrace(request.Context(), r.tracer))
	}
	return r.client.Do(request)
}
