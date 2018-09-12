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

[ -d vendor/github.com/ligato/crd-example/ ] && (echo "Run: rm -rf vendor/github.com/ligato/crd-example;dep ensure";exit 1)
test -z "$(env GO111MODULE=on go fmt ./...)" || (echo "Run go fmt ./... and recommit your code";exit 1)
env GO111MODULE=on go get -u github.com/golang/protobuf/protoc-gen-go
env GO111MODULE=on go get -u k8s.io/code-generator
env GO111MODULE=on go generate ./...
env GO111MODULE=on go build ./...
if [ "$(go version | grep 1.11)" == "" ]; then
	env GO111MODULE=on go test $RACE_ENABLED ./...
else
	env GO111MODULE=on go test -mod=vendor $RACE_ENABLED ./...
fi
env GO111MODULE=on go install ./...
