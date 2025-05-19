#!/bin/sh
CMD=/app/0_Introduction/matrixMul/matrixMul

if [ -n "${DAEMON_MODE}" ]; then
  CMD="$CMD -d"
fi

if [ -n "${MATRIX_SIZE}" ]; then
  CMD="$CMD
    -wA=${MATRIX_SIZE}
    -hA=${MATRIX_SIZE}
    -wB=${MATRIX_SIZE}
    -hB=${MATRIX_SIZE}
  "
fi

exec ${CMD}
