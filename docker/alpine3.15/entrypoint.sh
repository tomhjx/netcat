#!/bin/sh
set -e

if [ "$1" == 'bash' ]; then
    exec "$@"
    return 1
fi

/mybin/run.sh $@