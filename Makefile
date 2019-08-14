export ROOT:=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))
export BIN:=$(ROOT)/bin
export GOBIN:=$(BIN)
export PATH:=$(BIN):$(PATH)
export PROJECT=engine
export PROTOTOOL_VERSION=1.8.0
export DOKKU_HOST=elbix.dev
APP_NAME:=$(PROJECT)
DEFAULT_PASS=bita123
GO=$(shell which go)
GIT=$(shell which git)
CURL:=$(shell which curl)
CHMOD=$(shell which chmod)
DOCKER=$(shell which docker)
SSH=$(shell which ssh)
DB_PASS?=$(DEFAULT_PASS)
DB_USER?=$(APP_NAME)
DB_NAME?=$(APP_NAME)
WORK_DIR=$(ROOT)/tmp
LONG_HASH?=$(shell git log -n1 --pretty="format:%H" | cat)
SHORT_HASH?=$(shell git log -n1 --pretty="format:%h"| cat)
COMMIT_DATE?=$(shell git log -n1 --date="format:%D-%H-%I-%S" --pretty="format:%cd"| sed -e "s/\//-/g")
COMMIT_COUNT?=$(shell git rev-list HEAD --count| cat)
BUILD_DATE=$(shell date "+%D/%H/%I/%S"| sed -e "s/\//-/g")
VERSION="github.com/fzerorubigd/$(PROJECT)/pkg/version"
FLAGS="-X $(VERSION).hash=$(LONG_HASH) -X $(VERSION).short=$(SHORT_HASH) -X $(VERSION).date=$(COMMIT_DATE) -X $(VERSION).count=$(COMMIT_COUNT) -X $(VERSION).build=$(BUILD_DATE)"
LD_ARGS=-ldflags $(FLAGS)
GET=cd $(ROOT) && $(GO) get -u -v $(LD_ARGS)
BUILD=cd $(ROOT) && $(GO) build -v $(LD_ARGS)
INSTALL=cd $(ROOT) && $(GO) install -v $(LD_ARGS)
CG_SERVICES_POSTGRES_USER=$(DB_USER)
CG_SERVICES_POSTGRES_PASSWORD=$(DB_PASS)
CG_SERVICES_POSTGRES_DB=$(DB_NAME)
where-am-i = $(CURDIR)/$(word $(words $(MAKEFILE_LIST)),$(MAKEFILE_LIST))

# Default target is lint
lint: $(BIN)/golangci-lint
	LINT_GOGC=5 GOGC=5 $(BIN)/golangci-lint run


$(BIN)/jwtRS256.key:
	ssh-keygen -t rsa -b 4096 -m PEM -f $(BIN)/jwtRS256.key -N ''

$(BIN)/jwtRS256.key.bup: $(BIN)/jwtRS256.key
	openssl rsa -in $(BIN)/jwtRS256.key -pubout -outform PEM -out $(BIN)/jwtRS256.key.pub

rsa_file: $(BIN)/jwtRS256.key.bup $(BIN)/jwtRS256.key

$(BIN)/golangci-lint:
	$(CURL) -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(BIN) v1.17.1

clean:
	$(GIT) clean -fX ./

vendor:
	GO111MODULE=on $(GO) get ./cmd/... ./pkg/... ./modules/...
	GO111MODULE=on $(GO) mod tidy
	GO111MODULE=on $(GO) mod vendor

# Include modules make file
include $(wildcard $(ROOT)/modules/*/module.mk)

need_root:
	@[ "$(shell id -u)" -eq "0" ] || exit 1

not_root:
	@[ "$(shell id -u)" != "0" ] || exit 1

database-setup: need_root
	sudo -u postgres psql -U postgres -d postgres -c "CREATE USER $(DB_USER) WITH PASSWORD '$(DB_PASS)';" || sudo -u postgres psql -U postgres -d postgres -c "ALTER USER $(DB_USER) WITH PASSWORD '$(DB_PASS)';"
	sudo -u postgres psql -U postgres -d postgres -c "CREATE USER $(DB_USER)_test WITH PASSWORD '$(DB_PASS)';" || sudo -u postgres psql -U postgres -d postgres -c "ALTER USER $(DB_USER)_test WITH PASSWORD '$(DB_PASS)';"
	sudo -u postgres psql -U postgres -c "CREATE DATABASE $(DB_NAME);" || echo "Database $(DB_NAME) is already there?"
	sudo -u postgres psql -U postgres -c "CREATE DATABASE $(DB_NAME)_test;" || echo "Database $(DB_NAME)_test is already there?"
	sudo -u postgres psql -U postgres -c "GRANT ALL ON DATABASE $(DB_NAME) TO $(DB_USER);"
	sudo -u postgres psql -U postgres -c "GRANT ALL ON DATABASE $(DB_NAME)_test TO $(DB_USER)_test;"

$(BIN)/prototool:
	$(CURL) -sSL https://github.com/uber/prototool/releases/download/v$(PROTOTOOL_VERSION)/prototool-$(shell uname -s)-$(shell uname -m) -o $(BIN)/prototool
	$(CHMOD) +x $(BIN)/prototool

$(BIN)/protoc-gen-go:
	$(GET) github.com/golang/protobuf/protoc-gen-go

$(BIN)/protoc-gen-gogo:
	$(GET) github.com/gogo/protobuf/protoc-gen-gogo

$(BIN)/protoc-gen-grpc-gateway:
	$(GET) github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway

$(BIN)/protoc-gen-swagger:
	$(GET) github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

$(BIN)/protoc-gen-grpchan:
	$(GET) github.com/fullstorydev/grpchan/cmd/protoc-gen-grpchan

$(BIN)/go-bindata:
	$(GET) github.com/shuLhan/go-bindata/cmd/go-bindata

$(BIN)/reflex:
	$(GET) github.com/cespare/reflex

swagger-to-go:
	$(INSTALL) ./cmd/swagger-to-go

generators:
	$(INSTALL) ./cmd/protoc-gen-wrapper
	$(INSTALL) ./cmd/protoc-gen-model

tools-migration-qollenge: $(BIN)/go-bindata $(addsuffix -migration,$(dir $(wildcard $(ROOT)/modules/*/)))
	$(INSTALL) ./cmd/qollenge/qmigration

