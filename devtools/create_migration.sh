#!/usr/bin/env -S bash -e
# Script for creating a new migration file.

migrate create -ext sql -dir migrations -seq $1
HEADER="BEGIN;

-- Write your migration here.

END;"
for f in $(ls migrations | tail -n 2); do echo "$HEADER" >>"migrations/$f"; done