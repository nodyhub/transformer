#!/bin/bash
set -euo pipefail

CODEQL_VERSION="2.23.8"
ARCH=$(uname -m)
OS=$(uname -s)

echo "Detected OS: $OS, Arch: $ARCH"


if [ "$OS" = "Darwin" ]; then
    # Native Mac (Universal Binary - works on ARM64 and Intel)
    DOWNLOAD_URL="https://github.com/github/codeql-cli-binaries/releases/download/v${CODEQL_VERSION}/codeql-osx64.zip"
elif [ "$ARCH" = "x86_64" ]; then
    DOWNLOAD_URL="https://github.com/github/codeql-cli-binaries/releases/download/v${CODEQL_VERSION}/codeql-linux64.zip"
else
    echo "CodeQL installation is only supported on x86_64 and Mac for now. Skipping."
    exit 0
fi

echo "Downloading from: $DOWNLOAD_URL"
wget -qO codeql.zip "$DOWNLOAD_URL"
unzip -q codeql.zip -d /opt
rm codeql.zip
ln -sf /opt/codeql/codeql /usr/local/bin/codeql

# Verify (This will trigger Rosetta if on Linux ARM64)
codeql version