tools-migration-cerulean: $(BIN)/go-bindata $(addsuffix -migration,$(dir $(wildcard $(ROOT)/modules/*/)))
	$(INSTALL) ./cmd/cerulean/cmigration

mig-up-qollenge: tools-migration-qollenge
	$(BIN)/qmigration -action=up

mig-down-qollenge: tools-migration-qollenge
	$(BIN)/qmigration -action=down

mig-up-cerulean: tools-migration-cerulean
	$(BIN)/cmigration -action=up

mig-down-cerulean: tools-migration-cerulean
	$(BIN)/cmigration -action=down

test-qollenge: tools-migration-qollenge
	E_SERVICES_POSTGRES_USER="$(DB_USER)_test" E_SERVICES_POSTGRES_DB="$(DB_NAME)_test" $(BIN)/qmigration -action=down-all
	E_SERVICES_POSTGRES_USER="$(DB_USER)_test" E_SERVICES_POSTGRES_DB="$(DB_NAME)_test" $(BIN)/qmigration -action=up
	$(GO) test ./pkg/... ./modules/misc/... ./modules/user/... -coverprofile cover.cp
	E_SERVICES_POSTGRES_USER="$(DB_USER)_test" E_SERVICES_POSTGRES_DB="$(DB_NAME)_test" $(BIN)/qmigration -action=down-all

test-cerulean: tools-migration-cerulean
	E_SERVICES_POSTGRES_USER="$(DB_USER)_test" E_SERVICES_POSTGRES_DB="$(DB_NAME)_test" $(BIN)/cmigration -action=down-all
	E_SERVICES_POSTGRES_USER="$(DB_USER)_test" E_SERVICES_POSTGRES_DB="$(DB_NAME)_test" $(BIN)/cmigration -action=up
	$(GO) test ./pkg/... ./modules/misc/... ./modules/user/... ./modules/accounting/... -coverprofile cover.cp
	E_SERVICES_POSTGRES_USER="$(DB_USER)_test" E_SERVICES_POSTGRES_DB="$(DB_NAME)_test" $(BIN)/cmigration -action=down-all


proto: $(BIN)/prototool $(BIN)/protoc-gen-go $(BIN)/protoc-gen-grpc-gateway $(BIN)/protoc-gen-swagger $(BIN)/protoc-gen-grpchan $(BIN)/protoc-gen-gogo generators
	$(BIN)/prototool generate

swagger-ui: $(BIN)/go-bindata
	$(GIT) clone --depth 1 https://github.com/swagger-api/swagger-ui.git $(ROOT)/tmp/swagger-ui
	rm -rf $(ROOT)/third_party/swagger-ui
	mv $(ROOT)/tmp/swagger-ui/dist $(ROOT)/third_party/swagger-ui
	rm -rf $(ROOT)/tmp/swagger-ui
	sed -i -e 's/https:\/\/petstore.swagger.io\/v2\/swagger\.json/\/v1\/swagger\/index\.json/g' $(ROOT)/third_party/swagger-ui/index.html
	cd $(ROOT)/third_party/swagger-ui && $(BIN)/go-bindata -nometadata -o $(ROOT)/pkg/grpcgw/swagger.gen.go -nomemcopy=true -pkg=grpcgw ./...

swagger: swagger-to-go proto $(addsuffix -swagger,$(dir $(wildcard $(ROOT)/modules/*/)))

code-gen: swagger

build-all:
	@echo "Building all binaries"
	$(INSTALL) ./cmd/...

run-server-qollenge: code-gen build-all rsa_file
	@echo "Running..."
	E_MODULES_TOKEN_JWT_PRIVATE=$(shell cat $(BIN)/jwtRS256.key | base64 -w 0) E_MODULES_TOKEN_JWT_PUBLIC=$(shell cat $(BIN)/jwtRS256.key.pub | base64 -w 0) $(BIN)/qserver 2>&1

tools-migration: tools-migration-qollenge tools-migration-cerulean

all: build-all tools-migration

test: test-qollenge test-cerulean

watch: $(BIN)/reflex
	$(BIN)/reflex -r '\.proto$$' make code-gen

deploy-qollenge:
	$(DOCKER) build --build-arg APP_NAME=qollenge --build-arg APP_PREFIX=q -t dokku/qollenge:$(COMMIT_COUNT) .
	$(DOCKER) save dokku/qollenge:$(COMMIT_COUNT) | $(SSH) -o "StrictHostKeyChecking no" root@$(DOKKU_HOST) "docker load"
	$(SSH) -o "StrictHostKeyChecking no" root@$(DOKKU_HOST) "dokku tags:deploy qollenge $(COMMIT_COUNT)"

deploy-cerulean:
	$(DOCKER) build --build-arg APP_NAME=cerulean --build-arg APP_PREFIX=c -t dokku/cerulean:$(COMMIT_COUNT) .
	$(DOCKER) save dokku/cerulean:$(COMMIT_COUNT) | $(SSH) -o "StrictHostKeyChecking no" root@$(DOKKU_HOST) "docker load"
	$(SSH) -o "StrictHostKeyChecking no" root@$(DOKKU_HOST) "dokku tags:deploy cerulean $(COMMIT_COUNT)"

.PHONY: swagger-to-go proto swagger build-server run-server generate vendor
