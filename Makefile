BUILD_TAG := $(shell git describe --tags)

.PHONY: build
build:
	go build -ldflags "-X main.version=$(BUILD_TAG)" -o sqlpkg

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	go vet ./...
	golangci-lint run --print-issued-lines=false --out-format=colored-line-number ./...
