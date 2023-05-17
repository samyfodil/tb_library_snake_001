#!/bin/bash

set -x

. /utils/wasm.sh

echo "Building ${FILENAME}"

build debug "${FILENAME}"

exit $?