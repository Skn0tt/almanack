#!/bin/bash

set -eu -o pipefail

# Get the directory that this script file is in
THIS_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

cd "$THIS_DIR"

function _default() {
	start-api
}

function _die() {
	>&2 echo "Fatal: ${@}"
	exit 1
}

function _installed() {
	which "$1" >/dev/null 2>&1;
}

function help() {
	local SCRIPT=$0
	cat <<EOF
Usage

	$SCRIPT <task> <args>

Tasks:

EOF
	compgen -A function | grep -e '^_' -v | sort | xargs printf ' - %s\n'
	exit 2
}

function redis-start() {
	docker run --name redis-container --rm -p 6379:6379 redis:latest
}

function redis-cli() {
	docker run --name redis-cli -it --rm \
		--link redis-container:redis \
		redis:latest \
		redis-cli -h redis -p 6379
}

function sql() {
	_installed sqlc || _die "sqlc not installed"
	go generate ./...
}

function migrate() {
	cd sql/schema
	tern migrate -c prod.conf
}

function migrate:prod() {
	cd sql/schema
	tern migrate -c prod.conf
}

function build:prod() {
	./build.sh
}

function test() {
	test:backend
	test:frontend
	test:misc
}

function test:frontend() {
	yarn run test
}

function test:backend() {
	go test ./... -v
}

function test:misc() {
	shellcheck ./run.sh
}

function format() {
	yarn run lint
	gofmt -s -w .
	shfmt -w ./run.sh
}

function start-api() {
	set -x
	go run ./funcs/almanack-api -src-feed "$ARC_FEED" -cache
}

TIMEFORMAT="Task completed in %1lR"
time "${@:-_default}"
