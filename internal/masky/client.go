package masky

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/lucas-clemente/quic-go"
	"github.com/wkj9893/masky/internal/dns"
	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/tls"
)

type Client struct {
	config  ClientConfig
	cache   map[string]string // isocode cache
	session quic.Session
	mutex   sync.Mutex
}

func NewClient(config ClientConfig) (*Client, error) {
	if config.Dns != "" {
		dns.SetResolver(config.Dns)
	}
	c := &Client{
		config: config,
		cache:  make(map[string]string),
	}
	s, err := quic.DialAddr(c.config.Addr, tls.DefaultTLSConfig, DefaultQuicConfig)
	if err != nil {
		return nil, err
	}
	// auth
	err = s.SendMessage([]byte(c.config.Password))
	if err != nil {
		return nil, err
	}
	message, err := s.ReceiveMessage()
	if err != nil {
		return nil, err
	}
	if string(message) != "ok" {
		return nil, errors.New("fail to auth")
	}
	c.session = s
	return c, nil
}

func (c *Client) Mode() Mode {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.config.Mode
}

func (c *Client) GetConfig() ClientConfig {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.config
}

func (c *Client) SetConfig(config ClientConfig) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.config = config
}

func (c *Client) MarshalCache() ([]byte, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return json.MarshalIndent(c.cache, "", "  ")
}

func (c *Client) GetCache(host string) (string, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	isocode, ok := c.cache[host]
	return isocode, ok
}

func (c *Client) SetCache(host, isocode string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache[host] = isocode
}

func (c *Client) ConectRemote() (*Stream, error) {
	stream, err := c.session.OpenStream()
	if err == nil {
		return &Stream{stream}, nil
	}
	log.Warn(err)
	newClient, err := NewClient(c.config)
	if err != nil {
		return nil, err
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.session = newClient.session
	stream, err = c.session.OpenStream()
	if err != nil {
		return nil, err
	}
	log.Info("reconnect server successfully")
	return &Stream{stream}, nil
}
