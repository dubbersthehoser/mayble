#!/bin/sh

# Linux build and packages from a Ubuntu AMD64 host with ARM64 libs and compiler.

# Allow 'local'
# shellcheck disable=3043

set -eu

echo "-----------------------------------"
echo "Building and Packaging Ubuntu Linux"
echo "-----------------------------------"

NAME="$(grep '^Name' ./FyneApp.toml | cut -d'"' -f 2)"
ICON="$(grep '^Icon' ./FyneApp.toml | cut -d'"' -f 2)"
VERSION="$(grep '^Version' ./FyneApp.toml | cut -d'"' -f2)"

AMD64_NAME="${NAME}-amd64"
ARM64_NAME="${NAME}-arm64"

if [ -f /etc/os-release ]; then
	if ! grep -q 'ubuntu' /etc/os-release; then
		printf "invalid distro\n" 1>&2
		exit 1
	fi
else
	printf "couldn't check distro\n" 1>&2
	exit 1
fi

if [ "$(dpkg --print-architecture)" != 'amd64' ]; then
	printf "invalid host architecture %s\n" "$(dpkg --print-architecture)" 1>&2
	exit 1
fi

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

AMD64_PACKAGES="           \
    libgl1-mesa-dev:amd64  \
    libx11-dev:amd64       \
    libxrandr-dev:amd64    \
    libxxf86vm-dev:amd64   \
    libxi-dev:amd64        \
    libxcursor-dev:amd64   \
    libxinerama-dev:amd64  \
    libxext-dev:amd64      \
    libxfixes-dev:amd64    \
    libxdamage-dev:amd64   \
    libxrender-dev:amd64   \
    libx11-xcb-dev:amd64"


echo
echo "-------------------------"
echo "Installing AMD64 Packages"
echo "-------------------------"


# shellcheck disable=2086
sudo apt install $AMD64_PACKAGES

HAVE_ARM64='true'

if printf "%s" "$(dpkg --print-foreign-architectures)" | grep -q 'arm64'; then
	echo
	echo "-------------------------"
	echo "Installing ARM64 Packages"
	echo "-------------------------"
	# shellcheck disable=2086
	sudo apt install $ARM64_PACKAGES
	
else 
	echo "ARM64 is not set up for cross compile"
	HAVE_ARM64='false'
fi

echo
echo "-----------------"
echo "Building Binaries"
echo "-----------------"

mkdir -vp ./bin 

echo "Building: ./bin/${AMD64_NAME}"
go build -o "./bin/${AMD64_NAME}" .
echo "Completed: ./bin/${AMD64_NAME}"

if [ "${HAVE_ARM64}" = 'true' ]; then
	echo "Building: ./bin/${ARM64_NAME}"
	GOOS=linux                  \
	CGO='aarch64-linux-gnu-gcc' \
	CC='aarch64-linux-gnu-gcc'  \
	CGO_ENABLED=1               \
	GOARCH='arm64'              \
	go build  -o "./bin/${ARM64_NAME}"
	echo "Completed: ./bin/${ARM64_NAME}"
fi


clean_staging() {
	if [ -d ./staging ]; then
		rm -vrf ./staging
	fi

	mkdir -vp "./staging/${NAME}/share/applications" \
	          "./staging/${NAME}/share/bin"          \
	          "./staging/${NAME}/share/pixmaps"
}

files_to_staging() {
	local bin_path="${1}"
	cp -va "./${ICON}" "./staging/${NAME}/share/pixmaps/${NAME}.${ICON##*.}"
	cp -va "$bin_path" "./staging/${NAME}/share/bin/${NAME}"

	cat << EOF > "./staging/${NAME}/share/applications/${NAME}.desktop"
[Desktop Entry]
Type=Application
Name=${NAME}
Exec=${NAME}
Icon=${NAME}.${ICON##*.}
GenericName=Book Management
Categories=Office;Database;
Keywords=books;office;
EOF
}

packup_staging() {
	local arch_name="${1}"
	mkdir -p ./dist
	tar -cvf "./dist/${arch_name}.tar.gz" -C ./staging "./${NAME}"
}

echo
echo "-----------------"
echo "Building Packages"
echo "-----------------"

echo "AMD64..."
clean_staging
files_to_staging "./bin/${AMD64_NAME}"
packup_staging "${AMD64_NAME}-${VERSION}"
echo "Completed."

if [ -f "./bin/${ARM64_NAME}" ]; then
	echo "ARM64..."
	clean_staging
	files_to_staging "./bin/${ARM64_NAME}"
	packup_staging "${ARM64_NAME}-${VERSION}"
	echo "Completed."
fi

echo
echo "------------------"
echo "Archives Completed"
echo "------------------"
find ./dist -type f
