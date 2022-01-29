package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/lucas-clemente/quic-go"
	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"
	"github.com/wkj9893/masky/internal/socks"
	"github.com/wkj9893/masky/internal/tls"
)

var port string
var s quic.Session

func init() {
	port = "2022"
	parseArgs(os.Args[1:])
}

func main() {
	tlsConf, err := tls.GenerateTLSConfig()
	if err != nil {
		panic(err)
	}
	l, err := quic.ListenAddr(":"+port, tlsConf, masky.DefaultQuicConfig)
	if err != nil {
		panic(err)
	}
	s, err = l.Accept(context.Background())
	if err != nil {
		panic(err)
	}
	for {
		stream, err := s.AcceptStream(context.Background())
		if err != nil {
			panic(err)
		}
		go handleStream(stream)
	}
}

func handleStream(stream quic.Stream) {
	c := masky.NewConn(stream)

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
		dst, err := masky.Dial(addr.String())
		if err != nil {
			log.Warn(err)
			return
		}
		go masky.Copy(c, dst)
		go masky.Copy(dst, c)
	} else { // http
		req, err := http.ReadRequest(c.Reader())
		if err != nil {
			log.Error(err)
			return
		}
		log.Info(fmt.Sprintf("%v -> %v -> %v", s.RemoteAddr(), s.LocalAddr(), req.URL.Hostname()))
		if req.Method == http.MethodConnect {
			dst, err := masky.Dial(req.Host)
			if err != nil {
				log.Warn(err)
				return
			}
			go masky.Copy(c, dst)
			go masky.Copy(dst, c)
		} else {
			defer c.Close()
			client := http.Client{Transport: &http.Transport{
				Proxy: nil, // unset proxy in case recursion
			}}
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

func parseArgs(args []string) {
	for _, arg := range args {
		switch {
		case arg == "-h", arg == "--help":
			fmt.Println("masky server")
		case strings.HasPrefix(arg, "--port="):
			port = arg[len("--port="):]
		}
	}
}
