#!/bin/bash

set -euo pipefail

readonly top=$(git rev-parse --show-toplevel)

pushd $top/sshtest/testdata

docker build -t sshtest-ubuntu:latest -f Dockerfile.ubuntu  .

popd
