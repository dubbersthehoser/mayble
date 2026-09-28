#!/bin/sh

set -eu

# remove the old version from kin's chrome book linux env.
mayble_installed="$(apt-cache pkgnames | grep 'mayble')"
if [ -n "$mayble_installed" ]; then
	sudo apt remove mayble 
	status="$?"
	mayble_installed="$(apt-cache pkgnames | grep 'mayble')"
	if [ "$status" -ne 0 ] || [ -n "$mayble_intalled" ] ; then
		echo "failed to be removed."
		echo "aborting install."
		exit 1;
	fi
fi

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

# 1. download release.
echo "## Downloading Release ##"
RELEASE_URL="$(curl -s https://api.github.com/repos/dubbersthehoser/mayble/releases/latest" \
	| grep browser_download_url                                                         \
	| cut -d'\"' -f 4                                                                   \
	| grep $ARCH)"

ARCHIVE="${RELEASE_URL##*/}"

curl -sL "${RELEASE_URL}" -o "${ARCHIVE}"

# 2. extract it.
tar -xvf "$ARCHIVE" 

# 3. install

pushd ./mayble

[ ! -x ./install.sh ] && chmod 744 ./install.sh
./install.sh user-install

popd

rm -vrf ./mayble
