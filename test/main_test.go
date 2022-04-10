package test

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/wkj9893/masky/internal/client"
	"github.com/wkj9893/masky/internal/relay"
	"github.com/wkj9893/masky/internal/server"
)

func TestMain(t *testing.T) {
	clientConfig, err := client.ParseConfig("client.yaml")
	if err != nil {
		t.Error(err)
	}

	relayConfig, err := relay.ParseConfig("relay.yaml")
	if err != nil {
		t.Error(err)
	}

	serverConfig, err := server.ParseConfig("server.yaml")
	if err != nil {
		t.Error(err)
	}

	go func() {
		client.Run(clientConfig)
	}()

	go func() {
		relay.Run(relayConfig)
	}()

	go func() {
		server.Run(serverConfig)
	}()

	time.Sleep(time.Millisecond)
	os.Setenv("http_proxy", fmt.Sprintf("http://127.0.0.1:%v", clientConfig.Port))
	if _, err := http.Get("http://example.com"); err != nil {
		t.Error(err)
	}
}
