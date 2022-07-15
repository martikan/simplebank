#!/bin/sh

set -e

echo "run db migration..."
echo "postgresql://$DATASOURCE_USER:$DATASOURCE_PASSWORD@$SERVER_ADDRESS/$DATASOURCE_DB"
/app/migrate -path /app/migration -database "postgresql://$DATASOURCE_USER:$DATASOURCE_PASSWORD@$SERVER_ADDRESS/$DATASOURCE_DB" -verbose up

echo "start the app..."
exec "$@"