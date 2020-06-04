package t3go

import (
	"fmt"
	"time"

	"github.com/libp2p/go-nat"
	"go.uber.org/atomic"
)

const description = "mapping_port"

// NAT ...
type NAT interface {
}

// MapConfig ...
type MapConfig struct {
	Timeout     time.Duration
	Description string
}

type natClient struct {
	nat      nat.NAT
	stop     atomic.Bool
	port     int
	protocol string
	extport  int
	cfg      *MapConfig
}

// MapConfigOptions ...
type MapConfigOptions func(c *MapConfig)

// MappingOnPort ...
func MappingOnPort(protocol string, port int, opts ...MapConfigOptions) (NAT, error) {
	n, err := nat.DiscoverGateway()
	if err != nil {
		panic(err)
	}
	c := defaultMapConfig()
	for _, opt := range opts {
		opt(c)
	}
	extport, err := n.AddPortMapping(protocol, port, c.Description, c.Timeout*time.Second)
	if err != nil {
		return nil, err
	}
	cli := &natClient{
		cfg:      c,
		nat:      n,
		port:     port,
		protocol: protocol,
		extport:  extport,
	}
	return cli, nil
}

func (c *natClient) mapping() (err error) {
	c.stop.Store(false)
	c.extport, err = c.nat.AddPortMapping(c.protocol, c.port, description, c.cfg.Timeout*time.Second)
	if err != nil {
		return err
	}

	go func() {
		t := time.NewTicker(30 * time.Second)
		defer func() {
			t.Stop()
			if e := recover(); e != nil {
				fmt.Println("panic error:", e)
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
