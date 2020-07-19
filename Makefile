PROJECT_NAME := "pipe"
PKG := "behnama/pipe"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)
GOPATH := ${HOME}/go

.PHONY: all dep build clean test coverage coverhtml lint unit-test

all: build

wire:
	@wire ./...

mockery:
	@mockery -all -dir ./interfaces -output ./interfaces/mocks

lint-test:
	@golint -set_exit_status ${PKG_LIST}

vet-test:
	@go vet ${PKG_LIST}

fmt-test:
	@export unformatted=$$(gofmt -l ${GO_FILES} ) ;\
	if [ ! -z $$unformatted ] ; then echo $$unformatted ; exit 1; fi ;\
	exit 0 ;

fmt:
	@go fmt ./...

migrate: 
	@docker-compose rm -sfv db
	@docker-compose rm -sfv db-runner
	@docker volume rm -f $$(docker volume ls | grep api | awk '{print $$2}')
	@docker-compose up --force-recreate --build -d db

test: unit-test race-test

unit-test: export EXEC_MODE = TEST
unit-test: dep migrate mockery wire 
	@go test -count=1 -short ${PKG_LIST}

race-test: export EXEC_MODE = TEST
race-test: dep migrate mockery wire 
	@go test -race -count=1 -short ${PKG_LIST}

dep: 
	@go mod download; \
	go mod verify;

build: dep migrate mockery wire 
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./

build-docker: mockery wire 
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/api

clean: 
	@rm -f $(PROJECT_NAME)

run: build
	./api
# msan: dep
# 	@go test -msan -short ${PKG_LIST}

# coverage:
# 	./tools/coverage.sh;

# coverhtml: 
# 	./tools/coverage.sh html;
