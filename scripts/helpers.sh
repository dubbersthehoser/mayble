#!/bin/sh
#
# Helper functions for scripts
#
# Allow 'local'
# shellcheck disable=3043

#######################################
# Get configuration values.
# Arguments:
#   The key for value to be printed
#
config_get() {

	local key="$1"
	case "$key" in
	app-version)
		grep '^Version' ./FyneApp.toml | cut -d'"' -f 2
	;;
	app-icon)
		grep '^Icon' ./FyneApp.toml | cut -d'"' -f 2
	;;
	app-name)
		grep '^Name' ./FyneApp.toml | cut -d'"' -f 2
	;;
	*)
		printf "config_get %s: key not found"
		exit 1
	;;
	esac
}

#######################################
# Print out ubuntu/debian dependencies for insalling.
# Arguments:
#   Optional suffix for package architectures e.g 'arm64' and 'amd64'
# Globals:
#   None
#
ubuntu_packages() {
	local suffix
	local arch="${1:-__EMPTY__}"
	local deps="
	libgl1-mesa-dev  \
	libx11-dev       \
	libxrandr-dev    \
	libxxf86vm-dev   \
	libxi-dev        \
	libxcursor-dev   \
	libxinerama-dev  \
	libxext-dev      \
	libxfixes-dev    \
	libxdamage-dev   \
	libxrender-dev   \
	libx11-xcb-dev"

	case "$arch" in
		amd64|AMD64)
			suffix=":amd64"
		;;
		arm64|ARM64)
			suffix=":arm64"
		;;
		__EMPTY__)
			suffix=""
		;;
		*)
			printf "Invalid Aarchitecture '%s'\n" "$arch" 1>&2
			exit 1
		;;
	esac

	items=""

	for item in $deps; do
		items="${items} ${item}${suffix}"
	done

	printf "%s\n" "${items}"
}

linux_files_to_staging() {
	local bin_path="${1}"
	local version
	local name
	local icon
	version="$(config_get app-version)"
	icon="$(config_get app-icon)"
	name="$(config_get app-name)"

	cp -va "./${icon}" "./staging/${name}/share/pixmaps/${name}.${icon##*.}"
	cp -va "$bin_path" "./staging/${name}/share/bin/${name}"

	printf "%s\n" "mayble-${version}" > "./staging/${name}/version.txt"

	cat << EOF > "./staging/${name}/share/applications/${name}.desktop"
[Desktop Entry]
Type=Application
Name=${name}
Exec=${name}
Icon=${name}.${icon##*.}
GenericName=Book Management
Categories=Office;Database;
Keywords=books;office;
EOF
}


packup_linux_to_dist() {
	local name="${1}"
	mkdir -p ./dist
	tar -avcf "./dist/${name}.tar.gz" -C ./staging "./$(config_get app-name)"
}


clear_staging() {
	[ -d ./staging ] && rm -vrf ./staging
	mkdir -v ./staging
}

setup_staging_to_linux() {
	clear_staging
	local name
	name="$(config_get app-name)"
	mkdir -vp "./staging/${name}/share/applications" \
	          "./staging/${name}/share/bin"          \
	          "./staging/${name}/share/pixmaps"
}

log_fatal() {
	local msg="$1"
	printf "%s: %s %s\n" "$(date '+F %R')" "${0##*/}" "$msg" 1>&2
	exit 1
}
