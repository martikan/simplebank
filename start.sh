#!/bin/sh

set -e

echo "run db migration..."
echo "postgres://$DATASOURCE_USER:$DATASOURCE_PASSWORD@$SERVER_ADDRESS/$DATASOURCE_DB?sslmode=disable"
/app/migrate -path /app/migration -database "postgres://$DATASOURCE_USER:$DATASOURCE_PASSWORD@$SERVER_ADDRESS/$DATASOURCE_DB?sslmode=disable" -verbose up

echo "start the app..."
exec "$@"