BINARY_NAME=zenfocus
VERSION=1.0.0
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

.PHONY: all build test clean install

all: build

build:
	@echo "Building $(BINARY_NAME) v$(VERSION)..."
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) .

test:
	@echo "Running tests..."
	go test -v -race ./...

install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	install -m 755 bin/$(BINARY_NAME) /usr/local/bin/

clean:
	@echo "Cleaning up..."
	rm -rf bin/
