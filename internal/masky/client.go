package masky

import (
	"errors"

	"github.com/lucas-clemente/quic-go"
	"github.com/wkj9893/masky/internal/dns"
	"github.com/wkj9893/masky/internal/tls"
)

type Client struct {
	config  ClientConfig
	session quic.Session
}

func NewClient(config ClientConfig) (*Client, error) {
	if config.Dns != "" {
		dns.SetResolver(config.Dns)
	}
	c := &Client{
		config: config,
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

func (c *Client) Config() ClientConfig {
	return c.config
}

func (c *Client) ConectRemote() (*Stream, error) {
	stream, err := c.session.OpenStream()
	if err != nil {
		return nil, err
	}
	return &Stream{stream}, nil
}
