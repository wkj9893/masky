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
	mu      sync.RWMutex
}

func NewClient(config ClientConfig) (*Client, error) {
	if config.Dns != "" {
		dns.SetResolver(config.Dns)
	}
	c := &Client{
		config: config,
		cache:  make(map[string]string),
	}
	s, err := quic.DialAddr(c.config.Addr, tls.DefaultTLSConfig, ClientQuicConfig)
	if err != nil {
		return nil, err
	}
	c.session = s
	if err := c.auth(); err != nil {
		_ = c.session.CloseWithError(DefaultApplicationErrorCode, err.Error())
		return nil, err
	}
	return c, nil
}

func (c *Client) auth() error {
	if err := c.session.SendMessage([]byte(c.config.Password)); err != nil {
		return err
	}
	if message, err := c.session.ReceiveMessage(); err != nil {
		return err
	} else if string(message) != "ok" {
		return errors.New("fail to auth")
	}
	return nil
}

func (c *Client) Mode() Mode {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config.Mode
}

func (c *Client) GetConfig() ClientConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config
}

func (c *Client) MarshalCache() ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return json.MarshalIndent(c.cache, "", "  ")
}

func (c *Client) GetFromCache(host string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	isocode, ok := c.cache[host]
	return isocode, ok
}

func (c *Client) SetConfig(config ClientConfig) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.config = config
}

func (c *Client) SetCache(host, isocode string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[host] = isocode
}

func (c *Client) ConectRemote() (*Stream, error) {
	if stream, err := c.session.OpenStream(); err == nil {
		return &Stream{stream}, nil
	} else {
		log.Warn(err, "try to reconnect server")
	}
	// try to reconnect server
	if err := c.reconnect(); err != nil {
		return nil, err
	}
	if stream, err := c.session.OpenStream(); err == nil {
		return &Stream{stream}, nil
	} else {
		return nil, err
	}
}

func (c *Client) reconnect() error {
	newClient, err := NewClient(c.config)
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.session = newClient.session
	log.Info("reconnect server successfully")
	return nil
}
