package geoip

import (
	"net"

	"github.com/oschwald/maxminddb-golang"
)

func Lookup(ip net.IP) (string, error) {
	db, err := maxminddb.Open("./GeoLite2-Country.mmdb") // https://github.com/wkj9893/geoip/releases/latest/download/GeoLite2-Country.mmdb
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
