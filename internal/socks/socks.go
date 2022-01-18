package socks

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/wkj9893/masky/internal/geoip"
	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"
	"github.com/wkj9893/masky/internal/mode"
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

func HandleConn(c *masky.Conn, config masky.Config) {
	addr, err := Handshake(c)
	if err != nil {
		return
	}
	var dst io.ReadWriteCloser

	if config.Mode == mode.Direct {
		if dst, err = masky.Dial("tcp", addr.String()); err != nil {
			log.Error(err)
			return
		}
	} else if config.Mode == mode.Global {
		if dst, err = masky.ConectRemote(config.Addr); err != nil {
			log.Error(err)
			return
		}
		if _, err := dst.Write(append([]byte{5}, addr...)); err != nil {
			log.Error(err)
			return
		}
	} else if config.Mode == mode.Rule {
		isocode, err := lookup(addr)
		if err != nil {
			log.Error(err)
			return
		}
		if isocode == "CN" {
			if dst, err = masky.Dial("tcp", addr.String()); err != nil {
				return
			}
		} else {
			if dst, err = masky.ConectRemote(config.Addr); err != nil {
				return
			}
			if _, err = dst.Write(append([]byte{5}, addr...)); err != nil {
				return
			}
		}
	} else {
		log.Error("Unknown Mode")
		return
	}
	go masky.Copy(dst, c)
	go masky.Copy(c, dst)
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
		ip, err := net.LookupIP(string(addr[2 : len(addr)-2]))
		if err != nil {
			return "", err
		}
		return geoip.Lookup(ip[0])
	}
	return geoip.Lookup(net.IP(addr[1 : len(addr)-2]))
}