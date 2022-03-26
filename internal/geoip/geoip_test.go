package geoip

import (
	"net"
	"testing"
)

func TestLookup(t *testing.T) {
	name = "../../Country.mmdb"
	tests := []string{"example.com", "google.com"}
	for _, test := range tests {
		if _, err := Lookup(net.IP(test)); err != nil {
			t.FailNow()
		}
	}
}
