BUILD_TAG := $(shell git describe --tags)

.PHONY: build
build:
	go build -ldflags "-X main.version=$(BUILD_TAG)" -o sqlpkg

.PHONY: test
test:
	go test -v ./...
