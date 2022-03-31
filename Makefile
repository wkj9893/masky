lint:
	golangci-lint run

Country.mmdb:
	curl -o Country.mmdb https://raw.githubusercontent.com/P3TERX/GeoLite.mmdb/download/GeoLite2-Country.mmdb

GOBUILD = CGO_ENABLED=0 go build -trimpath -ldflags '-w -s'

build-client: Country.mmdb
	$(GOBUILD) ./cmd/masky-client

build-server:
	$(GOBUILD) ./cmd/masky-server

.PHONY: lint build-client build-server 