#!/bin/sh

set -eu

COMMAND="${1:-}"

if [ -z "${COMMAND}" ]; then
	echo "command not given"
	exit 1
fi

case "${COMMAND}" in
	build)
	echo "build has been run"
	#go build ./cmd/gui -o ./bin/mayble
	;;
	test)
	go test 
	;;
	debian | deb)
	package_deb
	;;
	*)
	echo "command not found"
	;;
esac

package_deb() {
	fyne package
}
