package socks

import (
	"errors"
	"fmt"
	"io"
	"net"
)

const (
	maxBufLen      int  = 256
	atypIPv4       byte = 1
	atypDomainName byte = 3
	atypIPv6       byte = 4
)

type Addr []byte // RFC 1928 section 5		ATYP + ADDR + PORT
var buf [256]byte

func (addr Addr) String() string {
	port := fmt.Sprint(256*int(addr[len(addr)-2]) + int(addr[len(addr)-1]))
	if addr[0] == atypIPv4 || addr[0] == atypIPv6 {
		return net.JoinHostPort(net.IP(addr[1:len(addr)-2]).String(), port)
	}
	return net.JoinHostPort(string(addr[2:len(addr)-2]), port)
}

func Handshake(rw io.ReadWriter, port int) (Addr, error) {
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
	// read ATYP DST.ADDR DST.PORT
	addr, err := ReadAddr(rw, buf)
	if err != nil {
		return nil, err
	}
	// write VER REP RSV ATYP BND.ADDR BND.PORT
	if _, err = rw.Write([]byte{5, 0, 0, 1, 127, 0, 0, 1, byte(port / 256), byte(port % 256)}); err != nil {
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
	case atypIPv4:
		if _, err := io.ReadFull(r, b[1:net.IPv4len+3]); err != nil {
			return nil, err
		}
		return b[:net.IPv4len+3], nil
	case atypDomainName:
		if _, err := io.ReadFull(r, b[1:2]); err != nil {
			return nil, err
		}
		if _, err := io.ReadFull(r, b[2:b[1]+4]); err != nil {
			return nil, err
		}
		return b[:b[1]+4], nil
	case atypIPv6:
		if _, err := io.ReadFull(r, b[1:net.IPv6len+3]); err != nil {
			return nil, err
		}
		return b[:net.IPv6len+3], nil
	}
	return nil, errors.New("atyp not supported")
}
