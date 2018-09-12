#!/bin/bash
#
# This script builds the networkservicemesh
#

set -xe

export RACE_ENABLED="-race"
while test $# -gt 0; do
	case "$1" in
		--race-test-disabled)
			export RACE_ENABLED=""
			;;
		*)
			break
			;;
	esac
	shift
done

GOTESTOPTS=""
if [ "$(go version | grep 1.11)" != "" ]; then
	export GO111MODULE=on
	GOTESTOPTS="-mod=vendor"
fi

[ -d vendor/github.com/ligato/crd-example/ ] && (echo "Run: rm -rf vendor/github.com/ligato/crd-example;dep ensure";exit 1)
test -z "$(go fmt ./...)" || (echo "Run go fmt ./... and recommit your code";exit 1)
go get -u github.com/golang/protobuf/protoc-gen-go
go generate ./...
go build ./...
go test "${GOTESTOPTS}" $RACE_ENABLED ./...
go install ./...
