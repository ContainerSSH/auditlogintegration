package auditlogintegration_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"testing"

	"github.com/containerssh/auditlog"
	"github.com/containerssh/auditlog/codec/binary"
	"github.com/containerssh/auditlog/message"
	"github.com/containerssh/auditlog/storage/file"
	"github.com/containerssh/geoip"
	"github.com/containerssh/log"
	"github.com/containerssh/service"
	"github.com/containerssh/sshserver"
	"github.com/containerssh/structutils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"

	"github.com/containerssh/auditlogintegration"
)

func TestConnectMessages(t *testing.T) {
	logger, err := log.New(
		log.Config{
			Level: log.LevelDebug,
			Format: log.FormatText,
		},
		"audit",
		os.Stdout,
	)
	assert.NoError(t, err)

	dir, err := ioutil.TempDir("temp", "testcase")
	assert.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	lifecycle := createTestServer(t, dir, logger)

	onReady := make(chan struct{})
	lifecycle.OnRunning(
		func(_ service.Service, _ service.Lifecycle) {
			onReady <- struct{}{}
		},
	)
	go func() {
		_ = lifecycle.Run()
	}()

	<-onReady

	processClientConnection(t)

	lifecycle.Stop(context.Background())

	checkStoredAuditMessages(t, dir, logger)
}

func createTestServer(t *testing.T, dir string, logger log.Logger) service.Lifecycle {
	geoipLookup, err := geoip.New(
		geoip.Config{
			Provider: "dummy",
		},
	)
	assert.NoError(t, err)

	auditLogHandler, err := auditlogintegration.New(
		auditlog.Config{
			Enable:  true,
			Format:  auditlog.FormatBinary,
			Storage: auditlog.StorageFile,
			File: file.Config{
				Directory: dir,
			},
		},
		&backendHandler{},
		geoipLookup,
		logger,
	)
	assert.NoError(t, err)

	sshConfig := sshserver.Config{
		Listen: "127.0.0.1:2222",
	}
	structutils.Defaults(&sshConfig)
	assert.NoError(t, sshConfig.GenerateHostKey())

	srv, err := sshserver.New(
		sshConfig,
		auditLogHandler,
		logger,
	)
	assert.NoError(t, err)

	lifecycle := service.NewLifecycle(srv)
	return lifecycle
}

func processClientConnection(t *testing.T) {
	clientConfig := &ssh.ClientConfig{
		User: "foo",
		Auth: []ssh.AuthMethod{ssh.Password("bar")},
	}
	clientConfig.HostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}
	sshConnection, err := ssh.Dial("tcp", "127.0.0.1:2222", clientConfig)
	assert.NoError(t, err)
	session, err := sshConnection.NewSession()
	assert.NoError(t, err)

	stdin, stdout, err := createPipe(session)
	assert.NoError(t, err)
	assert.NoError(t, session.Shell())

	_, err = stdin.Write([]byte("Check 1, 2..."))
	assert.NoError(t, err)
	assert.NoError(t, stdin.Close())

	_, err = ioutil.ReadAll(stdout)
	assert.NoError(t, err)

	assert.NoError(t, sshConnection.Close())
}

func checkStoredAuditMessages(t *testing.T, dir string, logger log.Logger) {
	storage, err := file.NewStorage(
		file.Config{
			Directory: dir,
		}, logger,
	)
	assert.NoError(t, err)
	entryChannel, errChannel := storage.List()
	var logReader io.ReadCloser
	select {
	case err := <-errChannel:
		assert.NoError(t, err)
	case entry := <-entryChannel:
		logReader, err = storage.OpenReader(entry.Name)
		assert.NoError(t, err)
	}
	assert.NotNil(t, logReader)
	if logReader == nil {
		return
	}

	decoder := binary.NewDecoder()
	messageChannel, errorChannel := decoder.Decode(logReader)
	var messages []message.Message
	var errors []error
loop:
	for {
		select {
		case msg, ok := <-messageChannel:
			if !ok {
				break loop
			}
			messages = append(messages, msg)
		case err, ok := <-errorChannel:
			if !ok {
				break loop
			}
			errors = append(errors, err)
		}
	}
	assert.Empty(t, errors)
	assert.NotEmpty(t, messages)
	assert.Equal(t, message.TypeConnect, messages[0].MessageType)
	assert.Equal(t, message.TypeAuthPassword, messages[1].MessageType)
	assert.Equal(t, message.TypeAuthPasswordSuccessful, messages[2].MessageType)
	assert.Equal(t, message.TypeHandshakeSuccessful, messages[3].MessageType)
	assert.Equal(t, message.TypeNewChannelSuccessful, messages[4].MessageType)
	assert.Equal(t, message.TypeChannelRequestShell, messages[5].MessageType)
	assert.Equal(t, message.TypeExit, messages[6].MessageType)
	assert.Equal(t, message.TypeDisconnect, messages[7].MessageType)
}

