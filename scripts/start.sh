#!/bin/sh

set -e

echo "run db migration"
goose -dir /app/schema up

echo "start app"
exec "$@"

