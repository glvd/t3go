package t3go

import (
	"net"
	"testing"
)

func TestNewTCPConnector(t *testing.T) {
	port, err := MappingOnPort("tcp", 10080)
	if err != nil {
		t.Fatal()
	}
	c, err := NewTCPConnector(&TCPConfig{
		Port:        10080,
		RemoteIP:    net.ParseIP("127.0.0.1"),
		BindPort:    12306,
		Certificate: nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Dial(); err != nil {
		t.Fatal(err)
	}
}