func createPipe(session *ssh.Session) (io.WriteCloser, io.Reader, error) {
	stdin, err := session.StdinPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to request stdin (%w)", err)
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to request stdout (%w)", err)
	}
	return stdin, stdout, nil
}

type backendHandler struct {}

func (b *backendHandler) OnUnsupportedChannelRequest(_ uint64, _ string, _ []byte) {}

func (b *backendHandler) OnFailedDecodeChannelRequest(
	_ uint64,
	_ string,
	_ []byte,
	_ error,
) {}

func (b *backendHandler) OnEnvRequest(_ uint64, _ string, _ string) error {
	return fmt.Errorf("env requests are not supported")
}

func (b *backendHandler) OnPtyRequest(
	_ uint64,
	_ string,
	_ uint32,
	_ uint32,
	_ uint32,
	_ uint32,
	_ []byte,
) error {
	return fmt.Errorf("pty requests are not supported")
}

func (b *backendHandler) OnExecRequest(
	_ uint64,
	_ string,
	_ io.Reader,
	_ io.Writer,
	_ io.Writer,
	_ func(exitStatus sshserver.ExitStatus),
) error {
	return fmt.Errorf("exec requests are not supported")
}

func (b *backendHandler) OnShell(
	_ uint64,
	stdin io.Reader,
	stdout io.Writer,
	_ io.Writer,
	onExit func(exitStatus sshserver.ExitStatus),
) error {
	go func() {
		_, _ = ioutil.ReadAll(stdin)
		_, _ = stdout.Write([]byte("Hello world!"))
		onExit(0)
	}()
	return nil
}

func (b *backendHandler) OnSubsystem(
	_ uint64,
	_ string,
	_ io.Reader,
	_ io.Writer,
	_ io.Writer,
	_ func(exitStatus sshserver.ExitStatus),
) error {
	return fmt.Errorf("subsystem requests are not supported")
}

func (b *backendHandler) OnSignal(_ uint64, _ string) error {
	return fmt.Errorf("signals are not supported")
}

func (b *backendHandler) OnWindow(_ uint64, _ uint32, _ uint32, _ uint32, _ uint32) error {
	return fmt.Errorf("window requests are not supported")
}

func (b *backendHandler) OnUnsupportedGlobalRequest(_ uint64, _ string, _ []byte) {
}

func (b *backendHandler) OnUnsupportedChannel(_ uint64, _ string, _ []byte) {
}

func (b *backendHandler) OnSessionChannel(_ uint64, _ []byte) (
	channel sshserver.SessionChannelHandler,
	failureReason sshserver.ChannelRejection,
) {
	return b, nil
}

func (b *backendHandler) OnAuthPassword(username string, password []byte) (
	response sshserver.AuthResponse,
	reason error,
) {
	if username == "foo" && bytes.Equal(password, []byte("bar")) {
		return sshserver.AuthResponseSuccess, nil
	}
	return sshserver.AuthResponseFailure, nil
}

func (b *backendHandler) OnAuthPubKey(_ string, _ string) (response sshserver.AuthResponse, reason error) {
	return sshserver.AuthResponseFailure, nil
}

func (b *backendHandler) OnHandshakeFailed(_ error) {}

func (b *backendHandler) OnHandshakeSuccess(username string) (
	connection sshserver.SSHConnectionHandler,
	failureReason error,
) {
	return b, nil
}

func (b *backendHandler) OnDisconnect() {
}

func (b *backendHandler) OnReady() error {
	return nil
}

func (b *backendHandler) OnShutdown(_ context.Context) {

}

func (b *backendHandler) OnNetworkConnection(
	_ net.TCPAddr,
	_ string,
) (sshserver.NetworkConnectionHandler, error) {
	return b, nil
}
