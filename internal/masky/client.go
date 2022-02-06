package masky

import (
	"github.com/lucas-clemente/quic-go"
	"github.com/wkj9893/masky/internal/dns"
	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/tls"
)

type Config struct {
	Port     string
	Mode     Mode
	Addr     string
	Dns      string
	LogLevel log.Level
}

type Client struct {
	config  Config
	session quic.Session
}

func NewClient(config Config) (*Client, error) {
	dns.SetResolver(config.Dns)
	c := &Client{
		config: config,
	}
	s, err := quic.DialAddr(c.config.Addr, tls.DefaultTLSConfig, DefaultQuicConfig)
	if err != nil {
		return nil, err
	}
	c.session = s
	return c, nil
}

func (c *Client) Config() Config {
	return c.config
}

func (c *Client) ConectRemote() (*Stream, error) {
	stream, err := c.session.OpenStream()
	if err != nil {
		return nil, err
	}
	return &Stream{stream}, nil
}
