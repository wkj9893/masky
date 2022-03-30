package masky

// import (
// 	"encoding/json"
// 	"sync"

// 	"github.com/lucas-clemente/quic-go"
// 	"github.com/wkj9893/masky/internal/client"
// 	"github.com/wkj9893/masky/internal/tls"
// )

// type Client struct {
// 	sync.RWMutex

// 	config client.Config
// 	cache  map[string]string // isocode cache
// 	m      map[quic.EarlySession]quic.Stream
// }

// func NewClient(config client.Config) (*Client, error) {
// 	_, err := quic.DialAddrEarly(config.Addr, tls.ClientTLSConfig, QuicConfig)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Client{
// 		config: config,
// 		cache:  map[string]string{},
// 		m:      map[quic.EarlySession]quic.Stream{},
// 	}, nil
// }

// func (c *Client) GetConfig() client.Config {
// 	c.RLock()
// 	defer c.RUnlock()
// 	return c.config
// }

// func (c *Client) SetConfig(config client.Config) {
// 	c.Lock()
// 	defer c.Unlock()
// 	c.config = config
// }

// func (c *Client) GetFromCache(host string) (string, bool) {
// 	c.RLock()
// 	defer c.RUnlock()
// 	isocode, ok := c.cache[host]
// 	return isocode, ok
// }

// func (c *Client) SetCache(host, isocode string) {
// 	c.Lock()
// 	defer c.Unlock()
// 	c.cache[host] = isocode
// }

// func (c *Client) MarshalCache() ([]byte, error) {
// 	c.RLock()
// 	defer c.RUnlock()
// 	return json.MarshalIndent(c.cache, "", "  ")
// }

// func (c *Client) GetConn() ([]byte, error) {
// 	c.RLock()
// 	defer c.RUnlock()
// 	return json.MarshalIndent(c.cache, "", "  ")
// }

// func (c *Client) find() quic.Stream {
// 	c.Lock()
// 	defer c.Unlock()
// 	for session, stream := range c.m {
// 		if !isActive(stream) {
// 			stream, err := session.OpenStream()
// 			if err != nil {
// 				delete(c.m, session)
// 				continue
// 			}
// 			c.m[session] = stream
// 			return stream
// 		}
// 	}
// 	return nil
// }

// func (c *Client) ConectRemote() (*Stream, error) {
// 	if stream := c.find(); stream != nil {
// 		return &Stream{stream}, nil
// 	}
// 	session, err := quic.DialAddrEarly(c.config.Addr, tls.ClientTLSConfig, QuicConfig)
// 	if err != nil {
// 		return nil, err
// 	}
// 	stream, err := session.OpenStream()
// 	if err != nil {
// 		return nil, err
// 	}
// 	c.Lock()
// 	c.m[session] = stream
// 	c.Unlock()
// 	return &Stream{stream}, nil
// }

// func isActive(s quic.Stream) bool {
// 	select {
// 	case <-s.Context().Done():
// 		return false
// 	default:
// 		return true
// 	}
// }
