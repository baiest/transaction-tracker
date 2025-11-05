#!/bin/sh
set -e

until nc -z -v -w30 postgres 5432
do
  sleep 2
done

migrate -path /app/migrations \
  -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable" up

exec /app/api-bin
