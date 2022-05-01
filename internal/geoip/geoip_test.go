package geoip

import (
	"net"
	"testing"
)

func TestLookup(t *testing.T) {
	ip, err := net.LookupIP("example.com")
	if err != nil {
		t.Error(err)
	}
	for _, i := range ip {
		if _, err := Lookup(i); err != nil {
			t.Error(err)
		}
	}
}
