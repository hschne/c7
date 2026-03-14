BINARY := c7
VERSION := dev
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE) -s -w"

.PHONY: build install test lint clean

build:
	go build $(LDFLAGS) -o $(BINARY) .

install: build
	cp $(BINARY) $(HOME)/.local/bin/$(BINARY)

test:
	go test -v ./...

lint:
	golangci-lint run ./...

clean:
	rm -f $(BINARY)
	rm -rf dist/
