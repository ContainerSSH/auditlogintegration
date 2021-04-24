package auditlogintegration

import (
	"github.com/containerssh/auditlog"
	sshserver "github.com/containerssh/sshserver/v2"
)

// NewHandler creates a new audit logging handler that logs all events as configured, and passes request to a provided backend.
func NewHandler(backend sshserver.Handler, logger auditlog.Logger) sshserver.Handler {
	return &handler{
		backend: backend,
		logger:  logger,
	}
}
