.PHONY: all build release install test clean

# Default target
all: build


VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT  ?= $(shell git rev-parse HEAD)
BUILDTIME ?= $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS = -ldflags "-X 'github.com/nodyhub/transformer/cmd.Version=$(VERSION)' -X 'github.com/nodyhub/transformer/cmd.Commit=$(COMMIT)' -X 'github.com/nodyhub/transformer/cmd.BuildTime=$(BUILDTIME)'"


# Build the binary
build:
	go build $(LDFLAGS) -o bin/transformer ./cmd/transformer && chmod +x bin/transformer

# Build binaries for multiple architectures
release:
	@echo "Building release binaries..."
	@mkdir -p bin/release
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/release/transformer-linux-amd64 ./cmd/transformer
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/release/transformer-linux-arm64 ./cmd/transformer
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/release/transformer-darwin-amd64 ./cmd/transformer
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/release/transformer-darwin-arm64 ./cmd/transformer
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/release/transformer-windows-amd64.exe ./cmd/transformer
	@echo "Release binaries built in bin/release/"

# Install the binary to $GOPATH/bin
install:
	go install $(LDFLAGS) ./cmd/transformer

# Run tests
test:
	go test ./...

# Run linter
lint:
	golangci-lint run

# Build Docker image
# Usage:
#   make image                                    # builds with MODULES=all, tag=latest
#   make image MODULES="git,ssh" TAG="v1.0"      # builds with custom MODULES and TAG
image: MODULES ?= all
image: TAG ?= latest
image:
	docker build --build-arg MODULES="$(MODULES)" -t transformer:$(TAG) .
	docker tag transformer:$(TAG) ghcr.io/nodyhub/transformer:$(TAG)



# Clean build artifacts
clean:
	rm -rf bin/
