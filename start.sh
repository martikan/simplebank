#!/bin/sh

set -e

echo "run db migration..."
/app/migrate -path /app/migration -database "postgresql://$DATASOURCE_USER:$DATASOURCE_PASSWORD@$DATASOURCE_URL/$DATASOURCE_DB" -verbose up

echo "start the app..."
exec "$@"