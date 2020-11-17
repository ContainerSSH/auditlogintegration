package auditlogintegration

import (
	"io"

	"github.com/containerssh/auditlog"
	"github.com/containerssh/sshserver"
)

type sessionChannelHandler struct {
	backend sshserver.SessionChannelHandler
	audit   auditlog.Channel
}

func (s *sessionChannelHandler) OnUnsupportedChannelRequest(requestID uint64, requestType string, payload []byte) {
	s.backend.OnUnsupportedChannelRequest(requestID, requestType, payload)
	s.audit.OnRequestUnknown(requestID, requestType, payload)
}

func (s *sessionChannelHandler) OnFailedDecodeChannelRequest(requestID uint64, requestType string, payload []byte, reason error) {
	s.backend.OnFailedDecodeChannelRequest(requestID, requestType, payload, reason)
	s.audit.OnRequestDecodeFailed(requestID, requestType, payload, reason.Error())
}

func (s *sessionChannelHandler) OnEnvRequest(requestID uint64, name string, value string) error {
	s.audit.OnRequestSetEnv(requestID, name, value)
	if err := s.backend.OnEnvRequest(requestID, name, value); err != nil {
		s.audit.OnRequestFailed(requestID, err)
		return err
	}
	return nil
}

func (s *sessionChannelHandler) OnExecRequest(
	requestID uint64,
	program string,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	onExit func(exitStatus sshserver.ExitStatus),
) error {
	stdin = s.audit.GetStdinProxy(stdin)
	stdout = s.audit.GetStdoutProxy(stdout)
	stderr = s.audit.GetStderrProxy(stderr)
	s.audit.OnRequestExec(requestID, program)
	if err := s.backend.OnExecRequest(
		requestID,
		program,
		stdin,
		stdout,
		stderr,
		func(exitStatus sshserver.ExitStatus) {
			s.audit.OnExit(uint32(exitStatus))
			onExit(exitStatus)
		},
	); err != nil {
		return err
	}
	return nil
}

func (s *sessionChannelHandler) OnPtyRequest(
	requestID uint64,
	term string,
	columns uint32,
	rows uint32,
	width uint32,
	height uint32,
	modeList []byte,
) error {
	s.audit.OnRequestPty(requestID, term, columns, rows, width, height, modeList)
	if err := s.backend.OnPtyRequest(requestID, term, columns, rows, width, height, modeList); err != nil {
		s.audit.OnRequestFailed(requestID, err)
		return err
	}
	return nil

}

func (s *sessionChannelHandler) OnShell(
	requestID uint64,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	onExit func(exitStatus sshserver.ExitStatus),
) error {
	stdin = s.audit.GetStdinProxy(stdin)
	stdout = s.audit.GetStdoutProxy(stdout)
	stderr = s.audit.GetStderrProxy(stderr)

	s.audit.OnRequestShell(requestID)
	if err := s.backend.OnShell(
		requestID,
		stdin,
		stdout,
		stderr,
		func(exitStatus sshserver.ExitStatus) {
			s.audit.OnExit(uint32(exitStatus))
			onExit(exitStatus)
		},
	); err != nil {
		s.audit.OnRequestFailed(requestID, err)
		return err
	}
	return nil
}

func (s *sessionChannelHandler) OnSignal(requestID uint64, signal string) error {
	s.audit.OnRequestSignal(requestID, signal)
	if err := s.backend.OnSignal(requestID, signal); err != nil {
		s.audit.OnRequestFailed(requestID, err)
		return err
	}
	return nil
}

func (s *sessionChannelHandler) OnSubsystem(
	requestID uint64,
	subsystem string,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	onExit func(exitStatus sshserver.ExitStatus),
) error {
	stdin = s.audit.GetStdinProxy(stdin)
	stdout = s.audit.GetStdoutProxy(stdout)
	stderr = s.audit.GetStderrProxy(stderr)
	s.audit.OnRequestSubsystem(requestID, subsystem)
	if err := s.backend.OnSubsystem(
		requestID,
		subsystem,
		stdin,
		stdout,
		stderr,
		func(exitStatus sshserver.ExitStatus) {
			s.audit.OnExit(uint32(exitStatus))
			onExit(exitStatus)
		},
	); err != nil {
		s.audit.OnRequestFailed(requestID, err)
		return err
	}
	return nil
}

func (s *sessionChannelHandler) OnWindow(
	requestID uint64,
	columns uint32,
	rows uint32,
	width uint32,
	height uint32,
) error {
	s.audit.OnRequestWindow(requestID, columns, rows, width, height)
	if err := s.backend.OnWindow(requestID, columns, rows, width, height); err != nil {
		s.audit.OnRequestFailed(requestID, err)
		return err
	}
	return nil
}
