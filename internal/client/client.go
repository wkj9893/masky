package client

import (
	"fmt"
	"net"
	"os"

	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"
	"github.com/wkj9893/masky/internal/tls"
	"gopkg.in/yaml.v3"
)

var tlsConf = tls.ClientTLSConfig()

func Run(config *Config) {
	log.SetLogLevel(config.LogLevel)
	setConfig(config)
	addr := fmt.Sprintf(":%v", config.Port)
	if !config.AllowLan {
		addr = fmt.Sprintf("127.0.0.1:%v", config.Port)
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	log.Info("client listen on port", config.Port)
	go func() {
		StartApi()
	}()
	for {
		if c, err := l.Accept(); err == nil {
			go handleConn(c, config)
		} else {
			panic(err)
		}
	}
}

func handleConn(c net.Conn, config *Config) {
	defer c.Close()
	conn := masky.NewConn(c)
	head, err := conn.Reader().Peek(1)
	if err != nil {
		log.Error(err)
		return
	}
	if head[0] == 5 {
		if err := handleSocks(conn, config); err != nil {
			log.Warn("client error:", err)
		}
	} else {
		if err := handleHttp(conn, config); err != nil {
			log.Warn("client error:", err)
		}
	}
}

func ParseConfig(name string) (*Config, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	if err := yaml.Unmarshal(data, c); err != nil {
		return nil, err
	}
	return c, nil
}
