package test

import (
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/wkj9893/masky/internal/client"
	"github.com/wkj9893/masky/internal/relay"
	"github.com/wkj9893/masky/internal/server"
)

var (
	id1 = uuid.MustParse("00000000-0000-0000-0000-000000000000")
	id2 = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	id3 = uuid.MustParse("22222222-2222-2222-2222-222222222222")

	httpClient = http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				Host:   "localhost:3000",
			}),
		},
	}

	socksClient = http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "socks5",
				Host:   "localhost:3000",
			}),
		},
	}
)

// test client -> server
// func TestClientServer(t *testing.T) {
// 	clientConfig := &client.Config{
// 		Port: 3000,
// 		Mode: client.GlobalMode,
// 		Proxies: []client.Proxy{
// 			{ID: id1},
// 			{ID: id2, Server: []string{"127.0.0.1:4000"}},
// 		},
// 	}
// 	serverConfig := &server.Config{
// 		Port: 4000,
// 		Proxies: []server.Proxy{
// 			{ID: id2},
// 		},
// 	}
// 	go func() {
// 		client.Run(clientConfig)
// 	}()
// 	go func() {
// 		server.Run(serverConfig)
// 	}()

// 	time.Sleep(time.Millisecond)
// 	get(t)
// }

// test client -> relay -> server
func TestClientRelayServer(t *testing.T) {
	clientConfig := &client.Config{
		Port: 3000,
		Mode: client.GlobalMode,
		Proxies: []client.Proxy{
			{ID: id1},
			{ID: id2, Server: []string{"127.0.0.1:4000", "127.0.0.1:5000"}},
		},
	}
	relayConfig := &relay.Config{
		Port: 4000,
		Proxies: []relay.Proxy{
			{ID: id2, Server: "127.0.0.1:5000"},
		},
	}
	serverConfig := &server.Config{
		Port: 5000,
		Proxies: []server.Proxy{
			{ID: id2},
		},
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
	get(t)
}

func get(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		if _, err := httpClient.Get("http://example.com"); err != nil {
			t.Error(err)
		}
		wg.Done()
	}()
	go func() {
		if _, err := httpClient.Get("https://example.com"); err != nil {
			t.Error(err)
		}
		wg.Done()
	}()
	go func() {
		if _, err := socksClient.Get("http://example.com"); err != nil {
			t.Error(err)
		}
		wg.Done()
	}()
	go func() {
		if _, err := socksClient.Get("https://example.com"); err != nil {
			t.Error(err)
		}
		wg.Done()
	}()
	wg.Wait()
}
