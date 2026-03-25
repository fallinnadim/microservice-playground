#!/bin/sh
set -e

echo "Running custom SQL migrations..."

for file in /migrations/*.sql; do
  echo "Executing $file"
  psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f "$file"
done

echo "Migrations completed!"