#!/bin/sh

set -eu

fyne-cross linux -arch amd64,arm64

fragments="
	./fyne-cross/dist/linux-amd64/mayble.tar.xz
	./fyne-cross/dist/linux-arm64/mayble.tar.xz
"

for path in ${fragments}; do
	if [ ! -f ${path} ]; then
		echo "file not found: '${path}'"
		exit 1
	fi
done

release_dir="/tmp/mayble-release/"

if [ ! -d "${release_dir}" ]; then
	mkdir "${release_dir}"
fi

for path in ${fragments}; do
	file="$(basename ${path})"
	platform="$(echo "${path}" | cut -d '/' -f 4)"
	release="${platform}_${file}"
	cp -f "${path}" "${release_dir}${release}"
	echo "${release_dir}${release}"
done

