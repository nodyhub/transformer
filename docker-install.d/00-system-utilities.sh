#!/bin/bash
set -euo pipefail

echo "Installing system libraries..."
apt-get update
apt-get install -y --no-install-recommends \
	ca-certificates \
	wget \
	curl \
	jq \
	unzip \
	bzip2 \
	xz-utils \
	gnupg \
	gnupg2 \
	gpg
