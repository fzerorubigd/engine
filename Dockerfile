FROM golang:1.12-alpine

ADD . /go/src/github.com/fzerorubigd/engine

RUN apk add --no-cache --virtual .build-deps git gcc g++ libc-dev make \
    && apk add --no-cache ca-certificates bash \
    && cd /go/src/github.com/fzerorubigd/engine && make all \
    && apk del .build-deps

FROM alpine:3.6

COPY --from=0 /go/src/github.com/fzerorubigd/engine/bin/server /bin/
COPY --from=0 /go/src/github.com/fzerorubigd/engine/bin/migration /bin/
ADD scripts/server.sh /bin/server.sh
ADD scripts/Procfile /bin/Procfile
ADD scripts/CHECKS /bin/CHECKS
RUN chmod a+x /bin/server.sh

EXPOSE 80

WORKDIR /bin
