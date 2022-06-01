package exit

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
	"github.com/wkj9893/masky/internal/socks"
	"github.com/wkj9893/masky/internal/tls"
	"golang.org/x/net/websocket"
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
				err := auth(p.ID, config)
				if err != nil {
					return
				}
				handleConn(masky.NewConn(c), config)
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
				handleConn(masky.NewConn(stream), config)
			})
		}
		s.ListenAndServeTLS(config.Cert, config.Key)
	default:
		tlsConf, err := tls.GenerateTLSConfig()
		if err != nil {
			panic(err)
		}
		l, err := quic.ListenAddrEarly(fmt.Sprintf(":%v", config.Port), tlsConf, nil)
		if err != nil {
			panic(err)
		}
		log.Info("start server successfully")
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
			var id uuid.UUID
			if _, err := io.ReadFull(conn, id[:]); err != nil {
				return
			}
			go handleConn(conn, config)
		}
	}
}

func handleConn(c *masky.Conn, config *Config) error {
	defer c.Close()
	head, err := c.Reader().Peek(1)
	if err != nil {
		return err
	}
	switch head[0] {
	case 5: // socks
		if _, err := c.Reader().ReadByte(); err != nil {
			return err
		}
		addr, err := socks.ReadAddr(c, make([]byte, 256))
		if err != nil {
			return err
		}
		dst, err := masky.Dial(addr.String())
		if err != nil {
			return err
		}
		masky.Relay(c, dst)
	default: // http
		req, err := http.ReadRequest(c.Reader())
		if err != nil {
			return err
		}
		if req.Method == http.MethodConnect {
			dst, err := masky.Dial(req.Host)
			if err != nil {
				return err
			}
			masky.Relay(c, dst)
		} else {
			req.RequestURI = ""
			resp, err := masky.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if err = resp.Write(c); err != nil {
				return err
			}
		}
	}
	return nil
}

func auth(id uuid.UUID, config *Config) error {
	for _, v := range config.Proxies {
		if v.ID == id {
			return nil
		}
	}
	return errors.New("cannot authorzize user, fail to find uuid")
}
