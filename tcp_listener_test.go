package t3go

import "testing"

func TestNewTCPListener(t *testing.T) {

	listener, err := NewTCPListener(&TCPConfig{
		Port:        10080,
		Certificate: nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = listener.Listen()
	if err != nil {
		t.Fatal(err)
	}
}
