#!/bin/bash
set -euo pipefail

echo "Installing Semgrep dependencies..."
apt-get update
apt-get install -y --no-install-recommends python3 python3-pip ca-certificates git
rm -rf /var/lib/apt/lists/*

echo "Installing Semgrep..."
pip3 install --no-cache-dir semgrep --break-system-packages
semgrep --version
