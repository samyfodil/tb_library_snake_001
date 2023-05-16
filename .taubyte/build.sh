#!/bin/bash

. /utils/wasm.sh

debug_build 1 "${FILENAME}"
ret=$?
echo -n $ret > /out/ret-code
exit $ret
