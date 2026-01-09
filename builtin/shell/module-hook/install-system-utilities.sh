#!/bin/bash

set -eu pipefail

echo "Installing system libraries for shell scripts..."
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

# clean up apt cache
rm -rf /var/lib/apt/lists/*

# exit 1