#!/bin/sh

set -eu

ARCH="$1"

PROJECT_ROOT="$PWD"

TAR_FILE="${PWD}/fyne-cross/dist/linux-${ARCH}/mayble.tar.xz"

if [ ! -e "$TAR_FILE" ]; then
	printf "tar file not found: %s\n" "${TAR_FILE}"
	exit 1
fi

VERSION="$(grep 'Version = ' FyneApp.toml | cut -d ' ' -f 5 | tr -d '"')"

BUILD_ROOT="./build/deb/mayble-${VERSION}"

DEB_DEBIAN="${BUILD_ROOT}/DEBIAN"

if [ -e "${BUILD_ROOT}" ]; then
	echo "[ CLEAN BUILD DIRECTORY ]"
	rm -vrf ${BUILD_ROOT}
fi

echo "[ SETTING UP ]"

mkdir -p "${DEB_DEBIAN}"

cd "${BUILD_ROOT}"
echo cd "${BUILD_ROOT}"

echo extracting...
tar -xvf "${TAR_FILE}"

ICON="$(grep 'Icon := .*' Makefile | cut -d ' ' -f 3 | tr -d '"')"

APPID="$(grep 'Name := .*' Makefile | cut -d ' ' -f 3 | tr -d '"')"

printf "APPID: %s\n" "${APPID}"
printf "ICON: %s\n" "${ICON}"

echo proper pathing...
mkdir -v -p ./usr/share/pixmaps ./usr/share/applications

mv -v ./usr/local/share/pixmaps/mayble.png "./usr/share/pixmaps/${ICON}"

mv -v ./usr/local/share/applications/mayble.desktop "./usr/share/applications/${APPID}.desktop"


echo removing empty directories...
find ./usr -type d -empty -delete 

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

cp -v LICENSE "${DEB_DEBIAN}/copyright"

echo "[ BUILD ]"

dpkg-deb --build "${BUILD_ROOT}"
