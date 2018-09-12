#!/bin/bash

# The only argument this script should ever be called with is '--verify-only'

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[@]}")/..
if [ "$(go version | grep 1.11)" == "" ]; then
	go get -u k8s.io/code-generator
	CODEGEN_PKG="${GOPATH}/src/k8s.io/code-generator"
else
	env GO111MODULE=on go get -u k8s.io/code-generator
	CODEGEN_PKG="${GOPATH}/pkg/mod/k8s.io/code-generator"
fi

echo "Calling ${CODEGEN_PKG}/generate-groups.sh"
"${CODEGEN_PKG}"/generate-groups.sh all \
  github.com/ligato/crd-example/pkg/client github.com/ligato/crd-example/pkg/apis \
  crdexample.io:v1 \
  --output-base "${GOPATH}/src/" \
  --go-header-file "${SCRIPT_ROOT}/conf/boilerplate.txt"

echo "Generating other deepcopy funcs"
"${GOPATH}"/bin/deepcopy-gen \
  --input-dirs ./pkg/crdexample \
  --go-header-file "${SCRIPT_ROOT}/conf/boilerplate.txt" \
  --bounding-dirs ./pkg/crdexample \
  -O zz_generated.deepcopy \
  -o "${GOPATH}/src"

echo "Generating openapi structures"
"${GOPATH}"/bin/openapi-gen \
  --input-dirs ./pkg/apis/crdexample.io/v1 --input-dirs ./pkg/crdexample \
  --output-package github.com/ligato/crd-example/pkg/apis/crdexample.io/v1 \
  --go-header-file "${SCRIPT_ROOT}/conf/boilerplate.txt"
