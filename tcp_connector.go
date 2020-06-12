package t3go

import (
	"context"
	"github.com/portmapping/go-reuse"
	"net"
	"time"
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
	local := &net.TCPAddr{IP: c.cfg.LocalIP, Port: c.cfg.LocalBindPort}
	remote := &net.TCPAddr{IP: c.cfg.RemoteIP, Port: c.cfg.RemotePort}
	tcp, err := reuse.DialTCP("tcp", local, remote)
	if err != nil {
		return err
	}
	go receiveHandle(tcp)
	go sendHandle(tcp)

	//err = writeHead(tcp, &Head{
	//	Type:    RequestPing,
	//	Tunnel:  0,
	//	Version: 0,
	//})
	//if err != nil {
	//	return err
	//}
	//
	//resp, err := readReply(tcp)
	//if err != nil {
	//	return err
	//}
	//fmt.Printf("%+v\n", *resp)
	time.Sleep(30 * time.Minute)
	return nil
}

func readReply(conn net.Conn) (*Response, error) {
	rlt := make([]byte, 16)
	_, err := conn.Read(rlt)
	if err != nil {
		return nil, err
	}
	resp := &Response{
		Status: rlt[0],
	}
	if rlt[1] == ByteTrue {
		tmp := make([]byte, maxByteSize)
		n, err := conn.Read(tmp)
		if err != nil {
			return nil, err
		}
		resp.Data = tmp[:n]
	}
	return resp, err
}
