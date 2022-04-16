package test

import (
	"fmt"
	"net/http"
	"os"
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

	httpProxy := fmt.Sprintf("http://127.0.0.1:%v", clientConf.Port)
	socksProxy := fmt.Sprintf("socks5://127.0.0.1:%v", clientConf.Port)
	os.Setenv("http_proxy", httpProxy)
	if err := get("http://example.com"); err != nil {
		t.Error(err)
	}

	os.Setenv("https_proxy", httpProxy)
	if err := get("https://example.com"); err != nil {
		t.Error(err)
	}

	os.Setenv("http_proxy", socksProxy)
	if err := get("http://example.com"); err != nil {
		t.Error(err)
	}

	os.Setenv("https_proxy", socksProxy)
	if err := get("https://example.com"); err != nil {
		t.Error(err)
	}
}

func get(url string) error {
	c := http.Client{
		Timeout: 3 * time.Second,
	}
	if _, err := c.Get(url); err != nil {
		return err
	}
	return nil
}
