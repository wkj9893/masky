package test

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/wkj9893/masky/internal/client"
	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/server"
)

func TestMain(t *testing.T) {
	clientConf := &client.Config{
		Port:     3000,
		Mode:     client.GlobalMode,
		Addr:     "127.0.0.1:4000",
		AllowLan: false,
		LogLevel: log.InfoLevel,
	}
	serverConf := &server.Config{
		Port:     4000,
		LogLevel: log.InfoLevel,
	}
	go func() {
		client.Run(clientConf)
	}()
	go func() {
		server.Run(serverConf)
	}()

	time.Sleep(time.Millisecond)
	httpClient := http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			Proxy: func(r *http.Request) (*url.URL, error) {
				return &url.URL{
					Scheme: "http",
					Host:   fmt.Sprintf("localhost:%v", clientConf.Port),
				}, nil
			},
		},
	}
	socksClient := http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			Proxy: func(r *http.Request) (*url.URL, error) {
				return &url.URL{
					Scheme: "socks5",
					Host:   fmt.Sprintf("localhost:%v", clientConf.Port),
				}, nil
			},
		},
	}

	if _, err := httpClient.Get("http://example.com"); err != nil {
		t.Error(err)
	}
	if _, err := httpClient.Get("https://example.com"); err != nil {
		t.Error(err)
	}
	if _, err := socksClient.Get("http://example.com"); err != nil {
		t.Error(err)
	}
	if _, err := socksClient.Get("https://example.com"); err != nil {
		t.Error(err)
	}
}
