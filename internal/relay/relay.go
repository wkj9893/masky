package relay

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
	"github.com/marten-seemann/webtransport-go"
	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"
	"github.com/wkj9893/masky/internal/tls"
	"golang.org/x/net/websocket"
)

var (
	tlsConf = tls.ClientTLSConfig()
)

func Run(config *Config) {
	log.SetLogLevel(config.LogLevel)
	switch config.Type {
	case "websocket":
		m := http.NewServeMux()
		m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello World! %s", time.Now())
		})
		for _, p := range config.Proxies {
			m.Handle("/"+p.ID.String(), websocket.Handler(func(c *websocket.Conn) {
				p, err := auth(p.ID, config)
				if err != nil {
					return
				}
				handleConn(masky.NewConn(c), p, config)
			}))
		}
		http.ListenAndServeTLS(fmt.Sprintf(":%v", config.Port), config.Cert, config.Key, m)
	case "webtransport":
		s := webtransport.Server{
			H3: http3.Server{Addr: fmt.Sprintf(":%v", config.Port)},
		}
		m := http.NewServeMux()
		m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello World! %s", time.Now())
		})
		for _, p := range config.Proxies {
			m.HandleFunc("/"+p.ID.String(), func(w http.ResponseWriter, r *http.Request) {
				c, err := s.Upgrade(w, r)
				if err != nil {
					return
				}
				stream, err := c.AcceptStream(context.Background())
				if err != nil {
					return
				}
				handleConn(masky.NewConn(stream), &p, config)
			})
		}
		s.ListenAndServeTLS(config.Cert, config.Key)
	default:
		l, err := quic.ListenAddrEarly(fmt.Sprintf(":%v", config.Port), tlsConf, nil)
		if err != nil {
			panic(err)
		}
		for {
			c, err := l.Accept(context.Background())
			if err != nil {
				panic(err)
			}
			stream, err := c.AcceptStream(context.Background())
			if err != nil {
				panic(err)
			}
			conn := masky.NewConn(stream)
			id, err := conn.Reader().Peek(16)
			if err != nil {
				return
			}
			var i uuid.UUID
			copy(i[:], id)
			p, err := auth(uuid.UUID(i), config)
			if err != nil {
				return
			}
			handleConn(conn, p, config)
		}
	}
}

func handleConn(c io.ReadWriteCloser, p *Proxy, config *Config) error {
	defer c.Close()
	dst, err := conectRemote(p)
	if err != nil {
		return err
	}
	masky.Relay(c, dst)
	return nil
}

func auth(id uuid.UUID, config *Config) (*Proxy, error) {
	for _, v := range config.Proxies {
		if v.ID == id {
			return &v, nil
		}
	}
	return nil, errors.New("cannot authorzize user, fail to find uuid")
}
