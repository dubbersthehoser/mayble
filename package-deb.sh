#!/bin/sh

set -eu

ARCH="$1"

PROJECT_ROOT="$PWD"

TAR_FILE="${PWD}/fyne-cross/dist/linux-${ARCH}/Mayble.tar.xz"

if [ ! -e "$TAR_FILE" ]; then
	printf "tar file not found: %s\n" "${TAR_FILE}"
	exit 1
fi

VERSION="$(grep 'Version = ' FyneApp.toml | cut -d ' ' -f 5 | tr -d '"')"

BUILD_ROOT="./build/deb/mayble-${VERSION}"

DEB_DEBIAN="${BUILD_ROOT}/DEBIAN"

if [ -e "${BUILD_ROOT}" ]; then
	echo "[ CLEAN BUILD STAGING ]"
	rm -vrf ${BUILD_ROOT}
fi

mkdir -p "${DEB_DEBIAN}"

cp -v "${TAR_FILE}" "${BUILD_ROOT}"

cd "${BUILD_ROOT}"

tar -xf "${TAR_FILE}"

echo "${VERSION}"

cat > DEBIAN/control << EOF
Section: office
Priority: optional
Maintainer: Brandon Fredericks
Homepage: https://github.com/dubbersthehoser/mayble
Package: mayble
Version: ${VERSION}
Architecture: ${ARCH}
Depends: build-essential
Description: A simple book management system.
EOF

cd "${PROJECT_ROOT}"

dpkg-deb --build "${BUILD_ROOT}"
