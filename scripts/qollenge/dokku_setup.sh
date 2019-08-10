#!/bin/bash

dokku plugin:install https://github.com/dokku/dokku-letsencrypt.git || dokku plugin:update letsencrypt
dokku plugin:install https://github.com/dokku/dokku-redis.git redis || dokku plugin:update redis
dokku plugin:install https://github.com/dokku/dokku-postgres.git postgres || dokku plugin:update postgres
dokku plugin:install https://github.com/dokku/dokku-memcached.git memcached || dokku plugin:update memcached

dokku apps:create engine
dokku redis:create engine_redis
dokku redis:connect engine_redis engine
dokku postgres:create engine_postgres
dokku postgres:connect engine_postgres engine

dokku config:set --no-restart engine DOKKU_LETSENCRYPT_EMAIL=fzero@rubi.gd
dokku config:set --no-restart engine E_SERVICES_SENTRY_ENABLED=true
dokku config:set --no-restart engine E_SERVICES_SENTRY_PROJECT=2
dokku config:set --no-restart engine E_SERVICES_SENTRY_URL=https://sentry.qollenge.ir
dokku config:set --no-restart engine E_SERVICES_SENTRY_SECRET=${SENTRY_KEY}

// PUSH

dokku letsencrypt engine
