#!/usr/bin/env bash

docker run --rm -i \
--name stats-to-scale-postgres \
-p 5432:5432 \
--env TIMESCALEDB_TELEMETRY=off \
-e POSTGRES_PASSWORD=artilidus \
-e POSTGRES_USER=artilidus \
-e POSTGRES_DB=artilidus \
-v "${PWD}/../database/schema.sql":/docker-entrypoint-initdb.d/1_schema.sql \
timescale/timescaledb:latest-pg11-oss
