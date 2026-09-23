#!/bin/sh

# linux build and package

set -eu

# GOOS=linux CGO='aarch64-linux-gnu-gcc' CC='aarch64-linux-gnu-gcc' CGO_ENABLED=1 GOARCH='arm64' go build 

ARM64_PACKAGES="             \
    libgl1-mesa-dev:arm64    \
    libx11-dev:arm64         \
    libxrandr-dev:arm64      \
    libxxf86vm-dev:arm64     \
    libxi-dev:arm64          \
    libxcursor-dev:arm64     \
    libxinerama-dev:arm64    \
    libxext-dev:arm64        \
    libxfixes-dev:arm64      \
    libxdamage-dev:arm64     \
    libxrender-dev:arm64     \
    libx11-xcb-dev:arm64"

AMD64_PACKAGES="       \
    libgl1-mesa-dev    \
    libx11-dev         \
    libxrandr-dev      \
    libxxf86vm-dev     \
    libxi-dev          \
    libxcursor-dev     \
    libxinerama-dev    \
    libxext-dev        \
    libxfixes-dev      \
    libxdamage-dev     \
    libxrender-dev     \
    libx11-xcb-dev"


ARCH="${1:-}"

case "$ARCH" in 
  amd64)
  ;;
  arm64)
  	CGO='aarch64-linux-gnu-gcc'
	CC='aarch64-linux-gnu-gcc'
  	GOARCH="arm64"
  ;;
  *)
    printf "invalid architecture as argument\n" >&2
    exit 1
  ;;
esac
