#!/bin/bash

dokku plugin:install https://github.com/dokku/dokku-letsencrypt.git || dokku plugin:update letsencrypt
dokku plugin:install https://github.com/dokku/dokku-redis.git redis || dokku plugin:update redis
dokku plugin:install https://github.com/dokku/dokku-postgres.git postgres || dokku plugin:update postgres

dokku apps:create cerulean
dokku redis:create cerulean_redis
dokku redis:link cerulean_redis cerulean
dokku postgres:create cerulean_postgres
dokku postgres:link cerulean_postgres cerulean

dokku config:set --no-restart cerulean DOKKU_LETSENCRYPT_EMAIL=fzero@rubi.gd
dokku config:set --no-restart cerulean E_SERVICES_SENTRY_ENABLED=true
dokku config:set --no-restart cerulean E_SERVICES_SENTRY_PROJECT=3
dokku config:set --no-restart cerulean E_SERVICES_SENTRY_URL=https://sentry.elbix.dev
dokku config:set --no-restart cerulean E_SERVICES_SENTRY_SECRET=${SENTRY_KEY}

// PUSH

dokku letsencrypt qollenge
