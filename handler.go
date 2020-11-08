package auditlogintegration

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/containerssh/auditlog"
	"github.com/containerssh/sshserver"
	"github.com/google/uuid"
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

func (h *handler) OnNetworkConnection(ip net.Addr) (sshserver.NetworkConnectionHandler, error) {
	backend, err := h.backend.OnNetworkConnection(ip)
	if err != nil {
		return nil, err
	}
	connectionId, err := uuid.New().MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to generate unique connection ID for %s (%w)", ip.String(), err)
	}
	auditConnection, err := h.logger.OnConnect(connectionId, *ip.(*net.TCPAddr))
	if err != nil {
		return nil, fmt.Errorf(
			"failed to initialize audit logger for connection from %s  (%w)",
			ip.String(),
			err,
		)
	}

	return &networkConnectionHandler{
		backend: backend,
		audit:   auditConnection,
	}, nil
}
