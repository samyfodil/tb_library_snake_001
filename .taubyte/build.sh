#!/bin/bash

set -x

. /utils/wasm.sh

echo "Building ${FILENAME}"

build debug "${FILENAME}"
ret=$?

echo -n $ret > /out/ret-code

exit $ret
