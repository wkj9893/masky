lint:
	golangci-lint run

test:
	go test ./...

internal/geoip/Country.mmdb:
	curl -o internal/geoip/Country.mmdb https://raw.githubusercontent.com/P3TERX/GeoLite.mmdb/download/GeoLite2-Country.mmdb

GOBUILD = CGO_ENABLED=0 go build -trimpath -ldflags '-w -s'

build-client: internal/geoip/Country.mmdb
	$(GOBUILD) ./cmd/masky-client

build-server:
	$(GOBUILD) ./cmd/masky-server

.PHONY: lint test build-client build-server 
