package auditlogintegration

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/containerssh/auditlog"
	"github.com/containerssh/sshserver"
)

type handler struct {
	logger  auditlog.Logger
	backend sshserver.Handler
}

func (h *handler) OnReady() error {
	return h.backend.OnReady()
}

func (h *handler) OnShutdown(shutdownContext context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		h.backend.OnShutdown(shutdownContext)
	}()
	go func() {
		defer wg.Done()
		h.logger.Shutdown(shutdownContext)
	}()
	wg.Wait()
}

func (h *handler) OnNetworkConnection(client net.TCPAddr, connectionID []byte) (sshserver.NetworkConnectionHandler, error) {
	backend, err := h.backend.OnNetworkConnection(client, connectionID)
	if err != nil {
		return nil, err
	}
	auditConnection, err := h.logger.OnConnect(connectionID, client)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to initialize audit logger for connection from %s (%w)",
			client.IP.String(),
			err,
		)
	}

	return &networkConnectionHandler{
		backend: backend,
		audit:   auditConnection,
	}, nil
}
