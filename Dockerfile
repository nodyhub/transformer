# Multi-stage build for transformer with security scanning tools

# Stage 1: Build the Go binary
FROM golang:1.25.1-bookworm AS builder

WORKDIR /build

# Install build dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    make \
    && rm -rf /var/lib/apt/lists/*

# Set GO_PRIVATE if needed (uncomment and set your private repo domain)
ENV GO_PRIVATE=github.com/nodyhub/transformer/*

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN make build

# Stage 2: Runtime image with security tools
FROM ubuntu:24.04

# Build argument to specify which modules to install
# Comma-separated list of modules (e.g., "trivy,semgrep,nmap")
# If empty or "all", installs all modules
# Available modules: trivy, semgrep, nuclei, codeql, trufflehog, nmap, gosec, osv-scanner
ARG MODULES=""

WORKDIR /app

# Copy transformer binary from builder
COPY --from=builder /build/bin/transformer /usr/local/bin/transformer

# Copy installation scripts
COPY docker-install.d /tmp/install.d

# Run installation scripts based on MODULES argument
RUN chmod +x /tmp/install.d/*.sh && \
    bash /tmp/install.d/install-modules.sh "$MODULES" && \
    rm -rf /tmp/install.d

# Create workspace directory
RUN mkdir -p /workspace

WORKDIR /workspace

# Verify transformer installation
RUN transformer --version || echo "Transformer installed"

# Set transformer as default command
CMD ["transformer", "--help"]
