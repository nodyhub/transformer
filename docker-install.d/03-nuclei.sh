#!/bin/bash
set -euo pipefail

echo "Installing Nuclei dependencies..."
apt-get update
apt-get install -y --no-install-recommends ca-certificates wget unzip
rm -rf /var/lib/apt/lists/*

echo "Installing Nuclei..."
NUCLEI_VERSION="3.3.6"
wget -qO nuclei.zip "https://github.com/projectdiscovery/nuclei/releases/download/v${NUCLEI_VERSION}/nuclei_${NUCLEI_VERSION}_linux_amd64.zip"
unzip nuclei.zip -d /usr/local/bin
rm nuclei.zip
chmod +x /usr/local/bin/nuclei
nuclei -version
