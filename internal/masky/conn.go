package masky

import (
	"bufio"
	"io"
	"net"
	"time"

	"github.com/wkj9893/masky/internal/geoip"
	"github.com/wkj9893/masky/internal/log"
)

type Conn struct {
	r   *bufio.Reader
	rwc io.ReadWriteCloser
}

func NewConn(c io.ReadWriteCloser) *Conn {
	return &Conn{r: bufio.NewReader(c), rwc: c}
}

func (c *Conn) Reader() *bufio.Reader {
	return c.r
}

func (c *Conn) Read(b []byte) (n int, err error) {
	return c.r.Read(b)
}

func (c *Conn) Write(b []byte) (n int, err error) {
	return c.rwc.Write(b)
}

func (c *Conn) Close() error {
	return c.rwc.Close()
}

const (
	defaultDialTimeout = 5 * time.Second
)

func Dial(addr string) (net.Conn, error) {
	return net.DialTimeout("tcp", addr, defaultDialTimeout)
}

func Relay(left, right io.ReadWriteCloser) {
	ch := make(chan int)
	go func() {
		if _, err := io.Copy(left, right); err == nil {
			left.Close()
			right.Close()
		}
		ch <- 1
	}()
	if _, err := io.Copy(right, left); err == nil {
		left.Close()
		right.Close()
	}
	<-ch
}

func lookup(host, port string) (string, error) {
	t := time.Now()
	ip, err := net.LookupIP(host)
	if err != nil {
		return "", err
	}
	for _, i := range ip {
		if isocode, err := geoip.Lookup(i); err == nil && isocode != "" {
			log.Info(time.Since(t), host, isocode)
			return isocode, nil
		}
	}
	if _, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), 200*time.Millisecond); err == nil {
		return "CN", nil
	}
	return "", nil
}

func Lookup(host, port string, c *Client) (string, error) {
	// t := time.Now()
	if isocode, ok := c.GetFromCache(host); ok {
		// log.Info(time.Since(t), "get from cache:", host, isocode)
		return isocode, nil
	}
	isocode, err := lookup(host, port)
	if err != nil {
		// lookup host error: no such host
		return "", err
	}
	c.SetCache(host, isocode)
	// log.Info(time.Since(t), "lookup host:", host, isocode)
	return isocode, nil
}
