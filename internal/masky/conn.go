package masky

import (
	"bufio"
	"io"
	"net"
	"sync"
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
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		if _, err := io.Copy(left, right); err == nil {
			left.Close()
			right.Close()
		}
		wg.Done()
	}()
	go func() {
		if _, err := io.Copy(right, left); err == nil {
			left.Close()
			right.Close()
		}
		wg.Done()
	}()
	wg.Wait()
}

func lookup(host, port string) (string, error) {
	ip, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		return "", err
	}
	if isocode, err := geoip.Lookup(ip.IP); err == nil && isocode != "" {
		return isocode, nil
	}
	if _, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), 100*time.Millisecond); err != nil {
		return "", err
	}
	return "CN", nil
}

func Lookup(host, port string, c *Client) string {
	t := time.Now()
	if isocode, ok := c.GetCache(host); ok {
		log.Info(time.Since(t), host, isocode)
		return isocode
	}
	isocode, err := lookup(host, port)
	if err != nil {
		log.Warn(err)
	}
	c.SetCache(host, isocode)
	log.Info(time.Since(t), host, isocode)
	return isocode
}
