#!/bin/bash
set -euo pipefail

echo "Installing Trufflehog dependencies..."
apt-get update
apt-get install -y --no-install-recommends ca-certificates wget tar
rm -rf /var/lib/apt/lists/*

echo "Installing Trufflehog..."
TRUFFLEHOG_VERSION="3.92.4"
# Detect architecture
ARCH=$(uname -m)
if [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
    TRUFFLEHOG_ARCH="arm64"
else
    TRUFFLEHOG_ARCH="amd64"
fi
wget -qO trufflehog.tar.gz "https://github.com/trufflesecurity/trufflehog/releases/download/v${TRUFFLEHOG_VERSION}/trufflehog_${TRUFFLEHOG_VERSION}_linux_${TRUFFLEHOG_ARCH}.tar.gz"
tar -xzf trufflehog.tar.gz -C /usr/local/bin
rm trufflehog.tar.gz
chmod +x /usr/local/bin/trufflehog
trufflehog --version
