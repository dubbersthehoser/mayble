#!/bin/sh

set -eu

if "$(uname)" != Linux; then
  printf "Invalid OS $(uname). Linux only\n"
fi

ARCH=${-:NONE}

case "$(arch)" in 
  x86_64)  ARCH=amd64 ;;
  aarch64) ARCH=amd64 ;;
  *) 
    printf "Invalid CPU architecture $(arch)\n" 1>&2
    exit 1
  ;;
esac

dl_url="$(curl -s https://api.github.com/repos/dubbersthehoser/mayble/releases/latest | \
  grep browser_download_url | \
  cut -d'"' -f 4 | /
  grep ${ARCH})"

DL_DIR="mayble-dist"
ARCHIVE="mayble.tar.xz"

mkdir -vp "$DL_DIR"

# download archive
printf "Downloading Archive:\n  - %s\n" "$dl_url"

curl -sL "$dl_url" -o "${DL_DIR}/${ARCHIVE}"

# extract package.

cd -v "$DL_DIR"

printf "Extracting:\n  - %s\n" "${ARCHIVE}"

tar -xf "$ARCIVE"

# install package to user install.

make user-install # make may not be installed?
