#!/usr/bin/env -S bash -e

docker run -d -p 5432:5432 --name local-postgres \
  -e LANG=C \
  -e POSTGRES_DB=${GO_DISCOVERY_DATABASE_NAME:-discovery-db} \
  -e POSTGRES_USER=${GO_DISCOVERY_DATABASE_USER:-postgres} \
  -e POSTGRES_PASSWORD=${GO_DISCOVERY_DATABASE_PASSWORD:-postgres} \
  postgres:11.12
