#!/bin/sh

if [[ $DATASTORE_WRITER_HOST && $DATASTORE_WRITER_PORT ]]; then
    while ! nc -z "$DATASTORE_WRITER_HOST" "$DATASTORE_WRITER_PORT"; do
        echo "waiting for $DATASTORE_WRITER_HOST:$DATASTORE_WRITER_PORT"
        sleep 1
    done

    echo "$DATASTORE_WRITER_HOST:$DATASTORE_WRITER_PORT is available"
fi

exec "$@"
