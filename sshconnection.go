package auditlogintegration

import (
	"github.com/containerssh/auditlog"
	"github.com/containerssh/auditlog/message"
	"github.com/containerssh/sshserver"
)

type sshConnectionHandler struct {
	backend sshserver.SSHConnectionHandler
	audit   auditlog.Connection
}

func (s *sshConnectionHandler) OnUnsupportedGlobalRequest(requestID uint64, requestType string, payload []byte) {
	//todo audit payload
	s.audit.OnGlobalRequestUnknown(requestType)
	s.backend.OnUnsupportedGlobalRequest(requestID, requestType, payload)
}

func (s *sshConnectionHandler) OnUnsupportedChannel(channelID uint64, channelType string, extraData []byte) {
	//todo audit extraData
	s.audit.OnNewChannelFailed(message.MakeChannelID(channelID), channelType, "unsupported channel type")
	s.backend.OnUnsupportedChannel(channelID, channelType, extraData)
}

func (s *sshConnectionHandler) OnSessionChannel(channelID uint64, extraData []byte) (channel sshserver.SessionChannelHandler, failureReason sshserver.ChannelRejection) {
	backend, err := s.backend.OnSessionChannel(channelID, extraData)
	if err != nil {
		return nil, err
	}
	auditChannel := s.audit.OnNewChannelSuccess(message.MakeChannelID(channelID), "session")
	return &sessionChannelHandler{
		backend: backend,
		audit:   auditChannel,
	}, nil
}
