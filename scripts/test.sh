#!/bin/bash -e

# The root of the build directory
ROOT=$(
  unset CDPATH
  root=$(dirname "${BASH_SOURCE}")
  cd "${root}"
  pwd
)

export GOPATH=${ROOT}:${ROOT}/Godeps/_workspace

go test ./... $*
