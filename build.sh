#!/bin/sh

set -eu

COMMAND="$1"

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
	go test \
	./internal/app     \
	./internal/config  \
	./internal/gui/... \
	./internal/listing \
	./internal/emiter \
	./internal/command \
	./internal/porting/... \
	./internal/searching \
	./internal/storage \
	./internal/sqlite \
	"$2" \

	;;
	*)
	echo "command not found"
	;;
esac
