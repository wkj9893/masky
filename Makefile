lint:
	golangci-lint run

GO_FLAGS = -trimpath -ldflags '-w -s'

build-client:
	CGO_ENABLED=0 go build $(GO_FLAGS) ./cmd/masky-client

build-server:
	CGO_ENABLED=0 go build $(GO_FLAGS) ./cmd/masky-server

docker:
	docker build --tag masky-server .

docker-run:
	docker run --rm -p 1080:1080/udp masky-server

.PHONY: lint build-client build-server docker docker-run

