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

func (s *sessionChannelHandler) OnUnsupportedChannelRequest(requestType string, payload []byte) {
	//todo log unsuccessful requests
	s.backend.OnUnsupportedChannelRequest(requestType, payload)
	//todo audit log payload
	s.audit.OnRequestUnknown(requestType)
}

func (s *sessionChannelHandler) OnFailedDecodeChannelRequest(requestType string, payload []byte, reason error) {
	//todo log unsuccessful requests
	s.backend.OnFailedDecodeChannelRequest(requestType, payload, reason)
	//todo audit log payload
	s.audit.OnRequestDecodeFailed(requestType, reason.Error())
}

func (s *sessionChannelHandler) OnEnvRequest(name string, value string) error {
	//todo log unsuccessful requests
	if err := s.backend.OnEnvRequest(name, value); err != nil {
		return err
	}
	s.audit.OnRequestSetEnv(name, value)
	return nil
}

func (s *sessionChannelHandler) OnExecRequest(program string, stdin io.Reader, stdout io.Writer, stderr io.Writer, onExit func(exitStatus uint32)) error {
	//todo log unsuccessful requests
	stdin = s.audit.GetStdinProxy(stdin)
	stdout = s.audit.GetStdoutProxy(stdout)
	stderr = s.audit.GetStderrProxy(stderr)
	//todo audit log exit
	if err := s.backend.OnExecRequest(program, stdin, stdout, stderr, onExit); err != nil {
		return err
	}
	s.audit.OnRequestExec(program)
	return nil
}

func (s *sessionChannelHandler) OnPtyRequest(term string, columns uint32, rows uint32, width uint32, height uint32, modeList []byte) error {
	//todo log unsuccessful requests
	if err := s.backend.OnPtyRequest(term, columns, rows, width, height, modeList); err != nil {
		return err
	}
	//todo audit log term, dimenstions, and modelist
	s.audit.OnRequestPty(uint(columns), uint(rows))
	return nil

}

func (s *sessionChannelHandler) OnShell(stdin io.Reader, stdout io.Writer, stderr io.Writer, onExit func(exitStatus uint32)) error {
	//todo log unsuccessful requests
	stdin = s.audit.GetStdinProxy(stdin)
	stdout = s.audit.GetStdoutProxy(stdout)
	stderr = s.audit.GetStderrProxy(stderr)
	//todo audit log exit
	if err := s.backend.OnShell(stdin, stdout, stderr, onExit); err != nil {
		return err
	}
	s.audit.OnRequestShell()
	return nil
}

func (s *sessionChannelHandler) OnSignal(signal string) error {
	//todo log unsuccessful requests
	if err := s.backend.OnSignal(signal); err != nil {
		return err
	}
	s.audit.OnRequestSignal(signal)
	return nil
}

func (s *sessionChannelHandler) OnSubsystem(subsystem string, stdin io.Reader, stdout io.Writer, stderr io.Writer, onExit func(exitStatus uint32)) error {
	//todo log unsuccessful requests
	stdin = s.audit.GetStdinProxy(stdin)
	stdout = s.audit.GetStdoutProxy(stdout)
	stderr = s.audit.GetStderrProxy(stderr)
	//todo audit log exit
	if err := s.backend.OnSubsystem(subsystem, stdin, stdout, stderr, onExit); err != nil {
		return err
	}
	s.audit.OnRequestSubsystem(subsystem)
	return nil
}

func (s *sessionChannelHandler) OnWindow(columns uint32, rows uint32, width uint32, height uint32) error {
	//todo log unsuccessful requests
	if err := s.backend.OnWindow(columns, rows, width, height); err != nil {
		return err
	}
	//todo log window dimensions
	s.audit.OnRequestWindow(uint(columns), uint(rows))
	return nil
}
