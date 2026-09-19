#!/bin/sh

set -eu

mkdir -p ./bin
go build -o ./bin/testbackend ./cmd/testbackend

pids=""

cleanup() {
	trap - EXIT INT TERM
	[ -z "$pids" ] || kill $pids 2>/dev/null || true
	wait $pids 2>/dev/null || true
}

trap cleanup EXIT
trap 'exit 0' INT TERM

PORT=9001 BACKEND_ID=backend-1 ./bin/testbackend & pids="$pids $!"
PORT=9002 BACKEND_ID=backend-2 ./bin/testbackend & pids="$pids $!"
PORT=9003 BACKEND_ID=backend-3 ./bin/testbackend & pids="$pids $!"

wait $pids
