#!/bin/bash
set -euo pipefail

echo "Installing Trivy dependencies..."
apt-get update
apt-get install -y --no-install-recommends ca-certificates wget gnupg gnupg2 gpg

echo "Installing Trivy..."
wget -qO - https://aquasecurity.github.io/trivy-repo/deb/public.key | gpg --dearmor -o /usr/share/keyrings/trivy.gpg
echo "deb [signed-by=/usr/share/keyrings/trivy.gpg] https://aquasecurity.github.io/trivy-repo/deb generic main" | tee /etc/apt/sources.list.d/trivy.list
apt-get update
apt-get install -y trivy
rm -rf /var/lib/apt/lists/*
trivy --version
