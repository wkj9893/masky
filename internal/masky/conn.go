package masky

import (
	"bufio"
	"io"
	"net"
	"time"
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
