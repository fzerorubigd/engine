export ROOT=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))
export BIN=$(ROOT)/bin
export GOPATH=$(abspath $(ROOT)/../..)
export GOBIN=$(BIN)
export PATH:=$(GOBIN):$(PATH)
APP_NAME=balloon
DEFAULT_PASS=bita123
GO=$(shell which go)
GIT=$(shell which git)
CURL:=$(shell which curl)
CHMOD=$(shell which chmod)
DB_PASS?=$(DEFAULT_PASS)
DB_USER?=$(APP_NAME)
DB_NAME?=$(APP_NAME)
WORK_DIR=$(ROOT)/tmp
LONG_HASH?=$(shell git log -n1 --pretty="format:%H" | cat)
SHORT_HASH?=$(shell git log -n1 --pretty="format:%h"| cat)
COMMIT_DATE?=$(shell git log -n1 --date="format:%D-%H-%I-%S" --pretty="format:%cd"| sed -e "s/\//-/g")
COMMIT_COUNT?=$(shell git rev-list HEAD --count| cat)
BUILD_DATE=$(shell date "+%D/%H/%I/%S"| sed -e "s/\//-/g")
FLAGS="-X version.hash=$(LONG_HASH) -X version.short=$(SHORT_HASH) -X version.date=$(COMMIT_DATE) -X version.count=$(COMMIT_COUNT) -X version.build=$(BUILD_DATE)"
LD_ARGS=-ldflags $(FLAGS)
GET=cd $(ROOT) && $(GO) get -u -v $(LD_ARGS)
CG_SERVICES_POSTGRES_USER=$(DB_USER)
CG_SERVICES_POSTGRES_PASSWORD=$(DB_PASS)
CG_SERVICES_POSTGRES_DB=$(DB_NAME)
where-am-i = $(CURDIR)/$(word $(words $(MAKEFILE_LIST)),$(MAKEFILE_LIST))

$(BIN)/prototool:
	$(CURL) -sSL https://github.com/uber/prototool/releases/download/v1.3.0/prototool-$(shell uname -s)-$(shell uname -m) -o $(BIN)/prototool
	$(CHMOD) +x $(BIN)/prototool

$(BIN)/protoc-gen-go:
	$(GET) github.com/golang/protobuf/protoc-gen-go

$(BIN)/protoc-gen-grpc-gateway:
	$(GET) github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway

$(BIN)/protoc-gen-swagger:
	$(GET) github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

proto: $(BIN)/prototool $(BIN)/protoc-gen-go $(BIN)/protoc-gen-grpc-gateway $(BIN)/protoc-gen-swagger
	$(BIN)/prototool all

