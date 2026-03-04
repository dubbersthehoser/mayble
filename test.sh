#!/bin/sh

set -e

if [ $# -eq 0 ]; then
	echo "no arguments or package were given" 1>&2
	exit 1
fi

cmd="$1"
shift

case "$cmd" in
	total)
		go test ./internal/... -coverprofile=cover.out
		go tool cover -func=cover.out
	;;

	all)
		go test ./internal/... -cover
	;;
	*)
		go test $@
	;;
esac



