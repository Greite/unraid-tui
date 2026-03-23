BINARY_NAME=unraid-tui
BUILD_DIR=bin
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-s -w -X github.com/Greite/unraid-tui/cmd.version=$(VERSION) -X github.com/Greite/unraid-tui/cmd.commit=$(COMMIT) -X github.com/Greite/unraid-tui/cmd.date=$(DATE)

.PHONY: build test lint fmt run clean install uninstall release-dry

build:
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) .

test:
	go test ./...

test-verbose:
	go test -v ./...

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	go vet ./...

fmt:
	go run golang.org/x/tools/cmd/goimports@latest -w .
	gofmt -w .

run: build
	$(BUILD_DIR)/$(BINARY_NAME)

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

uninstall:
	rm -f /usr/local/bin/$(BINARY_NAME)

release-dry:
	goreleaser release --snapshot --clean

clean:
	rm -rf $(BUILD_DIR) dist coverage.out coverage.html
