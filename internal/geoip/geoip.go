package geoip

import (
	"io"
	"net"
	"net/http"
	"os"

	"github.com/oschwald/maxminddb-golang"
	"github.com/wkj9893/masky/internal/log"
)

const name = "Country.mmdb"

func download() (err error) {
	resp, err := http.Get("https://cdn.jsdelivr.net/gh/P3TERX/GeoLite.mmdb@download/GeoLite2-Country.mmdb")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

func Lookup(ip net.IP) (string, error) {
	if _, err := os.Open(name); err != nil {
		if err := download(); err != nil {
			log.Warn("fail to download" + name)
			return "", err
		}
	}
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
