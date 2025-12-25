#!/bin/bash
set -euo pipefail

echo "Installing Gosec dependencies..."
apt-get update
apt-get install -y --no-install-recommends ca-certificates wget tar
rm -rf /var/lib/apt/lists/*

echo "Installing Gosec..."
GOSEC_VERSION="2.18.2"
# Detect architecture
ARCH=$(uname -m)
if [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
    GOSEC_ARCH="arm64"
else
    GOSEC_ARCH="amd64"
fi
wget -qO gosec.tar.gz "https://github.com/securego/gosec/releases/download/v${GOSEC_VERSION}/gosec_${GOSEC_VERSION}_linux_${GOSEC_ARCH}.tar.gz"
tar -xzf gosec.tar.gz -C /usr/local/bin gosec
rm gosec.tar.gz
chmod +x /usr/local/bin/gosec
gosec --version
