package engine

import (
	"net"
	"os"
)

func newAgentConnection() (net.Conn, error) {
	sockFile := os.Getenv("SSH_AUTH_SOCK")
	return net.Dial("unix", sockFile)
}
