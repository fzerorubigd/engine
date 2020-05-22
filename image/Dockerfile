FROM golang:1.14-buster

RUN go get github.com/fzerorubigd/didebaan/cmd/...

WORKDIR /elbix.dev/engine

ENTRYPOINT ["/go/bin/didebaan"]