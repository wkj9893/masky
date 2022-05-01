package geoip

import (
	"net"

	_ "embed"

	"github.com/oschwald/maxminddb-golang"
)

//go:embed Country.mmdb
var b []byte

func Lookup(ip net.IP) (string, error) {
	db, err := maxminddb.FromBytes(b)
	if err != nil {
		return "", err
	}
	defer db.Close()

	var record struct {
		Country struct {
			ISOCode string `maxminddb:"iso_code"`
		} `maxminddb:"country"`
	}
	if err := db.Lookup(ip, &record); err != nil {
		return "", err
	}
	return record.Country.ISOCode, nil
}
