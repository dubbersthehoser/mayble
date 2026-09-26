#!/bin/sh
#
# Linux build and package from a Debian as AMD64 host with ARM64 packages.
# When dpkg --print-foreign-architectures has arm64 added then the ARM64 binary will be created.
# Otherwise the AMD64 version will be created only.
#
# Allow 'local'
# shellcheck disable=3043

set -eu

package_linux() {
	local bin="${1:-}"
	local name="${2:-}"

	[ -z "$bin" ] && log_fatal "binary argumnet was not given"
	[ ! -f "$bin" ] && log_fatal " '${bin}' binary does not exists or is a directory"


	echo "-- setting up staging --"
	setup_staging_to_linux
	echo "-- files to staging --"
	linux_files_to_staging "$bin"
	echo "-- packing up staging --"
	packup_linux_to_dist "$name"
	echo "-- clearing staging --"
	clear_staging
}

echo "------------------------------------------"
echo "Building and Packaging for Linux on Debian"
echo "------------------------------------------"

. scripts/helpers.sh

NAME="$(config_get app-name)"
AMD64_NAME="${NAME}-amd64"
ARM64_NAME="${NAME}-arm64"

if [ -f /etc/os-release ]; then
	if ! grep -q 'debian' /etc/os-release; then
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

echo
echo "----------------------------"
echo "Installing AMD64 Depedencies"
echo "----------------------------"


# shellcheck disable=2046
sudo apt install $(ubuntu_packages amd64)

HAVE_ARM64='true'

if printf "%s" "$(dpkg --print-foreign-architectures)" | grep -q 'arm64'; then
	echo
	echo "----------------------------"
	echo "Installing ARM64 Depedencies"
	echo "----------------------------"
	# shellcheck disable=2046
	sudo apt install $(ubuntu_packages arm64)
	
else 
	echo "Warning: system is not set up for ARM64 cross compile"
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


echo
echo "-----------------"
echo "Building Packages"
echo "-----------------"

echo
echo "[ AMD64 ]"

package_linux "./bin/${AMD64_NAME}" "${AMD64_NAME}"
echo "Completed."

if [ -f "./bin/${ARM64_NAME}" ]; then
	echo
	echo "[ ARM64 ]"
	package_linux "./bin/${ARM64_NAME}" "${ARM64_NAME}"
	echo "Completed."
fi

echo
echo "------------------"
echo "Packages Completed"
echo "------------------"
find ./dist -type f
