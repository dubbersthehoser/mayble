#!/bin/sh

set -eu

die() {
	arg="$1"
	echo "$arg" 1>&2
	exit 1
}

VERSION_FILE="version.txt"

[ -f "$VERSION_FILE" ] || die "missing version file: '${VERSION_FILE}'"

VERSION="$(grep '^[0-9]\+\.[0-9]\+\.[0-9]\+$' "$VERSION_FILE")" || die "grep failed: $?"

[ -z "$VERSION" ] && die "could not find version"

FYNE_TOML="FyneApp.toml"
if [ -f "$FYNE_TOML" ]; then 
	TMP=$(mktemp)
	line="  Version = \"$VERSION\" # modified by $0"
	sed "s:^  Version = .*$:$line:" "$FYNE_TOML" >> "$TMP"
	mv -v $TMP "$FYNE_TOML"
else 
	echo "$FYNE_TOML not found." 1>&2
fi

APP_CONFIG="./internal/config/config.go"
if [ -f "$APP_CONFIG" ]; then 
	TMP=$(mktemp)
	line="const Version string = \"$VERSION\" // modified by $0"
	sed "s:^const Version string = .*:$line:" "$APP_CONFIG" >> "$TMP"
	mv -v "$TMP" "$APP_CONFIG"
else 
	echo "$APP_CONFIG not found." 1>&2
fi
