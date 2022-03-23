package masky

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/wkj9893/masky/internal/dns"
	"github.com/wkj9893/masky/internal/tls"
)

type Client struct {
	sync.RWMutex

	config ClientConfig
	cache  map[string]string // isocode cache
	// session quic.EarlySession
}

func NewClient(config ClientConfig) (*Client, error) {
	if config.Dns != "" {
		dns.SetResolver(config.Dns)
	}
	_, err := quic.DialAddrEarly(config.Addr, tls.ClientTLSConfig, QuicConfig)
	if err != nil {
		return nil, err
	}
	return &Client{
		config: config,
		cache:  map[string]string{},
	}, nil
}

func (c *Client) Mode() Mode {
	c.RLock()
	defer c.RUnlock()
	return c.config.Mode
}

func (c *Client) GetConfig() ClientConfig {
	c.RLock()
	defer c.RUnlock()
	return c.config
}

func (c *Client) MarshalCache() ([]byte, error) {
	c.RLock()
	defer c.RUnlock()
	return json.MarshalIndent(c.cache, "", "  ")
}

func (c *Client) GetFromCache(host string) (string, bool) {
	c.RLock()
	defer c.RUnlock()
	isocode, ok := c.cache[host]
	return isocode, ok
}

func (c *Client) SetConfig(config ClientConfig) {
	c.Lock()
	defer c.Unlock()
	c.config = config
}

func (c *Client) SetCache(host, isocode string) {
	c.Lock()
	defer c.Unlock()
	c.cache[host] = isocode
}

var (
	a = 0
	b = 0
)

func (c *Client) ConectRemote() (*Stream, error) {
	session, err := quic.DialAddrEarly(c.config.Addr, tls.ClientTLSConfig, QuicConfig)
	if err != nil {
		return nil, err
	}
	stream, err := session.OpenStream()
	if err != nil {
		return nil, err
	}
	go func() {
		time.Sleep(2 * time.Second)
		a++
		if session.ConnectionState().TLS.Used0RTT {
			b++
		}
		fmt.Println(a, b)
	}()
	return &Stream{stream}, nil
}
