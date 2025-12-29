# Multi-stage build for transformer with security scanning tools

# Stage 1: Security tools base
FROM ubuntu:24.04 AS sec-tools

# Build argument to specify which modules to install
ARG MODULES=""

# Install dependencies and security tools
WORKDIR /app
COPY docker-install.d /tmp/install.d
RUN chmod +x /tmp/install.d/*.sh && \
    bash /tmp/install.d/install-modules.sh "$MODULES" && \
    rm -rf /tmp/install.d

# Optionally install golang for tools that need it at runtime
RUN apt-get update && apt-get install -y --no-install-recommends golang && rm -rf /var/lib/apt/lists/*

# Create workspace directory
RUN mkdir -p /workspace

# Stage 2: Build the Go binary
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

# Stage 3: Final image
FROM sec-tools
WORKDIR /workspace
COPY --from=builder /build/bin/transformer /usr/local/bin/transformer

# Verify transformer installation
RUN transformer --version || echo "Transformer installed"

# Set transformer as default command
CMD ["transformer", "--help"]
