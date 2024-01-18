package engine

import (
	"github.com/Microsoft/go-winio"
	"net"
)

func newAgentConnection() (net.Conn, error) {
	pipeName := "\\\\.\\pipe\\openssh-ssh-agent"
	return winio.DialPipe(pipeName, nil)
}
