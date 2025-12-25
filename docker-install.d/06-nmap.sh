#!/bin/bash
set -euo pipefail

echo "Installing Nmap..."
apt-get update
apt-get install -y --no-install-recommends ca-certificates nmap
rm -rf /var/lib/apt/lists/*
nmap --version
