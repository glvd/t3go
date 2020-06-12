package t3go

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"net"

	"github.com/panjf2000/ants/v2"
	"github.com/portmapping/go-reuse"
)

const maxByteSize = 65520
const (
	// RequestPing ...
	RequestPing = 0x01
	// RequestConnect ...
	RequestConnect = 0x02
)
const (
	// ResponseFailed ...
	ResponseFailed = 0x00
	// ResponseSuccess ...
	ResponseSuccess = 0x01
)

const (
	// ByteFalse ...
	ByteFalse = 0x00
	// ByteTrue ...
	ByteTrue = 0x01
)

// TCPConfig ...
type TCPConfig struct {
	RemotePort    int
	RemoteIP      net.IP
	LocalIP       net.IP
	LocalPort     int
	LocalBindPort int
	Certificate   []tls.Certificate
}

// TCPListener ...
type TCPListener struct {
	cfg    *TCPConfig
	ctx    context.Context
	cancel context.CancelFunc
	pool   *ants.PoolWithFunc
}

// Head ...
type Head struct {
	Type    uint8
	Tunnel  uint8
	Version uint8
}

// Response ...
type Response struct {
	Status uint8
	Data   []byte
}

// HasData ...
func (r Response) HasData() uint8 {
	if r.Data == nil {
		return ByteFalse
	}
	return ByteTrue
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

func receiveHandle(conn net.Conn) error {
	buf := make([]byte, 1024)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			return err
		}
		fmt.Println("read", string(buf))
	}
}

func sendHandle(conn net.Conn) error {
	buf := make([]byte, 1024)
	copy(buf, "hello world")
	for {
		_, err := conn.Write(buf)
		if err != nil {
			return err
		}
		//time.Sleep(1 * time.Second)
	}
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
	fmt.Println("address", conn.RemoteAddr())
	//head, err := readHead(conn)
	//if err != nil {
	//	return
	//}
	//err = processRun(head.Type, conn)
	go receiveHandle(conn)
	go sendHandle(conn)
	time.Sleep(30 * time.Minute)
}

func processRun(types uint8, conn net.Conn) error {
	switch types {
	case RequestPing:
		return reply(conn, &Response{
			Status: ResponseSuccess,
			Data:   []byte("pong"),
		})
	case RequestConnect:

	}
	return fmt.Errorf("not supported")
}

func reply(conn net.Conn, resp *Response) error {
	rlt := make([]byte, 16)
	rlt[0] = resp.Status
	rlt[1] = resp.HasData()
	_, err := conn.Write(rlt)
	if err != nil {
		return err
	}
	if resp.HasData() == ByteTrue {
		_, err = conn.Write(resp.Data)
	}
	return err
}

func writeHead(conn net.Conn, head *Head) error {
	h := make([]byte, 16)
	h[0] = head.Type
	h[1] = head.Tunnel
	h[2] = head.Version
	_, err := conn.Write(h)
	return err
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

// Listen ...
func (l *TCPListener) Listen() (err error) {
	addr := &net.TCPAddr{
		IP:   l.cfg.LocalIP,
		Port: l.cfg.RemotePort,
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

// Stop ...
func (l *TCPListener) Stop() {
	if l.cancel != nil {
		l.cancel()
	}
}

// NewConnector ...
func (l *TCPListener) NewConnector(conn net.Conn) {
	defer conn.Close()
}
