package masky

import (
	"bufio"
	"crypto/tls"
	"io"

	"github.com/lucas-clemente/quic-go"
)

type Conn struct {
	r   *bufio.Reader
	rwc io.ReadWriteCloser
}

func New(c io.ReadWriteCloser) *Conn {
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

func ConectRemote(addr string) (quic.Stream, error) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"masky"},
	}
	session, err := quic.DialAddr(addr, tlsConf, nil)
	if err != nil {
		return nil, err
	}
	stream, err := session.OpenStream()
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func Copy(dst io.WriteCloser, src io.ReadCloser) {
	if _, err := io.Copy(dst, src); err != nil {
		src.Close()
		dst.Close()
	}
}
