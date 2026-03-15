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
		go test ./internal/... -coverprofile=cover.out
		go tool cover -html cover.out -o cover.html
	;;
	clear)
		rm cover.out cover.html
	;;

	*)
		go test "$cmd" $@ -coverprofile=cover.out
	;;
esac
go tool cover -html cover.out -o cover.html



