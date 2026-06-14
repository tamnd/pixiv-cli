BINARY  := bin/pixiv
PKG     := ./cmd/pixiv
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w \
	-X github.com/tamnd/pixiv-cli/cli.Version=$(VERSION) \
	-X github.com/tamnd/pixiv-cli/cli.Commit=$(COMMIT) \
	-X github.com/tamnd/pixiv-cli/cli.Date=$(DATE)

.PHONY: build install test vet fmt clean run

build:
	@mkdir -p $(dir $(BINARY))
	CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(BINARY) $(PKG)

install:
	CGO_ENABLED=0 go install -trimpath -ldflags "$(LDFLAGS)" $(PKG)

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w -s .

clean:
	rm -rf bin dist

run: build
	./$(BINARY) $(ARGS)
