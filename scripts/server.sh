#!/bin/sh

# Postgres
export E_SERVICES_POSTGRES_USER=$(echo ${DATABASE_URL} | grep -o "://[^:]*" | sed 's/^.\{3\}//')
export E_SERVICES_POSTGRES_PASSWORD=$(echo ${DATABASE_URL} | grep -o ":[^:@]*@" | sed 's/^.\{1\}//' | sed 's/.\{1\}$//')
export E_SERVICES_POSTGRES_DB=$(echo ${DATABASE_URL} | grep -o "/[^/]*$" | sed 's/^.\{1\}//')
export E_SERVICES_POSTGRES_HOST=$(echo ${DATABASE_URL} | grep -o "@[^:]*" | sed 's/^.\{1\}//')
export E_SERVICES_POSTGRES_PORT=$(echo ${DATABASE_URL} | grep -o ":[0-9]\+" | sed 's/^.\{1\}//')

# Redis
export E_SERVICES_REDIS_ADDRESS=$(echo ${REDIS_URL} | grep -o "@[^:]*" | sed 's/^.\{1\}//')
export E_SERVICES_REDIS_PORT=$(echo ${REDIS_URL} | grep -o "[0-9]\+$")
export E_SERVICES_REDIS_PASSWORD=$(echo ${REDIS_URL} | grep -o ":[^:@]*@" | sed 's/^.\{1\}//' | sed 's/.\{1\}$//')

export PORT=80
/bin/migration --action=up
/bin/server
