package auditlogintegration

import (
	"github.com/containerssh/auditlog"
	"github.com/containerssh/sshserver"
)

// NewHandler creates a new audit logging handler that logs all events as configured, and passes request to a provided backend.
//goland:noinspection GoUnusedExportedFunction
func NewHandler(backend sshserver.Handler, logger auditlog.Logger) sshserver.Handler {
	return &handler{
		backend: backend,
		logger:  logger,
	}
}
