#!/bin/sh

set -eu

arch=""

case "$(uname -m)" in
	x86_64) arch="amd64" ;;
	aarch64) arch="arm64" ;;
	*) echo "failed to determine machine hardware" ; exit 1 ;;
esac

pkg="linux-${arch}_mayble.tar.xz"

curl -LO "https://github.com/dubbersthehoser/mayble/releases/latest/download/${pkg}"

tmpdir="$(mktemp -d)"

mv "${pkg}" "${tmpdir}"

cd "${tmpdir}"

tar -xf "${pkg}"

