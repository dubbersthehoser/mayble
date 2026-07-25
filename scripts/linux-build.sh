#!/bin/sh

set -eu

DIST="./dist"
VERSION="$(sh ./scripts/version.sh --version)"

[ ! -d "$DIST" ] && mkdir -v "$DIST"

go clean -cache

for ARCH in "amd64" "arm64"; do
	CGO_ENABLED=1  \
	GOOS="linux"   \
	GOARCH="$ARCH" \
	CC=""   \
	go build -v -o "$DIST/mayble-${VERSION}_linux-${ARCH}" .
done
