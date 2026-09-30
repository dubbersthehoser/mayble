#!/bin/sh

set -eu

# remove the old version from kin's chrome book linux env.
echo "-- Checking for Installed v1.0.0 Version"
if apt-cache show 'mayble' > /dev/null; then
	echo "  unistalling..."
	if sudo apt remove -y 'mayble' ; then
		:
	else
		echo "  failed to uninstall" 1>&2
		echo "  aborting install." 1>&2
		exit 1;
	fi
	if apt-cache show 'mayble' > /dev/null; then
		echo "  failed to uninstall" 1>&2
		echo "  aborting install." 1>&2
		exit 1;
	fi
	echo "  uninstall completed"
fi

ARCH=""

echo "-- Checking Architecture"

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

curl -sL "${RELEASE_URL}" -o "${ARCHIVE}" > /dev/null

# 2. extract it.
echo "  extracting: ${ARCHIVE}"
tar -xvf "$ARCHIVE" 

# 3. install
cd ./mayble

echo "-- Install"
[ ! -x ./install.sh ] && chmod 744 ./install.sh
./install.sh user-install

cd ..
echo "  completed"
