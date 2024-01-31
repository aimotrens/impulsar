package sshexecutor

import (
	"net"

	"github.com/Microsoft/go-winio"
)

func newAgentConnection() (net.Conn, error) {
	pipeName := "\\\\.\\pipe\\openssh-ssh-agent"
	return winio.DialPipe(pipeName, nil)
}
