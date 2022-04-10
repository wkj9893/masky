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
		Addr:     "127.0.0.1:3000",
		AllowLan: true,
		LogLevel: log.InfoLevel,
	}
	serverConf := &server.Config{
		Port:     3000,
		LogLevel: log.InfoLevel,
	}
	go func() {
		client.Run(clientConf)
	}()
	go func() {
		server.Run(serverConf)
	}()

	time.Sleep(time.Millisecond)
	os.Setenv("http_proxy", fmt.Sprintf("http://127.0.0.1:%v", clientConf.Port))
	if _, err := http.Get("http://example.com"); err != nil {
		t.Error(err)
	}
}
