package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"

	"github.com/lucas-clemente/quic-go"
	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"
	q "github.com/wkj9893/masky/internal/quic"
	"github.com/wkj9893/masky/internal/socks"
)

var (
	port int
	// logLevel int
)

func init() {
	flag.IntVar(&port, "port", 2022, "Listen Port")
	// flag.IntVar(&logLevel, "logLevel", 0, "Log Level")
	flag.Parse()
}

func main() {
	l, err := quic.ListenAddr(fmt.Sprintf(":%v", port), q.GenerateTLSConfig(), nil)
	if err != nil {
		log.Error(err)
	}
	for {
		s, err := l.Accept(context.Background())
		if err != nil {
			log.Error(err)
			return
		}
		go handleSession(s)
	}
}

func handleSession(s quic.Session) {
	stream, err := s.AcceptStream(context.Background())
	if err != nil {
		log.Error(err)
		return
	}
	c := masky.New(stream)

	head, err := c.Reader().Peek(1)
	if err != nil {
		log.Error(err)
		return
	}
	if head[0] == 5 { // socks
		if _, err = c.Reader().ReadByte(); err != nil {
			log.Error(err)
			return
		}
		addr, err := socks.ReadAddr(c, make([]byte, 256))
		if err != nil {
			log.Error(err)
			return
		}
		log.Info(fmt.Sprintf("%v -> %v -> %v", s.RemoteAddr(), s.LocalAddr(), addr))
		dst, err := net.Dial("tcp", addr.String())
		if err != nil {
			return
		}
		go masky.Copy(dst, c)
		go masky.Copy(c, dst)
	} else { // http
		req, err := http.ReadRequest(c.Reader())
		if err != nil {
			return
		}
		log.Info(fmt.Sprintf("%v -> %v -> %v", s.RemoteAddr(), s.LocalAddr(), req.URL.Hostname()))
		if req.Method == http.MethodConnect {
			dst, err := net.Dial("tcp", req.Host)
			if err != nil {
				log.Error(err)
				return
			}
			if _, err = fmt.Fprintf(c, "%v %v \r\n\r\n", req.Proto, http.StatusOK); err != nil {
				log.Error(err)
				return
			}
			go masky.Copy(dst, c)
			go masky.Copy(c, dst)
			return
		} else {
			defer c.Close()
			client := http.Client{}

			req.RequestURI = ""
			resp, err := client.Do(req)
			if err != nil {
				log.Error(err)
				return
			}
			if err = resp.Write(c); err != nil {
				log.Error(err)
			}
		}
	}
}
