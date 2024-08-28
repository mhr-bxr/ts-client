.PHONY: all build 

BIN_DIR := ./bin
version := $(shell git rev-parse --short=12 HEAD)
timestamp := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
version := $(or $(version), $(shell cat /app/build-release | tr -d '\n'))

all: build

clean:
	rm -f $(BIN_DIR)/tsc

build: lint
	rm -f $(BIN_DIR)/tsc
	go build -o $(BIN_DIR)/tsc -v -ldflags \
		"-X main.rev=$(version) -X main.bts=$(timestamp)" cmd/tsc/main.go

lint:
	golangci-lint run

test: lint
	go test ./...
