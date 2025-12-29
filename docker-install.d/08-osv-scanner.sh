#!/bin/bash
set -euo pipefail

echo "Installing OSV-Scanner dependencies..."
apt-get update
apt-get install -y --no-install-recommends ca-certificates wget curl
rm -rf /var/lib/apt/lists/*

echo "Installing OSV-Scanner..."
OSV_VERSION="2.3.1"
# Detect architecture
ARCH=$(uname -m)
if [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
    OSV_ARCH="arm64"
else
    OSV_ARCH="amd64"
fi
wget -qO /usr/local/bin/osv-scanner "https://github.com/google/osv-scanner/releases/download/v${OSV_VERSION}/osv-scanner_linux_${OSV_ARCH}"
chmod +x /usr/local/bin/osv-scanner
osv-scanner --version
