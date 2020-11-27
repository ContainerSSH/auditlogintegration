package auditlogintegration

import (
	"fmt"

	"github.com/containerssh/geoip/geoipprovider"
	"github.com/containerssh/log"
	"github.com/containerssh/sshserver"

	"github.com/containerssh/auditlog"

	"github.com/containerssh/configuration"
)

// New creates a new handler based on the application config and the required dependencies. If audit logging is not
// enabled the backend will be returned directly.
//goland:noinspection GoUnusedExportedFunction
func New(
	config configuration.AppConfig,
	backend sshserver.Handler,
	geoIPLookupProvider geoipprovider.LookupProvider,
	logger log.Logger,
) (sshserver.Handler, error) {
	if !config.Audit.Enable {
		return backend, nil
	}

	auditLogger, err := auditlog.New(
		config.Audit,
		geoIPLookupProvider,
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create audit logger (%v)", err)
	}

	handler := NewHandler(
		backend,
		auditLogger,
	)
	return handler, nil
}
