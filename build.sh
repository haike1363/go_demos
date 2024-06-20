#!/bin/bash
SHELL_DIR=$(
  cd "$(dirname "$0")" || exit
  pwd
)
cd "${SHELL_DIR}" || exit

set -eu

MODULE=$1

echo build ${MODULE}

export GOROOT=/Users/easechen/apps/go1.20.3
export PATH=$GOROOT/bin:$PATH

go mod tidy
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
export GOPRIVATE=git.code.oa.com,git.woa.com

BUILD_FLAGS="-tags static_all"
go mod tidy
go build ${BUILD_FLAGS} -v -o bin/${MODULE}_test ./pkg/${MODULE}/main.go
