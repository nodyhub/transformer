# Multi-stage build for transformer with binary and dependencies

# Stage 1: Build the Go binary
FROM golang:1.25.2-bookworm AS builder

# install make and git
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    make git && rm -rf /var/lib/apt/lists/*

WORKDIR /build

# Set GO_PRIVATE if needed (uncomment and set your private repo domain)
ENV GO_PRIVATE=github.com/nodyhub/transformer/*

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN make build

# Stage 2: Transformer runtime image
FROM ubuntu:24.04

# Copy the built binary from the builder stage
COPY --from=builder /build/bin/transformer /usr/local/bin/transformer

# Verify transformer installation
RUN echo "Transformer version:" && transformer version

# Build argument to specify which modules to install, default is empty and installs none, 
# use "all" to install all available modules, comma separated list or 
# run 'transformer modules list' to see available modules
ARG MODULES="builtin"

# Install specified modules
WORKDIR /app
RUN transformer modules install "${MODULES}"

# Create workspace directory
WORKDIR /workspace

# Set transformer as default command
CMD ["transformer", "--help"]
