package masky

import (
	"bufio"
	"fmt"
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

func Lookup(host string, port string) string {
	t := time.Now()
	ip, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		log.Warn(err)
		return ""
	}
	if isocode, err := geoip.Lookup(ip.IP); err == nil {
		fmt.Println(time.Since(t), host, ip, isocode)
		return isocode
	}
	if _, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), 100*time.Millisecond); err == nil {
		return "CN"
	} else {
		log.Warn(err)
	}
	return ""
}
