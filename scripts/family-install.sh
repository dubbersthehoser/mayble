#!/bin/sh


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
# install to the user directory.


