#!/bin/sh

set -eux

# Note: LLM Generated
# Setup apt for arm64 packages for releases.

sudo dpkg --add-architecture arm64

# Restrict the runner's existing repositories to amd64.
if [ -f /etc/apt/sources.list.d/ubuntu.sources ]; then
	sudo sed -i '/^Types:/a Architectures: amd64' \
	/etc/apt/sources.list.d/ubuntu.sources
else
	sudo sed -i -E 's|^(deb )|\1[arch=amd64] |' \
	/etc/apt/sources.list
fi

# Add arm64 repositories from the Ubuntu ports mirror.
sudo tee /etc/apt/sources.list.d/arm64.sources >/dev/null <<'EOF'
Types: deb
URIs: http://ports.ubuntu.com/ubuntu-ports/
Suites: noble noble-updates noble-backports
Components: main restricted universe multiverse
Architectures: arm64
Signed-By: /usr/share/keyrings/ubuntu-archive-keyring.gpg

Types: deb
URIs: http://ports.ubuntu.com/ubuntu-ports/
Suites: noble-security
Components: main restricted universe multiverse
Architectures: arm64
Signed-By: /usr/share/keyrings/ubuntu-archive-keyring.gpg
EOF

sudo apt-get update
