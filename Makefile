lint:
	golangci-lint run

Country.mmdb:
	curl -o Country.mmdb https://raw.githubusercontent.com/P3TERX/GeoLite.mmdb/download/GeoLite2-Country.mmdb 

GO_FLAGS = -trimpath -ldflags '-w -s'

build-client: Country.mmdb 
	CGO_ENABLED=0 go build $(GO_FLAGS) ./cmd/masky-client

build-server:
	CGO_ENABLED=0 go build $(GO_FLAGS) ./cmd/masky-server

dev: build-client
	cd web && pnpm build

.PHONY: lint build-client build-server dev

