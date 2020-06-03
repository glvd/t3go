package t3go

import (
	"context"
	"github.com/portmapping/go-reuse"
	"net"
)

// TCPConnector ...
type TCPConnector struct {
	cfg    *TCPConfig
	ctx    context.Context
	cancel context.CancelFunc
}

// NewTCPConnector ...
func NewTCPConnector(cfg *TCPConfig) (*TCPConnector, error) {
	tcp := &TCPConnector{
		ctx:    nil,
		cancel: nil,
		cfg:    cfg,
	}
	tcp.ctx, tcp.cancel = context.WithCancel(context.TODO())
	return tcp, nil
}

// Dial ...
func (c *TCPConnector) Dial() error {
	local := &net.TCPAddr{IP: net.IPv4zero, Port: c.cfg.BindPort}
	remote := &net.TCPAddr{IP: c.cfg.RemoteIP, Port: c.cfg.Port}
	tcp, err := reuse.DialTCP("tcp", local, remote)
	if err != nil {
		return err
	}
	err = writeHead(tcp, &Head{
		Type:    RequestPing,
		Tunnel:  0,
		Version: 0,
	})
	if err != nil {
		return err
	}
	return nil
}
