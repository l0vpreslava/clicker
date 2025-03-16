#!/usr/bin/env bash

set -xe

export CGO_ENABLED="1"
export CGO_CFLAGS="-mmacosx-version-min=10.12"
export CGO_LDFLAGS="-mmacosx-version-min=10.12"
export GOOS="darwin"

GOARCH="arm64" go build -o clicker_arm64 .
GOARCH="amd64" go build -o clicker_amd64 .
lipo -create -output clicker clicker_amd64 clicker_arm_64
