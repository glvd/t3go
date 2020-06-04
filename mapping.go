package t3go

import (
	"errors"
	"fmt"
	"time"

	"github.com/libp2p/go-nat"
	"go.uber.org/atomic"
)

const description = "mapping_port"

// NAT ...
type NAT interface {
	Port() int
	ExtPort() int
	Protocol() string
	StopMapping() error
	StartMapping() error
}

// MapConfig ...
type MapConfig struct {
	Timeout     time.Duration
	Description string
}

type natClient struct {
	nat      nat.NAT
	stop     *atomic.Bool
	port     int
	protocol string
	extport  int
	cfg      *MapConfig
}

// Protocol ...
func (c *natClient) Protocol() string {
	return c.protocol
}

// Port ...
func (c *natClient) Port() int {
	return c.port
}

// ExtPort ...
func (c *natClient) ExtPort() int {
	return c.extport
}

// StopMapping ...
func (c *natClient) StopMapping() error {
	c.stop.Store(true)
	return nil
}

// StartMapping ...
func (c *natClient) StartMapping() error {
	return c.mapping()
}

// MapConfigOptions ...
type MapConfigOptions func(c *MapConfig)

// MappingOnPort ...
func MappingOnPort(protocol string, port int, opts ...MapConfigOptions) (NAT, error) {
	n, err := nat.DiscoverGateway()
	if err != nil {
		return nil, fmt.Errorf("failed discover gateway on mapping: %w", err)
	}
	c := defaultMapConfig()
	for _, opt := range opts {
		opt(c)
	}
	cli := &natClient{
		stop:     atomic.NewBool(true),
		cfg:      c,
		nat:      n,
		port:     port,
		protocol: protocol,
	}
	err = cli.mapping()
	return cli, err
}

func (c *natClient) mapping() (err error) {
	if !c.stop.Load() {
		return errors.New("nat was on mapping")
	}
	c.stop.Store(false)
	c.extport, err = c.nat.AddPortMapping(c.protocol, c.port, description, c.cfg.Timeout*time.Second)
	if err != nil {
		return fmt.Errorf("port mapping failed: %w", err)
	}

	go func() {
		t := time.NewTicker(30 * time.Second)
		defer func() {
			t.Stop()
			if e := recover(); e != nil {
				fmt.Println("panic error on mapping:", e)
			}
		}()

		for {
			//check mapping every 30 second
			<-t.C
			if c.stop.Load() {
				return
			}
			_, err = c.nat.AddPortMapping(c.protocol, c.port, description, c.cfg.Timeout*time.Second)
			if err != nil {
				panic(err)
			}

		}
	}()

	return nil
}

func defaultMapConfig() *MapConfig {
	return &MapConfig{
		Timeout:     60,
		Description: description,
	}
}
