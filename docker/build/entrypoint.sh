#!/bin/sh
set -e

if [ "$1" == '/bin/sh' ]; then
    exec "$@"
    return 1
fi

/mybin/run.sh $@