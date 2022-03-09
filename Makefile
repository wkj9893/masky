lint:
	golangci-lint run

GO_FLAGS = -trimpath -ldflags '-w -s'

build-client:
	CGO_ENABLED=0 go build $(GO_FLAGS) ./cmd/masky-client

build-server:
	CGO_ENABLED=0 go build $(GO_FLAGS) ./cmd/masky-server

docker-build-client:
	docker build --tag masky-client -f ./client.Dockerfile .

docker-build-server:
	docker build --tag masky-server .

download:
	curl -o Country.mmdb https://raw.githubusercontent.com/P3TERX/GeoLite.mmdb/download/GeoLite2-Country.mmdb 

.PHONY: lint build-client build-server docker docker-run

