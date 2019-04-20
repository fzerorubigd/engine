FROM golang:1.11-alpine

ADD . /go/src/github.com/fzerorubigd/balloon

ENV GO111MODULE=off
ENV GOPATH=/go

RUN apk add --no-cache --virtual .build-deps git gcc g++ libc-dev make \
    && apk add --no-cache ca-certificates bash \
    && cd /go/src/github.com/fzerorubigd/balloon && make all \
    && apk del .build-deps

FROM alpine:3.6

COPY --from=0 /go/src/github.com/fzerorubigd/balloon/bin/server /bin/
COPY --from=0 /go/src/github.com/fzerorubigd/balloon/bin/worker /bin/
COPY --from=0 /go/src/github.com/fzerorubigd/balloon/bin/migration /bin/
ADD scripts/server.sh /bin/server.sh
ADD scripts/worker.sh /bin/worker.sh

RUN chmod a+x /bin/server.sh /bin/worker.sh

EXPOSE 80

CMD ["/bin/server.sh"]
