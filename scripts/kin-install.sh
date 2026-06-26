#!/bin/sh

# remove the old version from kin's chrome book linux env.
mayble_installed="$(apt-cache pkgnames | grep 'mayble')"
if [ -n "$mayble_installed" ]; then
	sudo apt remove mayble 
	status="$?"
	mayble_installed="$(apt-cache pkgnames | grep 'mayble')"
	if [ "$status" -ne 0 ] || [ -n "$mayble_intalled" ] ; then
		echo "failed to be removed."
		echo "aborting intall."
		exit 1;
	fi
fi

echo "continue install"

# 1. download release.
# 2. extract it.
# 3. check arch.
# 4. install the binary and dot desktop file.
# 5. clean up.

