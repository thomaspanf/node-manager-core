package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/rocket-pool/node-manager-core/log"
)

type NetworkSocketApiServer struct {
	logger   *slog.Logger
	handlers []IHandler
	port     uint16
	socket   net.Listener
	server   http.Server
	router   *mux.Router
}

func NewNetworkSocketApiServer(logger *slog.Logger, port uint16, handlers []IHandler, baseRoute string, apiVersion string) (*NetworkSocketApiServer, error) {
	// Create the router
	router := mux.NewRouter()

	// Create the manager
	server := &NetworkSocketApiServer{
		logger:   logger,
		handlers: handlers,
		port:     port,
		router:   router,
		server: http.Server{
			Handler: router,
		},
	}

	// Register each route
	nmcRouter := router.Host(baseRoute).PathPrefix("/api/v" + apiVersion).Subrouter()
	for _, handler := range server.handlers {
		handler.RegisterRoutes(nmcRouter)
	}

	return server, nil
}

// Starts listening for incoming HTTP requests
func (s *NetworkSocketApiServer) Start(wg *sync.WaitGroup) error {
	// Create the socket
	socket, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", s.port))
	if err != nil {
		return fmt.Errorf("error creating socket: %w", err)
	}
	s.socket = socket

	// Start listening
	wg.Add(1)
	go func() {
		err := s.server.Serve(socket)
		if !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("error while listening for HTTP requests", log.Err(err))
		}
		wg.Done()
	}()

	return nil
}

// Stops the HTTP listener
func (s *NetworkSocketApiServer) Stop() error {
	err := s.server.Shutdown(context.Background())
	if err != nil {
		return fmt.Errorf("error stopping listener: %w", err)
	}
	return nil
}
