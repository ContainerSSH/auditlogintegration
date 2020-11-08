package auditlogintegration

import (
	"github.com/containerssh/auditlog"
	"github.com/containerssh/sshserver"
)

type sshConnectionHandler struct {
	backend sshserver.SSHConnectionHandler
	audit   auditlog.Connection
}

func (s *sshConnectionHandler) OnUnsupportedGlobalRequest(requestType string, payload []byte) {
	//todo audit payload
	s.audit.OnGlobalRequestUnknown(requestType)
	s.backend.OnUnsupportedGlobalRequest(requestType, payload)
}

func (s *sshConnectionHandler) OnUnsupportedChannel(channelType string, extraData []byte) {
	//todo audit extraData
	s.audit.OnNewChannelFailed(channelType, "unsupported channel type")
	s.backend.OnUnsupportedChannel(channelType, extraData)
}

func (s *sshConnectionHandler) OnSessionChannel(extraData []byte) (channel sshserver.SessionChannelHandler, failureReason sshserver.ChannelRejection) {
	backend, err := s.backend.OnSessionChannel(extraData)
	if err != nil {
		return nil, err
	}
	auditChannel := s.audit.OnNewChannelSuccess("session")
	return &sessionChannelHandler{
		backend: backend,
		audit:   auditChannel,
	}, nil
}
