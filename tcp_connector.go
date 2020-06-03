package t3go

import (
	"context"
	"fmt"
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

	resp, err := readReply(tcp)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", *resp)

	return nil
}

func readReply(conn net.Conn) (*Response, error) {
	resp := &Response{}
	rlt := make([]byte, 16)
	n, err := conn.Read(rlt)
	if err != nil {
		return nil, err
	}
	resp.Status = rlt[0]
	next := rlt[1]
	fmt.Println("rlt", rlt, "limit", n)
	if next == ByteTrue {
		tmp := make([]byte, maxByteSize)
		n, err := conn.Read(tmp)
		if err != nil {
			return nil, err
		}
		fmt.Println("tmp", string(tmp[:n]))
		resp.Data = tmp[:n]
	}
	return resp, err
}
