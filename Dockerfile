FROM golang:1.14-alpine

ENV GOFLAGS="-mod=readonly"
ADD . /go/src/elbix.dev/engine

RUN apk add --no-cache --virtual .build-deps git gcc g++ libc-dev make \
    && apk add --no-cache ca-certificates bash \
    && cd /go/src/elbix.dev/engine && make all \
    && apk del .build-deps

FROM alpine:3.6

ARG APP_NAME
ARG APP_PREFIX

COPY --from=0 /go/src/elbix.dev/engine/bin/server /bin/server
COPY --from=0 /go/src/elbix.dev/engine/bin/migration /bin/migration
ADD scripts/server.sh /bin/server.sh
ADD scripts/migration.sh /bin/migration.sh
RUN echo "web: /bin/sh /bin/server.sh" > /bin/Procfile
RUN echo "/v1/misc/health" > /bin/CHECKS
RUN chmod a+x /bin/server.sh /bin/migration.sh

EXPOSE 80

WORKDIR /bin
