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

if [ "$(go version | grep 1.11)" != "" ]; then
	export GO111MODULE=on
fi

[ -d vendor/github.com/ligato/crd-example/ ] && (echo "Run: rm -rf vendor/github.com/ligato/crd-example;dep ensure";exit 1)
test -z "$(go fmt ./...)" || (echo "Run go fmt ./... and recommit your code";exit 1)
go get -u github.com/golang/protobuf/protoc-gen-go
go get -u k8s.io/code-generator
go generate ./...
go build ./...
if [ "$(go version | grep 1.11)" == "" ]; then
	go test $RACE_ENABLED ./...
else
	go test -mod=vendor $RACE_ENABLED ./...
fi
go install ./...
