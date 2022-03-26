package geoip

import (
	"net"

	"github.com/oschwald/maxminddb-golang"
)

var name = "Country.mmdb"

func Lookup(ip net.IP) (string, error) {
	db, err := maxminddb.Open(name)
	if err != nil {
		return "", err
	}
	defer db.Close()

	var record struct {
		Country struct {
			ISOCode string `maxminddb:"iso_code"`
		} `maxminddb:"country"`
	}

	err = db.Lookup(ip, &record)
	if err != nil {
		return "", err
	}
	return record.Country.ISOCode, nil
}
