#!/bin/bash

. /utils/wasm.sh

debug_build "${FILENAME}"
ret=$?
echo -n $ret > /out/ret-code
exit $ret
