package socks

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/wkj9893/masky/internal/geoip"
	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"
)

const (
	Version        = 5
	maxBufLen      = 256
	AtypIPv4       = 1
	AtypDomainName = 3
	AtypIPv6       = 4
)

type Addr []byte // RFC 1928 section 5		ATYP + ADDR + PORT

func (addr Addr) String() string {
	port := fmt.Sprint(256*int(addr[len(addr)-2]) + int(addr[len(addr)-1]))
	if addr[0] == AtypIPv4 || addr[0] == AtypIPv6 {
		return net.JoinHostPort(net.IP(addr[1:len(addr)-2]).String(), port)
	}
	return string(addr[2:len(addr)-2]) + ":" + port
}

func HandleConn(c *masky.Conn, client *masky.Client) {
	defer c.Close()
	addr, err := Handshake(c)
	if err != nil {
		log.Error(err)
		return
	}
	var dst io.ReadWriteCloser

	switch client.Config().Mode {
	case masky.DirectMode:
		if dst, err = masky.Dial(addr.String()); err != nil {
			log.Warn(err)
			return
		}
	case masky.GlobalMode:
		if dst, err = client.ConectRemote(); err != nil {
			log.Error(err)
			return
		}
		if _, err := dst.Write(append([]byte{5}, addr...)); err != nil {
			log.Error(err)
			return
		}
	case masky.RuleMode:
		isocode, err := lookup(addr)
		if err != nil {
			log.Warn(err)
			isocode = "CN"
		}
		if isocode == "CN" {
			if dst, err = masky.Dial(addr.String()); err != nil {
				log.Warn(err)
				return
			}
		} else {
			if dst, err = client.ConectRemote(); err != nil {
				log.Error(err)
				return
			}
			if _, err = dst.Write(append([]byte{5}, addr...)); err != nil {
				log.Error(err)
				return
			}
		}
	}
	defer dst.Close()
	masky.Relay(c, dst)
}

func Handshake(rw io.ReadWriter) (Addr, error) {
	buf := make([]byte, maxBufLen)
	// read VER, NMETHODS
	if _, err := io.ReadFull(rw, buf[:2]); err != nil {
		return nil, err
	}
	// read METHODS
	if _, err := io.ReadFull(rw, buf[:buf[1]]); err != nil {
		return nil, err
	}
	// write VER METHOD
	if _, err := rw.Write([]byte{5, 0}); err != nil {
		return nil, err
	}
	// read VER, CMD, RSV
	if _, err := io.ReadFull(rw, buf[:3]); err != nil {
		return nil, err
	}
	if buf[1] != 1 {
		return nil, errors.New("cmd not supported")
	}
	// read DST.ADDR DST.PORT
	addr, err := ReadAddr(rw, buf)
	if err != nil {
		return nil, err
	}
	if _, err = rw.Write([]byte{5, 0, 0, 1, 127, 0, 0, 1, 7, 229}); err != nil {
		return nil, err
	}
	return addr, nil
}

func ReadAddr(r io.Reader, b []byte) (Addr, error) {
	// read ATYP
	if _, err := io.ReadFull(r, b[:1]); err != nil {
		return nil, err
	}
	switch b[0] {
	case AtypIPv4:
		if _, err := io.ReadFull(r, b[1:net.IPv4len+3]); err != nil {
			return nil, err
		}
		return b[:net.IPv4len+3], nil
	case AtypDomainName:
		if _, err := io.ReadFull(r, b[1:2]); err != nil {
			return nil, err
		}
		if _, err := io.ReadFull(r, b[2:b[1]+4]); err != nil {
			return nil, err
		}
		return b[:b[1]+4], nil
	case AtypIPv6:
		if _, err := io.ReadFull(r, b[1:net.IPv6len+2]); err != nil {
			return nil, err
		}
		return b[:net.IPv6len+3], nil
	}
	return nil, errors.New("atyp not supported")
}

func lookup(addr Addr) (string, error) {
	if addr[0] == AtypDomainName {
		ip, err := net.ResolveIPAddr("ip", string(addr[2:len(addr)-2]))
		if err != nil {
			return "", err
		}
		return geoip.Lookup(ip.IP)
	}
	return geoip.Lookup(net.IP(addr[1 : len(addr)-2]))
}
