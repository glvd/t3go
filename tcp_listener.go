package t3go

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"github.com/portmapping/go-reuse"
)

type TCPConfig struct {
	Port        int
	Certificate []tls.Certificate
}

type TCPListener struct {
	cfg    *TCPConfig
	ctx    context.Context
	cancel context.CancelFunc
}

// NewTCPListener ...
func NewTCPListener(cfg *TCPConfig) *TCPListener {
	tcp := &TCPListener{
		ctx:    nil,
		cancel: nil,
		cfg:    cfg,
	}
	tcp.ctx, tcp.cancel = context.WithCancel(context.TODO())

	return tcp
}

func (l *TCPListener) Listen() (err error) {
	addr := &net.TCPAddr{
		IP:   net.IPv4zero,
		Port: l.cfg.Port,
	}

	var lis net.Listener
	if l.cfg.Certificate != nil {
		lis, err = reuse.ListenTLS("tcp", addr.String(), &tls.Config{
			Certificates: l.cfg.Certificate,
		})
	} else {
		lis, err = reuse.ListenTCP("tcp", addr)
	}
	if err != nil {
		return err
	}
	fmt.Println("listen tcp on address:", addr.String())
	for {
		conn, err := lis.Accept()
		if err != nil {
			continue
		}
		go l.NewConnector(conn)
	}
}

func (l *TCPListener) Stop() {
	if l.cancel != nil {
		l.cancel()
	}
}

func (l *TCPListener) NewConnector(conn net.Conn) {
	defer conn.Close()
}
