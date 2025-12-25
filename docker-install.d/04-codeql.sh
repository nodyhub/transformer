#!/bin/bash
set -euo pipefail

echo "Installing CodeQL dependencies..."
apt-get update
apt-get install -y --no-install-recommends ca-certificates wget unzip
rm -rf /var/lib/apt/lists/*

echo "Installing CodeQL..."
CODEQL_VERSION="2.15.5"

# Detect architecture
ARCH=$(uname -m)
if [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
    echo "Skipping CodeQL: ARM64 Linux binaries not available"
    echo "CodeQL requires x86_64 architecture"
    exit 0
fi

# Install for x86_64
wget -qO codeql.zip "https://github.com/github/codeql-cli-binaries/releases/download/v${CODEQL_VERSION}/codeql-linux64.zip"
unzip -q codeql.zip -d /opt
rm codeql.zip
ln -s /opt/codeql/codeql /usr/local/bin/codeql
codeql version
