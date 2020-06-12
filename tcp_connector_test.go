package t3go

import (
	"fmt"
	"net"
	"testing"
)

func TestNewTCPConnector(t *testing.T) {
	nat, err := MappingOnPort("tcp", 10080)
	if err != nil {
		t.Fatal()
	}
	fmt.Println("dail with port", nat.ExtPort())
	c, err := NewTCPConnector(&TCPConfig{
		RemotePort:    10080,
		RemoteIP:      net.ParseIP("127.0.0.1"),
		LocalBindPort: nat.ExtPort(),
		Certificate:   nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Dial(); err != nil {
		t.Fatal(err)
	}
}
