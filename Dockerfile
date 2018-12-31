FROM golang:1.11-alpine

ADD . /balloon

# I don't need to set GOPATH since the Makefile takes care of that
RUN apk add --no-cache --virtual .build-deps git gcc g++ libc-dev make \
    && apk add --no-cache ca-certificates bash \
    && cd /balloon && make all \
    && apk del .build-deps

FROM alpine:3.6

COPY --from=0 /balloon/bin/server /bin/
COPY --from=0 /balloon/bin/migration /bin/
ADD scripts/dokku-run.sh /bin/run.sh

EXPOSE 80

CMD ["/bin/sh", "/bin/run.sh"]