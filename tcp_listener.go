package t3go

import (
	"context"
	"crypto/tls"
	"fmt"
	ants2 "github.com/panjf2000/ants"
	"net"

	"github.com/panjf2000/ants/v2"
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
	pool   *ants2.PoolWithFunc
}

// NewTCPListener ...
func NewTCPListener(cfg *TCPConfig) (*TCPListener, error) {
	antsPool, err := ants.NewPoolWithFunc(ants.DefaultAntsPoolSize, tcpListenHandler, ants.WithNonblocking(false))
	if err != nil {
		return nil, err
	}
	tcp := &TCPListener{
		ctx:    nil,
		cancel: nil,
		pool:   antsPool,
		cfg:    cfg,
	}
	tcp.ctx, tcp.cancel = context.WithCancel(context.TODO())

	return tcp, nil
}

func tcpListenHandler(i interface{}) {
	conn, b := i.(net.Conn)
	if !b {
		return
	}
	var err error
	defer func() {
		if err != nil {
			conn.Close()
		}
	}()
	head, err := readHead(conn)
	if err != nil {
		return
	}
	err = processRun(head.Type, conn)
}

func processRun(types uint8, conn net.Conn) error {
	switch types {
	case RequestPing:
	case RequestConnect:
	}
	return fmt.Errorf("not supported")
}

type Head struct {
	Type    uint8 `json:"type"`
	Tunnel  uint8 `json:"tunnel"`
	Version uint8 `json:"version"`
}

func readHead(conn net.Conn) (*Head, error) {
	head := make([]byte, 16)
	read, err := conn.Read(head)
	if err != nil {
		return nil, err
	}
	if read < 8 {
		return nil, fmt.Errorf("wrong head size")
	}
	h := Head{
		Type:    head[0],
		Tunnel:  head[1],
		Version: head[2],
	}
	return &h, nil
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
		err = l.pool.Invoke(conn)
		if err != nil {
			continue
		}
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
