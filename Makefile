lint:
	golangci-lint run

Country.mmdb:
	curl -o Country.mmdb https://raw.githubusercontent.com/P3TERX/GeoLite.mmdb/download/GeoLite2-Country.mmdb 

build-web:
	cd web && pnpm build

build-client: Country.mmdb build-web
	CGO_ENABLED=0 go build $(GO_FLAGS) ./cmd/masky-client

build-server:
	CGO_ENABLED=0 go build $(GO_FLAGS) ./cmd/masky-server

GO_FLAGS = -trimpath -ldflags '-w -s'

.PHONY: lint build-client build-server build-web 

