#!/bin/sh

set -eu

die() {
	arg="$1"
	echo "$arg" 1>&2
	exit 1
}

apply_line() {
	FILE="$1"; shift
	PATT="$1"; shift
	LINE="$1"; shift

	if [ -f "$FILE" ]; then 
		TMP=$(mktemp)
		sed "s:$PATT:$LINE:" "$FILE" >> "$TMP"
		mv -v "$TMP" "$FILE"
	else 
		echo "$FILE not found." 1>&2
	fi
}

ARG="${1:-}"

VERSION_FILE="version.txt"

[ -f "$VERSION_FILE" ] || die "missing version file: '${VERSION_FILE}'"

VERSION="$(grep '^[0-9]\+\.[0-9]\+\.[0-9]\+$' "$VERSION_FILE")" || die "grep failed: $?"

case $ARG in

	# Print out the version with flag argument.
	-v | --version)
		printf "%s\n" "$VERSION"
		exit 0
		;;
	# Git tag current .
	-t | --tag)
		git tag -a "v${VERSION}" -m "application version ${VERSION}"
		exit 0
		;;
esac

# Error when invalid flag argument was given.
if [ -n "$ARG" ]; then
	die "invalid argument flag '$ARG'"
fi

[ -z "$VERSION" ] && die "could not find version"

apply_line "FyneApp.toml"      \
	    "^  Version = .*$"  \
	    "  Version = \"$VERSION\" # modified by $0" 

# Add version to package config to config.go file.
apply_line "./internal/config/config.go"  \
	    "^const Version string = .*"   \
	    "const Version string = \"$VERSION\" // modified by $0"

# change version to info file.
apply_line "./doc/info.md"  \
	    "^Mayble v.*"  \
	    "Mayble v$VERSION"

# change version to manual file.
apply_line "./doc/manual.md"  \
	    "^Mayble v.*"  \
	    "Mayble v$VERSION"

