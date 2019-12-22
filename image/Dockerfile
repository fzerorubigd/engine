FROM golang:1.13-buster

RUN go get github.com/fzerorubigd/didebaan/cmd/...

WORKDIR /elbix.dev/engine

ENTRYPOINT ["/go/bin/didebaan"]