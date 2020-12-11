package auditlogintegration

import (
	"github.com/containerssh/auditlog"
	"github.com/containerssh/sshserver"
)

type networkConnectionHandler struct {
	backend sshserver.NetworkConnectionHandler
	audit   auditlog.Connection
}

func (n *networkConnectionHandler) OnAuthPassword(
	username string,
	password []byte,
) (response sshserver.AuthResponse, reason error) {
	n.audit.OnAuthPassword(username, password)
	response, reason = n.backend.OnAuthPassword(username, password)
	switch response {
	case sshserver.AuthResponseSuccess:
		n.audit.OnAuthPasswordSuccess(username, password)
	case sshserver.AuthResponseFailure:
		n.audit.OnAuthPasswordFailed(username, password)
	case sshserver.AuthResponseUnavailable:
		if reason != nil {
			n.audit.OnAuthPasswordBackendError(username, password, reason.Error())
		} else {
			n.audit.OnAuthPasswordBackendError(username, password, "")
		}
	}
	return response, reason
}

func (n *networkConnectionHandler) OnAuthPubKey(
	username string,
	pubKey string,
) (
	response sshserver.AuthResponse,
	reason error,
) {
	n.audit.OnAuthPubKey(username, pubKey)
	response, reason = n.backend.OnAuthPubKey(username, pubKey)
	switch response {
	case sshserver.AuthResponseSuccess:
		n.audit.OnAuthPubKeySuccess(username, pubKey)
	case sshserver.AuthResponseFailure:
		n.audit.OnAuthPubKeyFailed(username, pubKey)
	case sshserver.AuthResponseUnavailable:
		if reason != nil {
			n.audit.OnAuthPubKeyBackendError(username, pubKey, reason.Error())
		} else {
			n.audit.OnAuthPubKeyBackendError(username, pubKey, "")
		}
	}
	return response, reason
}

func (n *networkConnectionHandler) OnHandshakeFailed(reason error) {
	n.backend.OnHandshakeFailed(reason)
	n.audit.OnHandshakeFailed(reason.Error())
}

func (n *networkConnectionHandler) OnHandshakeSuccess(
	username string,
) (
	connection sshserver.SSHConnectionHandler,
	failureReason error,
) {
	n.audit.OnHandshakeSuccessful(username)
	backend, err := n.backend.OnHandshakeSuccess(username)
	if err != nil {
		return nil, err
	}
	return &sshConnectionHandler{
		backend: backend,
		audit:   n.audit,
	}, nil
}

func (n *networkConnectionHandler) OnDisconnect() {
	n.audit.OnDisconnect()
	n.backend.OnDisconnect()
}
