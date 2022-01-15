lint:
	golangci-lint run

build:
	go build cmd/masky-client && go build cmd/masky-server

.PHONY: lint build
