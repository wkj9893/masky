package masky

import (
	"net"
	"sync"
	"time"

	"github.com/wkj9893/masky/internal/geoip"
)

var (
	mu sync.RWMutex

	m = map[string]string{}
)

func lookup(host, port string) (string, error) {
	ip, err := net.LookupIP(host)
	if err != nil {
		return "", err
	}
	for _, i := range ip {
		if isocode, err := geoip.Lookup(i); err == nil && isocode != "" {
			return isocode, nil
		}
	}
	if _, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), 200*time.Millisecond); err == nil {
		return "CN", nil
	}
	return "", nil
}

func Lookup(host, port string) (string, error) {
	if isocode, ok := get(host); ok {
		return isocode, nil
	}
	isocode, err := lookup(host, port)
	if err != nil {
		return "", err
	}
	set(host, isocode)
	return isocode, nil
}

func get(host string) (string, bool) {
	mu.RLock()
	isocode, ok := m[host]
	mu.RUnlock()
	return isocode, ok
}

func set(host, isocode string) {
	mu.Lock()
	m[host] = isocode
	mu.Unlock()
}
