#!/bin/bash

dokku plugin:install https://github.com/dokku/dokku-letsencrypt.git || dokku plugin:update letsencrypt
dokku plugin:install https://github.com/dokku/dokku-redis.git redis || dokku plugin:update redis
dokku plugin:install https://github.com/dokku/dokku-postgres.git postgres || dokku plugin:update postgres

dokku apps:create qollenge
dokku redis:create qollenge_redis
dokku redis:link qollenge_redis qollenge
dokku postgres:create qollenge_postgres
dokku postgres:link qollenge_postgres qollenge

dokku config:set --no-restart qollenge DOKKU_LETSENCRYPT_EMAIL=fzero@rubi.gd
dokku config:set --no-restart qollenge E_SERVICES_SENTRY_ENABLED=true
dokku config:set --no-restart qollenge E_SERVICES_SENTRY_PROJECT=2
dokku config:set --no-restart qollenge E_SERVICES_SENTRY_URL=https://sentry.elbix.dev
dokku config:set --no-restart qollenge E_SERVICES_SENTRY_SECRET=${SENTRY_KEY}

// PUSH

dokku letsencrypt cerulean
