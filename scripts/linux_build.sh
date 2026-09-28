#!/bin/sh

set -eu

[ ! -d "$DIST" ] && mkdir -v "$DIST"

go clean -cache

for ARCH in "amd64" "arm64"; do
	CGO_ENABLED=1  \
	GOOS="linux"   \
	GOARCH="$ARCH" \
	CC=""   \
	go build -v -o "$DIST/mayble-${VERSION}_linux-${ARCH}" .
done
