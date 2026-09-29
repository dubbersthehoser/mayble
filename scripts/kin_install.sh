#!/bin/sh

set -eu

# remove the old version from kin's chrome book linux env.
echo "-- Checking for Installed v1.0.0 Version"
mayble_installed="$(apt-cache pkgnames | awk 'mayble' )"
if [ -n "$mayble_installed" ]; then
	echo "-- Unistalling v1.0.0"
	sudo apt remove mayble 
	status="$?"
	mayble_installed="$(apt-cache pkgnames | awk '/mayble/')"
	if [ "$status" -ne 0 ] || [ -n "$mayble_intalled" ] ; then
		echo "  failed to be removed" 1>&2
		echo "  aborting install." 1>&2
		exit 1;
	fi
fi

ARCH=""

echo "-- Checking Arch"

case "$(arch)" in
	x86_64) 
		ARCH=amd64 
	;;
	aarch64) 
		ARCH=arm64 
	;;
	*)
		printf "invalid architecture $(arch)\n" 1>&2
		exit 1
	;;
esac

echo "  architecture: ${ARCH}"

# 1. download release.
echo "-- Downloading Release"
RELEASE_URL="$(curl -sL https://api.github.com/repos/dubbersthehoser/mayble/releases/latest \
	| awk '/browser_download_url/'                                                      \
	| cut -d\" -f4                                                                      \
	| awk "/$ARCH/")"

echo "  downloading: ${RELEASE_URL}"

ARCHIVE="${RELEASE_URL##*/}"

curl -sL "${RELEASE_URL}" -o "${ARCHIVE}"

# 2. extract it.
echo "  extracting: ${ARCHIVE}"
tar -xvf "$ARCHIVE" 

# 3. install
cd ./mayble

[ ! -x ./install.sh ] && chmod 744 ./install.sh
./install.sh user-install

cd ..
