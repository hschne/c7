BINARY := c7
VERSION := 0.1.0
LDFLAGS := -ldflags "-X main.version=$(VERSION) -s -w"

.PHONY: build install test clean

build:
	go build $(LDFLAGS) -o $(BINARY) .

# Build for common platforms
build-all:
	GOOS=linux   GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-amd64 .
	GOOS=darwin  GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-amd64 .
	GOOS=darwin  GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-windows-amd64.exe .

install: build
	cp $(BINARY) /usr/local/bin/$(BINARY)

test:
	go test -v ./...

clean:
	rm -f $(BINARY)
	rm -rf dist/
